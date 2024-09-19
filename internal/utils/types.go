package utils

import "github.com/google/uuid"

type UserInfo struct {
	ID     string
	Name   string
	Active bool
	// conn *websocket.Conn
	// token jwt.Token
}

type Message struct {
	SenderID uuid.UUID
	Username string
	Content  string
}
