package executor

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"vish/internal/ui"
)

// ColorizedWriter wraps an io.Writer to add dynamic color to lines
type ColorizedWriter struct {
	w io.Writer
}

func (cw *ColorizedWriter) Write(p []byte) (n int, err error) {
	lines := strings.Split(string(p), "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		
		// Apply dynamic color based on line content
		coloredLine := ui.GetColor(line).Render(line)
		if i < len(lines)-1 {
			coloredLine += "\n"
		}
		_, err = cw.w.Write([]byte(coloredLine))
		if err != nil {
			return 0, err
		}
	}
	return len(p), nil
}

func Execute(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	
	// Use colorized writers for stdout and stderr
	cmd.Stdout = &ColorizedWriter{w: os.Stdout}
	cmd.Stderr = &ColorizedWriter{w: os.Stderr}
	cmd.Stdin = os.Stdin
	
	return cmd.Run()
}
