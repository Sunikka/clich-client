package chatUI

import (
	"fmt"
	"log"
	"strings"

	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sunikka/clich-client/internal/utils"
	"golang.org/x/net/websocket"
)

var testNames = []string{"Paavo", "Gigachad", "Harold", "Taavi", "MogBot"}

// Styling
const (
	primaryColor   = lipgloss.Color("#32CD32")
	secondaryColor = lipgloss.Color("#767676")
	highlightColor = lipgloss.Color("#FFFFFF")
)

var (
	primaryStyle   = lipgloss.NewStyle().Foreground(primaryColor)
	secondaryStyle = lipgloss.NewStyle().Foreground(secondaryColor)
	highlightStyle = lipgloss.NewStyle().Foreground(highlightColor)
)

type wsMsg struct {
	// senderID string
	name    string
	message string
}

type Model struct {
	viewport       viewport.Model
	messages       []string
	textarea       textarea.Model
	senderStyle    lipgloss.Style
	recipientStyle lipgloss.Style
	ws             *websocket.Conn
	err            error
}

type (
	errMsg error
)

func InitialModel(ws *websocket.Conn) Model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = 280

	ta.SetWidth(80)
	ta.SetHeight(3)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(80, 23)
	vp.SetContent(secondaryStyle.Render(fmt.Sprintf("Welcome to the global chat!\nType a message and press enter to send.")))

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return Model{
		textarea:       ta,
		messages:       []string{},
		viewport:       vp,
		senderStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("36")),
		recipientStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		ws:             ws,
		err:            nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			_, err := m.ws.Write([]byte(m.textarea.Value()))
			if err != nil {
				log.Println("Error sending message: ", err)
			}
			// m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())

			// For testing
			// m.messages = append(m.messages, m.recipientStyle.Render("MogBot: ")+"Based")

			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}
	case utils.Message:

		m.messages = append(m.messages, m.senderStyle.Render(msg.Username+": ")+msg.Content)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))

		m.viewport.GotoBottom()

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)

}

func (m Model) View() string {
	// Print the logo
	asciiArt, err := os.ReadFile("assets/ascii_art.txt")
	if err != nil {
		log.Println("Error loading the ascii art: ", err)
	}

	return fmt.Sprintf(`
%s


%s
%s
	`,
		primaryStyle.Width(50).Render(string(asciiArt)),
		m.viewport.View(),
		m.textarea.View(),
	)
}
