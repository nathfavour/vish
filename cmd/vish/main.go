package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"github.com/nathfavour/tony/pkg/identity"
	"github.com/nathfavour/vish/internal/daemon"
	"github.com/nathfavour/vish/internal/pty"
)

func main() {
	daemonMode := flag.Bool("daemon", false, "Start in daemon mode")
	flag.Parse()

	if *daemonMode {
		d := daemon.NewDaemon("/tmp/vish.sock")
		if err := d.Start(); err != nil {
			log.Fatalf("Daemon failed: %v", err)
		}
		return
	}

	// Interactive mode (Client)
	args := flag.Args()
	if len(args) == 0 {
		startShell()
		return
	}

	runCommand(args[0], args[1:])
}

func runCommand(command string, args []string) {
	p, err := pty.Open()
	if err != nil {
		log.Fatalf("Failed to open PTY: %v", err)
	}

	cmd := exec.Command(command, args...)
	if err := p.Run(cmd); err != nil {
		log.Fatalf("Failed to run command: %v", err)
	}

	// Relay between PTY and current terminal
	go func() {
		os.Stdout.ReadFrom(p.Master)
	}()
	p.Master.ReadFrom(os.Stdin)
}

func startShell() {
	fmt.Println("vish: Agentic Hybrid Shell [Tony Kernel Active]")
	fmt.Println("Type 'exit' to quit.")
	
	// Tony Kernel Initialization (Manned Master)
	// In a real scenario, this seed would be loaded from secure storage
	masterSeed := [32]byte{0xDE, 0xAD, 0xBE, 0xEF}
	idManager := identity.NewManager(masterSeed)
	_ = idManager

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("vish> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "exit" {
			break
		}
		if input == "" {
			continue
		}
		
		parts := strings.Fields(input)
		runCommand(parts[0], parts[1:])
	}
}
