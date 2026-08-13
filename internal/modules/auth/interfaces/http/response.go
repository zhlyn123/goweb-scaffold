package http

import "time"

type RegisterResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type LoginResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
}
