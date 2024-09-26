package mainModel

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sunikka/clich-client/internal/auth"
	viewTypes "github.com/sunikka/clich-client/internal/models"
	chatUI "github.com/sunikka/clich-client/internal/models/chat"
	"github.com/sunikka/clich-client/internal/models/logging"
	loginUI "github.com/sunikka/clich-client/internal/models/login"
	"golang.org/x/net/websocket"
)

var p *tea.Program

type MainModel struct {
	app    *tea.Program
	state  viewTypes.SessionState
	login  tea.Model
	chat   tea.Model
	wsConn *websocket.Conn

	logger logging.Model

	// User related
	username string
	token    string
	loggedIn bool
}

func NewMainModel() tea.Model {

	return MainModel{
		state:  viewTypes.LoginView,
		login:  loginUI.InitialModel(20),
		logger: logging.NewLogWindow(100, 20),
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case viewTypes.SwitchViewMsg:
		m.state = viewTypes.SessionState(msg.State)

	case auth.LoginRequest:
		// return m, tea.Batch(auth.SendLoginRequestCmd(msg.Username, msg.Password))
		cmds = append(cmds, m.logger.Log("Login Request sent"), auth.SendLoginRequestCmd(msg.Username, msg.Password))

	case auth.LoginSuccess:
		cmds = append(cmds,
			m.establishChatConnection(),
			m.startMsgHandler(),
			m.SwitchView(viewTypes.ChatView),
		)

		m.username = msg.Username
		m.token = msg.Token
		m.loggedIn = true

	case auth.LoginFailure:
		cmds = append(cmds, m.logger.Log(msg.Error))
	case WsConnected:
		cmds = append(cmds, m.logger.Log("Websocket connected!"))
	}

	switch m.state {

	case viewTypes.LoginView:
		m.login, cmd = m.login.Update(msg)
		cmds = append(cmds, cmd)

	case viewTypes.ChatView:
		m.chat, cmd = m.chat.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update logger
	updatedLogger, cmd := m.logger.Update(msg)
	m.logger = updatedLogger.(logging.Model)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	logView := m.logger.View()

	switch m.state {
	case viewTypes.ChatView:
		m.chat = chatUI.InitialModel(m.wsConn)
		return m.chat.View()
	default:
		//		return m.login.View()
		return lipgloss.JoinHorizontal(lipgloss.Top, m.login.View(), logView)
	}

}

func (m MainModel) SwitchView(view viewTypes.SessionState) tea.Cmd {
	return func() tea.Msg {
		return viewTypes.SwitchViewMsg{
			State: view,
		}
	}
}
