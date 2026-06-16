package pty

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"
	"vish/internal/spine"

	"golang.org/x/sys/unix"
)

const EX_SPINE_SLEEP = 104

type PTY struct {
	Master *os.File
	Slave  *os.File
}

func Open() (*PTY, error) {
	mfd, err := unix.Open("/dev/ptmx", unix.O_RDWR|unix.O_NOCTTY|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, err
	}

	// Unlockpt (via TIOCSPTLCK)
	var lock int = 0
	if _, _, errno := syscall.Syscall(unix.SYS_IOCTL, uintptr(mfd), unix.TIOCSPTLCK, uintptr(unsafe.Pointer(&lock))); errno != 0 {
		unix.Close(mfd)
		return nil, errno
	}

	// Get pts name (via TIOCGPTN)
	var ptsn int
	if _, _, errno := syscall.Syscall(unix.SYS_IOCTL, uintptr(mfd), unix.TIOCGPTN, uintptr(unsafe.Pointer(&ptsn))); errno != 0 {
		unix.Close(mfd)
		return nil, errno
	}

	ptsName := fmt.Sprintf("/dev/pts/%d", ptsn)
	sfd, err := unix.Open(ptsName, unix.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		unix.Close(mfd)
		return nil, err
	}

	return &PTY{
		Master: os.NewFile(uintptr(mfd), "/dev/ptmx"),
		Slave:  os.NewFile(uintptr(sfd), ptsName),
	}, nil
}

func (p *PTY) Run(c *exec.Cmd) error {
	c.Stdout = p.Slave
	c.Stderr = p.Slave
	c.Stdin = p.Slave
	c.SysProcAttr = &syscall.SysProcAttr{
		Setsid:  true,
		Setctty: true,
	}

	if err := c.Start(); err != nil {
		return err
	}

	// Close slave in parent
	p.Slave.Close()

	// Wait for process in a separate goroutine
	go func() {
		var waitStatus syscall.WaitStatus
		_, err := syscall.Wait4(c.Process.Pid, &waitStatus, 0, nil)
		if err != nil {
			return
		}

		if waitStatus.Exited() && waitStatus.ExitStatus() == EX_SPINE_SLEEP {
			fmt.Printf("\033[2K\r\033[2m[spine: Agent sleeping for 60s - Session Suspended]\033[0m\n")
			
			client := spine.NewClient("/tmp/spine.sock")
			var ns [16]byte
			copy(ns[:], "default")
			
			// Blocking call
			err := client.Park(ns, 60*time.Second, []byte("state-v1"))
			if err != nil {
				fmt.Printf("Park error: %v\n", err)
				return
			}

			fmt.Printf("\033[2K\r\033[1;32m[spine: Agent Woke Up - Resuming Session]\033[0m\n")
			
			// Resurrection: Re-spawn the command
			// In a real system, we'd use the state from the vault
			newCmd := exec.Command(c.Path, c.Args[1:]...)
			p.Run(newCmd)
		}
	}()

	return nil
}
