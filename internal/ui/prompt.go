package ui

import (
	"os"
	"strings"
	"vish/internal/ecosystem"

	"github.com/charmbracelet/lipgloss"
)

var (
	PromptStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	ArrowStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	ManagedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Italic(true)
	UnmanagedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

// GetPrompt returns the styled prompt string.
func GetPrompt() string {
	managed, _ := ecosystem.IsManaged()
	status := ""
	if managed {
		status = ManagedStyle.Render(" (managed)")
	}

	cwd, _ := os.Getwd()
	// Get last part of path
	parts := strings.Split(cwd, "/")
	dir := parts[len(parts)-1]
	if dir == "" {
		dir = "/"
	}

	// Dynamic color for the directory
	dirStyle := GetColor(dir).Bold(true)

	return PromptStyle.Render("vish") + status + " " + dirStyle.Render(dir) + " " + ArrowStyle.Render("❯") + " "
}
