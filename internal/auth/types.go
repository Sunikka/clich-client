package auth

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginSuccess struct {
	Username string `json:"username"`
	Token    string
}

type LoginFailure struct {
	Error string `json:"error"`
}
