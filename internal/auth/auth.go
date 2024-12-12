package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

		if res == nil {
			return LoginFailure{Error: "Can't connect to the authentication service"}
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(res.Body)
			return LoginFailure{Error: fmt.Sprintf("Login failed, status: %d, response: %s", res.StatusCode, body)}
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return LoginFailure{Error: fmt.Sprintf("Login failed: %v | %v", err, res.StatusCode)}
		}
		fmt.Print(string(resBody))
		var result LoginSuccess
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

func SendRegisterRequestCmd(username, password string) tea.Cmd {
	return func() tea.Msg {
		RegisterURL := os.Getenv("SERVER_URL") + ":" + os.Getenv("AUTH_PORT") + "/v1/register"
		log.Fatal("RegisterURL: ", RegisterURL)
		data := RegisterRequest{
			Username: username,
			Password: password,
		}

		reqBody, err := json.Marshal(data)
		if err != nil {
			return RegisterFailure{Error: fmt.Sprintf("Register failed: %v", err)}
		}

		res, err := http.Post(RegisterURL, "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			return RegisterFailure{Error: fmt.Sprintf("Register failed: %v | %v", err, res.StatusCode)}
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(res.Body)
			return RegisterFailure{Error: fmt.Sprintf("Register failed, status: %d, response: %s", res.StatusCode, body)}
		}

		// NOTE: The JSON response body is currently not utilized in any way, need to think about it

		// resBody, err := io.ReadAll(res.Body)
		// if err != nil {
		// 	return RegisterFailure{Error: fmt.Sprintf("Register failed: %v | %v", err, res.StatusCode)}
		// }

		// var result RegisterSuccess
		// err = json.Unmarshal(resBody, &result)
		// if err != nil {
		// 	return RegisterFailure{Error: fmt.Sprintf("Error parsing register response: %v", err)}
		// }

		result := RegisterSuccess{Message: fmt.Sprintf("User %s signed up succesfully!", username)}
		return result
	}
}
