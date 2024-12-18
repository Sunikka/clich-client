package loginUI

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sunikka/clich-client/internal/auth"
	"github.com/sunikka/clich-client/internal/models/logging"
	"github.com/sunikka/clich-client/internal/theme"
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
	theme        *theme.Theme
	styles       []lipgloss.Style
	Err          error
}

func InitialModel(windowHeight int, theme *theme.Theme) Model {
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
		theme:        theme,
		styles:       applyStyles(theme),
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
				uname := m.Elements[username].TextInput.Value()
				pw := m.Elements[password].TextInput.Value()
				return m, func() tea.Msg { return auth.RegisterRequest{Username: uname, Password: pw} }
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

	disclaimer := "DISCLAIMER: This is a test environment. There is no encryption for your messages. Do not use a password you use in other services and do not type anything sensitive in the chat!"

	formWidth := 50

	// TODO: Clean up the highlighting code
	loginBtn := m.styles[continueStyle].Width(20).Render(m.Elements[login].ButtonText)
	if m.Focused == login {
		loginBtn = m.styles[highlightStyle].Width(20).Render(m.Elements[login].ButtonText)
	}
	registerBtn := m.styles[continueStyle].Width(20).Render(m.Elements[register].ButtonText)
	if m.Focused == register {
		registerBtn = m.styles[highlightStyle].Width(20).Render(m.Elements[register].ButtonText)
	}

	// Add fields and titles
	view := fmt.Sprintf(`%s


	%s
	%s

	%s
	%s

	%s
	%s


%s

		`,

		m.styles[inputStyle].Width(formWidth).Render("---------------- Login to CLICH ----------------"),
		m.styles[inputStyle].Width(formWidth-5).Render("Username"),
		m.Elements[username].TextInput.View(),
		m.styles[inputStyle].Width(formWidth-5).Render("Password"),
		m.Elements[password].TextInput.View(),
		loginBtn,
		registerBtn,
		disclaimer)
	view += m.styles[continueStyle].Width(formWidth).Render("\n\n Press ESC or CTRL+C to exit... \n")

	// Add styles (margin & border)
	view = m.styles[paddingStyle].Width(formWidth).Render(view)
	view = m.styles[marginStyle].Width(formWidth).Render(view)
	view = m.styles[borderStyle].Width(formWidth).Render(view)

	return m.styles[borderStyle].Width(formWidth + marginLeft).Height(m.WindowHeight).Render(view)
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
