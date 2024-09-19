package mainUI

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	viewTypes "github.com/sunikka/clich-client/internal/models"
	chatUI "github.com/sunikka/clich-client/internal/models/chat"
	loginUI "github.com/sunikka/clich-client/internal/models/login"
	"golang.org/x/net/websocket"
)

var p *tea.Program

type MainModel struct {
	state viewTypes.SessionState
	login tea.Model
	chat  tea.Model
}

func NewMainModel(ws *websocket.Conn) tea.Model {

	return MainModel{
		state: viewTypes.LoginView,
		login: loginUI.InitialModel(),

		chat: chatUI.InitialModel(ws),
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	// var cmds []tea.Cmd

	switch msg := msg.(type) {
	case viewTypes.SwitchViewMsg:
		fmt.Println(msg)
		m.state = viewTypes.SessionState(msg.State)
	}

	switch m.state {
	case viewTypes.LoginView:
		m.login, cmd = m.login.Update(msg)
		return m, cmd
	case viewTypes.ChatView:
		m.chat, cmd = m.chat.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m MainModel) View() string {
	switch m.state {
	case viewTypes.ChatView:
		return m.chat.View()

	default:
		return m.login.View()
	}

}
