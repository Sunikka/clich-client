package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
	mainUI "github.com/sunikka/clich-client/internal/models/main"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// The mainModel needs to reference the app, which is why it has to be initialized this way
	app := tea.NewProgram(mainUI.NewMainModel(nil))
	mainModel := mainUI.NewMainModel(app)

	app = tea.NewProgram(mainModel, tea.WithAltScreen())

	_, err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
