package logging

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styling
var (
	paddingStyle = lipgloss.NewStyle().PaddingRight(3).PaddingLeft(3)
	borderStyle  = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1).BorderStyle(lipgloss.HiddenBorder())
)

type logMsg struct {
	logText string
}

type Model struct {
	Window       viewport.Model
	WindowWidth  int
	WindowHeight int
	MaxLogs      int
	Logs         []string
}

func NewLogWindow(width, height int) Model {
	window := viewport.New(width, height)

	return Model{
		Window:       window,
		WindowWidth:  width,
		WindowHeight: height,
		MaxLogs:      height - 1,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case logMsg:
		m.Logs = append(m.Logs, msg.logText)

		if len(m.Logs) > m.MaxLogs {
			m.Logs = m.Logs[1:]
		}

		m.Window.SetContent(m.formatLogs())
		m.Window.GotoBottom()
	}

	return m, nil
}

func (m Model) View() string {
	logView := paddingStyle.Render(m.formatLogs())

	return borderStyle.Width(m.WindowWidth).Height(m.WindowHeight).Render(logView)
}

func (m Model) formatLogs() string {
	if len(m.Logs) == 0 {
		return fmt.Sprintf("%s\n", "no logs yet")
	}

	var logs string

	for _, msg := range m.Logs {
		logs += fmt.Sprintf("%s\n", msg)
	}

	return fmt.Sprintf("%s\n", logs)

}

func (m Model) Log(logStr string) tea.Cmd {
	return func() tea.Msg {
		return logMsg{logText: logStr}
	}
}
