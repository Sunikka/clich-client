package ChatModel

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/sunikka/clich-client/internal/models/logging"
	"github.com/sunikka/clich-client/internal/utils"
	"golang.org/x/net/websocket"
)

type WsConnected struct{}

type WsErr struct {
	ErrMsg string
}

type UserInfo struct {
	ID     uuid.UUID
	Name   string
	Active bool
	// token jwt.Token
}

func ReturnWsConnected() tea.Cmd {
	return func() tea.Msg {
		return WsConnected{}
	}
}

func ReturnWsErr(err error) tea.Cmd {
	return func() tea.Msg {
		return WsErr{ErrMsg: fmt.Sprintf("Websocket error: %v", err)}
	}
}

func (m ChatModel) Connect(username string) (ChatModel, tea.Cmd) {
	// Websocket connection
	url := os.Getenv("SERVER_URL_WS")
	origin := os.Getenv("CLIENT_ORIGIN")
	conn, err := websocket.Dial(url, "", origin)
	if err != nil {
		return m, ReturnWsErr(err)
	}

	if conn == nil {
		return m, ReturnWsErr(err)
	}

	m.Ws = conn
	m.connected = true

	user := UserInfo{
		ID:     uuid.New(),
		Name:   username,
		Active: true,
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return m, ReturnWsErr(err)
	}

	err = websocket.Message.Send(m.Ws, userJSON)
	if err != nil {
		return m, ReturnWsErr(err)
	}
	return m, ReturnWsConnected()
}

func (m *ChatModel) StartMsgHandler() tea.Cmd {

	if m.Ws == nil {
		return func() tea.Msg {
			return tea.Batch(
				logging.SendLogReq("Websocket is nil"),
				func() tea.Msg { return WsErr{ErrMsg: "Failed to start msg handler, websocket nil"} },
			)
		}
	}

	m.msgCh = make(chan string)

	go func() {
		defer close(m.msgCh)
		for {
			var msg string
			err := websocket.Message.Receive(m.Ws, &msg)
			if err != nil {
				// Add log channel later to log error
				fmt.Println(err)
				continue
			}
			m.msgCh <- msg
		}
	}()

	// 	return func() tea.Msg { return utils.Message{Username: username, Content: "test message"} }

	return m.TickMessageCheck()
}

func (m *ChatModel) testMsg(username string) tea.Cmd {
	return func() tea.Msg { return utils.Message{Username: username, Content: "test message"} }
}

// tickMessageCheck schedules periodic polling of WebSocket messages using tea.Tick.
func (m *ChatModel) TickMessageCheck() tea.Cmd {
	// This returns a tea.Tick which checks the message channel every 100ms
	return tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg {
		select {
		case msg, ok := <-m.msgCh:
			log.Println("Message got in ticker")
			if !ok {
				// If the channel is closed, indicate WebSocket closure
				return WsErr{ErrMsg: "WebSocket connection closed"}
			}
			// Return the received message as tea.Msg to the update loop
			return tea.Batch(
				func() tea.Msg { return utils.Message{Username: m.Username, Content: msg} },
				m.TickMessageCheck(), // Schedule the next tick
			)
		default:
			return m.TickMessageCheck()
		}

	})
}
