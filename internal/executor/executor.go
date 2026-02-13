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
	Writer io.Writer
}

func (cw *ColorizedWriter) Write(p []byte) (n int, err error) {
	// If the input is too small, just write it
	if len(p) < 2 {
		return cw.Writer.Write(p)
	}

	// Simple heuristic: if there's no newline, color the whole chunk
	if !strings.Contains(string(p), "\n") {
		colored := ui.GetColor(string(p)).Render(string(p))
		_, err = cw.Writer.Write([]byte(colored))
		return len(p), err
	}

	lines := strings.Split(string(p), "\n")
	for i, line := range lines {
		if line == "" && i < len(lines)-1 {
			cw.Writer.Write([]byte("\n"))
			continue
		}
		if line == "" {
			continue
		}
		
		coloredLine := ui.GetColor(line).Render(line)
		if i < len(lines)-1 {
			coloredLine += "\n"
		}
		_, err = cw.Writer.Write([]byte(coloredLine))
		if err != nil {
			return 0, err
		}
	}
	return len(p), nil
}

func Execute(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	
	// Use colorized writers for stdout and stderr
	cmd.Stdout = &ColorizedWriter{Writer: os.Stdout}
	cmd.Stderr = &ColorizedWriter{Writer: os.Stderr}
	cmd.Stdin = os.Stdin
	
	return cmd.Run()
}
