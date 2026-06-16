package daemon

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

type Session struct {
	ID     string
	Master *os.File
	Slave  *os.File
	Cmd    *exec.Cmd
}

type Daemon struct {
	socketPath string
	sessions   map[string]*Session
}

func NewDaemon(socketPath string) *Daemon {
	return &Daemon{
		socketPath: socketPath,
		sessions:   make(map[string]*Session),
	}
}

func (d *Daemon) Start() error {
	if _, err := os.Stat(d.socketPath); err == nil {
		os.Remove(d.socketPath)
	}

	ln, err := net.Listen("unix", d.socketPath)
	if err != nil {
		return err
	}

	fmt.Printf("vish daemon listening on %s\n", d.socketPath)

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		go d.handleClient(conn)
	}
}

func (d *Daemon) handleClient(conn net.Conn) {
	defer conn.Close()

	// 1. Read command from client (e.g. "START", "ATTACH", "LIST")
	// For now, let's assume simple start/attach logic
}

func OpenPTY() (master, slave *os.File, err error) {
	mfd, err := unix.Open("/dev/ptmx", unix.O_RDWR|unix.O_NOCTTY|unix.O_CLOEXEC, 0)
	if err != nil {
		return nil, nil, err
	}

	// Unlockpt (via TIOCSPTLCK)
	var lock int = 0
	if _, _, errno := syscall.Syscall(unix.SYS_IOCTL, uintptr(mfd), unix.TIOCSPTLCK, uintptr(unsafe.Pointer(&lock))); errno != 0 {
		unix.Close(mfd)
		return nil, nil, errno
	}

	// Get pts name (via TIOCGPTN)
	var ptsn int
	if _, _, errno := syscall.Syscall(unix.SYS_IOCTL, uintptr(mfd), unix.TIOCGPTN, uintptr(unsafe.Pointer(&ptsn))); errno != 0 {
		unix.Close(mfd)
		return nil, nil, errno
	}

	ptsName := fmt.Sprintf("/dev/pts/%d", ptsn)
	sfd, err := unix.Open(ptsName, unix.O_RDWR|unix.O_NOCTTY, 0)
	if err != nil {
		unix.Close(mfd)
		return nil, nil, err
	}

	return os.NewFile(uintptr(mfd), "/dev/ptmx"), os.NewFile(uintptr(sfd), ptsName), nil
}
