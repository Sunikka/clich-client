package loginUI

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/sunikka/clich-client/internal/theme"
)

type style int

const (
	marginLeft = 20
)

const (
	inputStyle style = iota
	continueStyle
	highlightStyle
	borderStyle
	marginStyle
	paddingStyle
	styleCount // Length of the styles slice
)

func applyStyles(theme *theme.Theme) []lipgloss.Style {
	styles := make([]lipgloss.Style, styleCount)
	styles[inputStyle] = lipgloss.NewStyle().Foreground(theme.PrimaryColor)
	styles[continueStyle] = lipgloss.NewStyle().Foreground(theme.SecondaryColor)
	styles[highlightStyle] = lipgloss.NewStyle().Foreground(theme.HighlightColor)
	styles[borderStyle] = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1).BorderStyle(lipgloss.HiddenBorder())

	styles[marginStyle] = lipgloss.NewStyle().MarginRight(marginLeft)
	styles[paddingStyle] = lipgloss.NewStyle().PaddingRight(3).PaddingLeft(3)

	return styles
}
