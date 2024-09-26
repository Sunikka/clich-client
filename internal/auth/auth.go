package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func SendLoginRequestCmd(username, password string) tea.Cmd {
	return func() tea.Msg {
		loginURL := os.Getenv("SERVER_URL") + ":" + os.Getenv("AUTH_PORT") + "/v1/login"

		data := LoginRequest{
			Username: username,
			Password: password,
		}

		reqBody, err := json.Marshal(data)
		if err != nil {
			return LoginFailure{Error: fmt.Sprintf("Login failed: %v", err)}
		}

		res, err := http.Post(loginURL, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			return LoginFailure{Error: fmt.Sprintf("Login failed: %v | %v", err, res.StatusCode)}
		}

		if res.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(res.Body)
			return LoginFailure{Error: fmt.Sprintf("Login failed, status: %d, response: %s", res.StatusCode, body)}
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return LoginFailure{Error: fmt.Sprintf("Login failed: %v | %v", err, res.StatusCode)}
		}

		result := LoginSuccess{}
		err = json.Unmarshal(resBody, &result)
		if err != nil {
			return LoginFailure{Error: fmt.Sprintf("Error parsing login response: %v", err)}
		}

		if result.Token == "" {
			return LoginFailure{Error: "Login failed, Invalid token"}
		}

		return result
	}
}
