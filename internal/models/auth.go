package models

import (
	"github.com/google/uuid"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenPayload struct {
	UserID uuid.UUID `json:"user_id"`
	Exp    time.Time `json:"exp"`
}

type TokenResponse struct {
	Token string    `json:"token"`
	Exp   time.Time `json:"exp"`
}
