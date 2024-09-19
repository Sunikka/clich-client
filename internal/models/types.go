package viewTypes

type SessionState int

const (
	LoginView SessionState = iota
	ChatView
)

type SwitchViewMsg struct {
	State    SessionState
	Username string
}
