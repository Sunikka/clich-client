package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	mainUI "github.com/sunikka/clich-client/internal/models/main"
)

var testNames = []string{"Paavo", "Gigachad", "Harold", "Taavi", "MogBot"}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	app := tea.NewProgram(mainUI.NewMainModel(), tea.WithAltScreen())

	// Sends client info as a JSON for now
	// To be switched to just a JWT token for authenticating and fetching the user info on server-side
	// when the login system and db has been implemented...

	_, err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
