package auth

type LoginCredentials struct {
	Username string
	Password string
	Role     string
}

type TokenClaims struct {
	UserID string
	Role   string
	Expiry int64
}
type LoginResult struct {
	Token    string `json:"token"`
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
