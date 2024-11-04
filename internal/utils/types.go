package utils

import "time"

type UserInfo struct {
	ID     string
	Name   string
	Active bool
	// conn *websocket.Conn
	// token jwt.Token
}

type Message struct {
	// SenderID uuid.UUID
	Username string
	Content  string
	SentAt   time.Time
}
