package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	mainUI "github.com/sunikka/clich-client/internal/models/main"
	"github.com/sunikka/clich-client/internal/utils"
	"golang.org/x/net/websocket"
)

var testNames = []string{"Paavo", "Gigachad", "Harold", "Taavi", "MogBot"}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	// Websocket connection
	url := os.Getenv("SERVER_URL")
	origin := os.Getenv("CLIENT_ORIGIN")

	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	app := tea.NewProgram(mainUI.NewMainModel(ws), tea.WithAltScreen())

	// Sends client info as a JSON for now
	// To be switched to just a JWT token for authenticating and fetching the user info on server-side
	// when the login system and db has been implemented...
	user := utils.UserInfo{
		ID:     uuid.NewString(),
		Name:   testNames[rand.IntN(len(testNames)-1)],
		Active: true,
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
	}

	err = websocket.Message.Send(ws, userJSON)
	if err != nil {
		log.Fatal(err)
	}

	// Message handler
	go func() {
		var msg string
		for {
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				log.Println("Error reading the message: ", err)
				return
			}
			app.Send(utils.Message{
				SenderID: uuid.New(),
				Username: user.Name,
				Content:  msg})
		}
	}()

	_, err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
