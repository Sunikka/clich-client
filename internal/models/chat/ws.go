package ChatModel

import (
	"encoding/json"
	"fmt"
	"io"
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

	m.msgCh = make(chan []byte)

	go func() {
		defer close(m.msgCh)
		for {
			var msg []byte
			err := websocket.Message.Receive(m.Ws, &msg)
			if err != nil {
				// TODO: Handle this more gracefully
				if err == io.EOF {
					log.Fatal("Websocket connection closed by server")
				}

				// TODO: Add log channel later to log error
				fmt.Println(err)
				continue
			}
			m.msgCh <- msg
		}
	}()

	return nil
}

// tickMessageCheck schedules periodic polling of WebSocket messages using tea.Tick.
func (m *ChatModel) TickMessageCheck() tea.Cmd {
	// This returns a tea.Tick which checks the message channel every 100ms
	return tea.Tick(time.Millisecond*100, func(time.Time) tea.Msg {
		select {
		case msg, ok := <-m.msgCh:
			if !ok {
				// If the channel is closed, indicate WebSocket closure
				return WsErr{ErrMsg: "WebSocket connection closed"}
			}
			// Return the received message as tea.Msg to the update loop
			// return utils.Message{Username: m.Username, Content: msg}
			return m.ParseMessage(msg)
		default:
			return nil
		}

	})
}

func (m ChatModel) SendMessage(msg string) error {
	payload := Message{
		Content: msg,
		SentAt:  time.Now().UTC(),
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = m.Ws.Write([]byte(payloadJSON))
	if err != nil {
		return err
	}

	return nil
}

func (m ChatModel) ParseMessage(buf []byte) tea.Msg {
	var msgParsed MessageReceived

	err := json.Unmarshal(buf, &msgParsed)
	if err != nil {
		return ReturnWsErr(err)
	}

	return utils.Message{Username: msgParsed.Sender, Content: msgParsed.Content, SentAt: msgParsed.SentAt}
}

// TODO: Replace this by protobuf
type Message struct {
	Content string    `json:"content"`
	SentAt  time.Time `json:"sent_at"`
}

type MessageReceived struct {
	Sender  string    `json:"sender"`
	Content string    `json:"content"`
	SentAt  time.Time `json:"sent_at"`
}
