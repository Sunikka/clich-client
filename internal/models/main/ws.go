package mainModel

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sunikka/clich-client/internal/utils"
	"golang.org/x/net/websocket"
)

type WsConnected struct {
	Conn *websocket.Conn
}

type WsErr struct {
	Error string
}

// Websocket connection
// Currently only to a global chat...
func (m *MainModel) establishChatConnection() tea.Cmd {
	return func() tea.Msg {
		err := godotenv.Load()
		if err != nil {
			log.Fatal(err)
		}
		// Websocket connection
		url := os.Getenv("SERVER_URL_WS")
		origin := os.Getenv("CLIENT_ORIGIN")

		ws, err := websocket.Dial(url, "", origin)
		if err != nil {
			return WsErr{Error: fmt.Sprintf("Websocket connection failed: %v", err)}
		}

		m.wsConn = ws
		return WsConnected{Conn: ws}
	}
}

//	userJSON, err := json.Marshal(user)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = websocket.Message.Send(ws, userJSON)
//	if err != nil {
//		log.Fatal(err)
//	}

func (m MainModel) startMsgHandler() tea.Cmd {
	return func() tea.Msg {
		go func() {
			for {
				var msg string
				err := websocket.Message.Receive(m.wsConn, &msg)
				if err != nil {
					log.Println("Error reading the message: ", err)
					return
				}
				m.app.Send(utils.Message{
					SenderID: uuid.New(),
					Username: m.username,
					Content:  msg})
			}
		}()

		return nil
	}
}
