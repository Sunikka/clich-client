package ChatModel

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/sunikka/clich-client/internal/theme"
)

type style int

const (
	marginLeft = 20
)

const (
	primaryStyle style = iota
	secondaryStyle
	highlightStyle
	borderStyle
	styleCount // Length of the styles slice
)

func applyStyles(theme *theme.Theme) []lipgloss.Style {
	styles := make([]lipgloss.Style, styleCount)
	styles[primaryStyle] = lipgloss.NewStyle().Foreground(theme.PrimaryColor)
	styles[secondaryStyle] = lipgloss.NewStyle().Foreground(theme.SecondaryColor)
	styles[highlightStyle] = lipgloss.NewStyle().Foreground(theme.HighlightColor)
	styles[borderStyle] = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1).BorderStyle(lipgloss.HiddenBorder())

	return styles
}
