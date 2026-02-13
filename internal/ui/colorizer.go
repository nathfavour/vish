package ui

import (
	"hash/fnv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Palette of vibrant colors
var palette = []string{
	"5",  // Magenta
	"6",  // Cyan
	"9",  // Red
	"10", // Green
	"11", // Yellow
	"12", // Blue
	"13", // Purple
	"14", // Aqua
	"15", // White
	"190", // Gold
	"201", // Hot Pink
	"208", // Orange
	"214", // Deep Orange
	"118", // Lime
	"45",  // Turquoise
}

// GetColor returns a lipgloss.Style based on the hash of the input string
func GetColor(s string) lipgloss.Style {
	h := fnv.New32a()
	h.Write([]byte(s))
	hash := h.Sum32()
	colorIndex := hash % uint32(len(palette))
	return lipgloss.NewStyle().Foreground(lipgloss.Color(palette[colorIndex]))
}

// Highlight syntax-highlights a shell command string dynamically
func Highlight(input string) string {
	if strings.TrimSpace(input) == "" {
		return input
	}

	var builder strings.Builder
	words := strings.Fields(input)
	
	// Simple word-based highlighting for real-time vibe
	// We can get more complex with the AST later
	for i, word := range words {
		if i > 0 {
			builder.WriteString(" ")
		}
		
		// If it looks like a flag
		if strings.HasPrefix(word, "-") {
			builder.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("244")).Render(word))
			continue
		}

		// Use the hash-based color for dynamic vibe
		builder.WriteString(GetColor(word).Render(word))
	}

	// Handle trailing spaces to keep typing natural
	if strings.HasSuffix(input, " ") {
		builder.WriteString(" ")
	}

	return builder.String()
}
