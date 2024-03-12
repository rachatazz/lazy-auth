package model

import "time"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	IsAll        bool   `json:"is_all"`
}

type TokenResponse struct {
	TokenType             string    `json:"token_type"`
	AccessToken           string    `json:"access_token"`
	TokenExpiresAt        time.Time `json:"token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}
