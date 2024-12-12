package auth

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginSuccess struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

type LoginFailure struct {
	Error string `json:"error"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// just going to print something like "User registered succesfully" to the UI
type RegisterSuccess struct {
	Message string
}

type RegisterFailure struct {
	Error string `json:"error"`
}
