package mainModel

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sunikka/clich-client/internal/auth"
	viewTypes "github.com/sunikka/clich-client/internal/models"
	chatUI "github.com/sunikka/clich-client/internal/models/chat"
	"github.com/sunikka/clich-client/internal/models/logging"
	loginUI "github.com/sunikka/clich-client/internal/models/login"
	"github.com/sunikka/clich-client/internal/utils"
)

var p *tea.Program

type MainModel struct {
	app   *tea.Program
	state viewTypes.SessionState
	login tea.Model
	chat  chatUI.ChatModel

	logger logging.Model

	// User related
	username string
	token    string
	loggedIn bool
}

func NewMainModel(app *tea.Program) tea.Model {

	return MainModel{
		app:    app,
		state:  viewTypes.LoginView,
		chat:   chatUI.NewChatModel(),
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

	case auth.RegisterRequest:
		cmds = append(cmds, m.logger.Log("Register Request sent"), auth.SendRegisterRequestCmd(msg.Username, msg.Password))

	case auth.RegisterSuccess:
		cmds = append(cmds, m.logger.Log(msg.Message))

	case auth.RegisterFailure:
		cmds = append(cmds, m.logger.Log(msg.Error))

	case auth.LoginRequest:
		// return m, tea.Batch(auth.SendLoginRequestCmd(msg.Username, msg.Password))
		cmds = append(cmds, m.logger.Log("Login Request sent"), auth.SendLoginRequestCmd(msg.Username, msg.Password))

	case auth.LoginSuccess:
		m.username = msg.Username
		m.token = msg.Token
		m.loggedIn = true
		chatWithConn, cmd := m.chat.Connect(m.username)
		m.chat = chatWithConn
		m.chat.Username = msg.Username
		cmds = append(cmds,
			m.logger.Log(fmt.Sprintf("Login succesful, welcome %s!", msg.Username)),
			m.SwitchView(viewTypes.ChatView),
			cmd,
		)

	case auth.LoginFailure:
		cmds = append(cmds, m.logger.Log(msg.Error))

	case chatUI.WsConnected:
		updatedChat, cmd := m.chat.Update(msg)
		m.chat = updatedChat.(chatUI.ChatModel)
		cmds = append(cmds, m.logger.Log("Connection established!"), cmd)

		startMsgHandlerCmd := m.chat.StartMsgHandler()
		cmds = append(cmds, startMsgHandlerCmd)
	case chatUI.WsErr:
		cmds = append(cmds, m.logger.Log(fmt.Sprintf("Websocket connection failed: %v", msg.ErrMsg)))
	case utils.Message:
		// updatedChat, _ := m.chat.Update(msg)
		// m.chat = updatedChat.(chatUI.ChatModel)

		cmds = append(cmds, m.logger.Log(fmt.Sprintf("Message received from %s", msg.Username)))

	case logging.LogRequest:
		cmds = append(cmds, m.logger.Log(msg.LogText))
	}

	switch m.state {

	case viewTypes.LoginView:
		m.login, cmd = m.login.Update(msg)
		cmds = append(cmds, cmd)

	case viewTypes.ChatView:
		updatedChat, cmd := m.chat.Update(msg)
		m.chat = updatedChat.(chatUI.ChatModel)
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
		//	return m.chat.View()
		return lipgloss.JoinHorizontal(lipgloss.Top, m.chat.View(), logView)
	default:
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
