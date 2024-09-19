package loginUI

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	viewTypes "github.com/sunikka/clich-client/internal/models"
)

type (
	errMsg error
)

// UI elements
type elementType int

const (
	inputElement elementType = iota
	buttonElement
)

const (
	username = iota
	password
	login = iota
	guestLogin
	register
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
)

// TODO: Maybe turn menu buttons into their own Tea Model?
type UIElement struct {
	Type       elementType
	TextInput  textinput.Model
	ButtonText string
}

type Model struct {
	Elements []UIElement
	Focused  int
	Err      error
}

func InitialModel() Model {
	elements := make([]UIElement, 5)

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

	guestLoginBtn := "Enter as a guest"
	elements[guestLogin] = UIElement{Type: buttonElement, ButtonText: guestLoginBtn}

	registerBtn := "Sign up"
	elements[register] = UIElement{Type: buttonElement, ButtonText: registerBtn}

	return Model{
		Elements: elements,
		Focused:  0,
		Err:      nil,
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
		case tea.KeyShiftTab:
			m.prevInput()
		case tea.KeyEnter:
			// If on the last field, submit login (move on to next view)
			// otherwise move to the next field
			if m.Focused == login {
				return m, func() tea.Msg { return viewTypes.SwitchViewMsg{State: viewTypes.ChatView} }
			}
			if m.Focused == guestLogin {
				username := m.Elements[username].TextInput.Value()

				return m, func() tea.Msg {
					return viewTypes.SwitchViewMsg{
						State:    viewTypes.ChatView,
						Username: username,
					}
				}
			}
			if m.Focused == register {
				return m, func() tea.Msg { return viewTypes.SwitchViewMsg{State: viewTypes.ChatView} }
			}

			m.nextInput()

			for i := range m.Elements {
				if m.Elements[i].Type == inputElement {
					m.Elements[i].TextInput.Blur()
				}
			}

			if m.Elements[m.Focused].Type == inputElement {
				m.Elements[m.Focused].TextInput.Focus()
			}
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

	// TODO: Clean up the highlighting code
	loginBtn := continueStyle.Width(20).Render(m.Elements[login].ButtonText)
	if m.Focused == login {
		loginBtn = highlightStyle.Width(20).Render(m.Elements[login].ButtonText)
	}

	guestLoginBtn := continueStyle.Width(20).Render(m.Elements[guestLogin].ButtonText)
	if m.Focused == guestLogin {
		guestLoginBtn = highlightStyle.Width(20).Render(m.Elements[guestLogin].ButtonText)
	}

	registerBtn := continueStyle.Width(20).Render(m.Elements[register].ButtonText)
	if m.Focused == register {
		registerBtn = highlightStyle.Width(20).Render(m.Elements[register].ButtonText)
	}

	// Add fields and titles
	view := fmt.Sprintf(`
	%s


	%s
	%s

	%s
	%s

	%s
	%s
	%s
		`,

		inputStyle.Width(50).Render("---------------- Login to CLICH ----------------"),
		inputStyle.Width(45).Render("Username"),
		m.Elements[username].TextInput.View(),
		inputStyle.Width(45).Render("Password"),
		m.Elements[password].TextInput.View(),
		loginBtn,
		guestLoginBtn,
		registerBtn)
	view += continueStyle.Width(50).Render("\n\n Press ESC or CTRL+C to exit... \n")
	return view
}

func (m *Model) nextInput() {
	m.Focused = (m.Focused + 1) % len(m.Elements)
}

func (m *Model) prevInput() {
	m.Focused--

	if m.Focused < 0 {
		m.Focused = len(m.Elements) - 1
	}
}
