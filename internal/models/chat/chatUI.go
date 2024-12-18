package ChatModel

import (
	"fmt"
	"log"
	"strings"

	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sunikka/clich-client/internal/models/logging"
	"github.com/sunikka/clich-client/internal/theme"
	"github.com/sunikka/clich-client/internal/utils"
	"golang.org/x/net/websocket"
)

type wsMsg struct {
	// senderID string
	name    string
	message string
}

type ChatModel struct {
	viewport       viewport.Model
	Messages       []string
	textarea       textarea.Model
	senderStyle    lipgloss.Style
	recipientStyle lipgloss.Style
	Username       string
	connected      bool
	Ws             *websocket.Conn
	err            error
	app            *tea.Program
	Debug          *log.Logger
	msgCh          chan []byte
	theme          *theme.Theme
	styles         []lipgloss.Style
}

type (
	errMsg error
)

func NewChatModel(theme *theme.Theme) ChatModel {
	styles := applyStyles(theme)

	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "| "
	ta.CharLimit = 280

	ta.SetWidth(80)
	ta.SetHeight(3)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(80, 20)
	vp.SetContent(styles[secondaryStyle].Render(fmt.Sprintf("Welcome to the global chat!\nType a message and press enter to send.")))

	ta.KeyMap.InsertNewline.SetEnabled(false)
	return ChatModel{
		textarea:       ta,
		Messages:       []string{},
		viewport:       vp,
		senderStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("36")),
		recipientStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		connected:      false,
		err:            nil,
		Debug:          log.New(os.Stderr, "DEBUG: ", log.Lshortfile|log.LstdFlags),
		theme:          theme,
		styles:         styles,
	}
}

func (m ChatModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
		cmds  []tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)
	cmds = append(cmds, tiCmd, vpCmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			if m.connected == false {
				return m, func() tea.Msg { return WsErr{ErrMsg: "No websocket connection"} }
			}
			m.SendMessage(m.textarea.Value())
			cmds = append(cmds, logging.SendLogReq("Message sent!"))
			m.viewport.SetContent(strings.Join(m.Messages, "\n"))
			m.textarea.Reset()
			m.viewport.GotoBottom()
		}

	case utils.Message:
		formedMessage := fmt.Sprintf("%s %s", m.senderStyle.Render(msg.Username+": "), msg.Content)

		m.Messages = append(m.Messages, formedMessage)
		m.viewport.SetContent(strings.Join(m.Messages, "\n"))

		m.viewport.GotoBottom()

	case errMsg:
		m.err = msg
		return m, nil
	}

	if m.connected {
		cmds = append(cmds, m.TickMessageCheck())
	}

	return m, tea.Batch(cmds...)
}

func (m ChatModel) View() string {
	// Print the logo
	asciiArt, err := os.ReadFile("assets/ascii_art.txt")
	if err != nil {
		log.Println("Error loading the ascii art: ", err)
	}

	view := fmt.Sprintf(`
%s


%s
%s
	`,
		m.styles[primaryStyle].Width(50).Render(string(asciiArt)),
		m.viewport.View(),
		m.textarea.View(),
	)

	return m.styles[borderStyle].Width(50).Render(view)

}
