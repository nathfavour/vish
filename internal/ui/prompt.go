package ui

import "github.com/charmbracelet/lipgloss"

var (
	PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Bold(true)
	ArrowStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
)

func GetPrompt() string {
	return PromptStyle.Render("vish") + " " + ArrowStyle.Render("❯") + " "
}
