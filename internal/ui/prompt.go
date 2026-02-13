package ui

import (
	"vish/internal/ecosystem"

	"github.com/charmbracelet/lipgloss"
)

var (
	PromptStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	ArrowStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	ManagedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Italic(true)
	UnmanagedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func GetPrompt() string {
	managed, _ := ecosystem.IsManaged()
	status := ""
	if managed {
		status = ManagedStyle.Render(" (managed)")
	}

	return PromptStyle.Render("vish") + status + " " + ArrowStyle.Render("❯") + " "
}
