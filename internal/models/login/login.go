package loginUI

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sunikka/clich-client/internal/auth"
	viewTypes "github.com/sunikka/clich-client/internal/models"
	"github.com/sunikka/clich-client/internal/models/logging"
)

type (
	errMsg error
)

// UI elements
type elementType int

const (
	inputElement elementType = iota
	buttonElement
	logElement
)

const (
	username = iota
	password
	login = iota
	register
)

const (
	marginLeft = 20
)

// Styling
const (
	primaryColor   = lipgloss.Color("#32CD32")
	secondaryColor = lipgloss.Color("#767676")
	highlightColor = lipgloss.Color("#FFFFFF")
)

var (
	inputStyle     = lipgloss.NewStyle().Foreground(primaryColor)
	continueStyle  = lipgloss.NewStyle().Foreground(secondaryColor)
	highlightStyle = lipgloss.NewStyle().Foreground(highlightColor)
	borderStyle    = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).Padding(1).BorderStyle(lipgloss.HiddenBorder())
	marginStyle    = lipgloss.NewStyle().MarginRight(marginLeft)
	paddingStyle   = lipgloss.NewStyle().PaddingRight(3).PaddingLeft(3)
)

// TODO: Maybe turn menu buttons into their own Tea Model?
type UIElement struct {
	Type       elementType
	TextInput  textinput.Model
	ButtonText string
}

type Model struct {
	Elements     []UIElement
	LogWindow    logging.Model
	WindowHeight int
	logMessages  []string
	Focused      int
	Err          error
}

func InitialModel(windowHeight int) Model {
	elements := make([]UIElement, 4)

	usernameInput := textinput.New()
	usernameInput.Placeholder = "username"
	usernameInput.Focus()
	elements[username] = UIElement{Type: inputElement, TextInput: usernameInput}

	passwordInput := textinput.New()
	passwordInput.Placeholder = "password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'
	elements[password] = UIElement{Type: inputElement, TextInput: passwordInput}

	loginBtn := "Sign in"
	elements[login] = UIElement{Type: buttonElement, ButtonText: loginBtn}

	registerBtn := "Sign up"
	elements[register] = UIElement{Type: buttonElement, ButtonText: registerBtn}

	return Model{
		Elements:     elements,
		Focused:      0,
		WindowHeight: windowHeight,
		Err:          nil,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.Elements))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyTab:
			m.nextInput()
			m.updateInputElement()

		case tea.KeyShiftTab:
			m.prevInput()
			m.updateInputElement()
		case tea.KeyEnter:
			// If on the last field, submit login (move on to next view)
			// otherwise move to the next field
			if m.Focused == login {
				//return m, func() tea.Msg { return viewTypes.SwitchViewMsg{State: viewTypes.ChatView} }
				uname := m.Elements[username].TextInput.Value()
				pw := m.Elements[password].TextInput.Value()
				return m, func() tea.Msg { return auth.LoginRequest{Username: uname, Password: pw} }

			}
			if m.Focused == register {
				return m, func() tea.Msg { return viewTypes.SwitchViewMsg{State: viewTypes.ChatView} }
			}

			m.nextInput()
			m.updateInputElement()
		}
	case errMsg:
		m.Err = msg
		return m, nil
	}

	for i := range m.Elements {
		if m.Elements[i].Type == inputElement {
			m.Elements[i].TextInput, cmds[i] = m.Elements[i].TextInput.Update(msg)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	formWidth := 50

	// TODO: Clean up the highlighting code
	loginBtn := continueStyle.Width(20).Render(m.Elements[login].ButtonText)
	if m.Focused == login {
		loginBtn = highlightStyle.Width(20).Render(m.Elements[login].ButtonText)
	}
	registerBtn := continueStyle.Width(20).Render(m.Elements[register].ButtonText)
	if m.Focused == register {
		registerBtn = highlightStyle.Width(20).Render(m.Elements[register].ButtonText)
	}

	// Add fields and titles
	view := fmt.Sprintf(`%s


	%s
	%s

	%s
	%s

	%s
	%s





		`,

		inputStyle.Width(formWidth).Render("---------------- Login to CLICH ----------------"),
		inputStyle.Width(formWidth-5).Render("Username"),
		m.Elements[username].TextInput.View(),
		inputStyle.Width(formWidth-5).Render("Password"),
		m.Elements[password].TextInput.View(),
		loginBtn,
		registerBtn)
	view += continueStyle.Width(formWidth).Render("\n\n Press ESC or CTRL+C to exit... \n")

	// Add styles (margin & border)
	view = paddingStyle.Width(formWidth).Render(view)
	view = marginStyle.Width(formWidth).Render(view)
	view = borderStyle.Width(formWidth).Render(view)

	return borderStyle.Width(formWidth + marginLeft).Height(m.WindowHeight).Render(view)
}

func (m *Model) nextInput() {
	if m.Focused < len(m.Elements)-1 {
		m.Focused = (m.Focused + 1) % len(m.Elements)
	} else {
		m.Focused = 0
	}
}

func (m *Model) prevInput() {
	m.Focused--

	if m.Focused < 0 {
		m.Focused = len(m.Elements) - 1
	}
}

func (m *Model) updateInputElement() {
	for i := range m.Elements {
		if m.Elements[i].Type == inputElement {
			m.Elements[i].TextInput.Blur()
		}
	}

	if m.Elements[m.Focused].Type == inputElement {
		m.Elements[m.Focused].TextInput.Focus()
	}

}
