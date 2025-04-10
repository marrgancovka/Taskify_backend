package models

import "github.com/google/uuid"

type BoardMember struct {
	BoardID uuid.UUID `json:"board_id"`
	UserID  uuid.UUID `json:"user_id"`
	RoleID  uuid.UUID `json:"role_id"`
}
type BoardMemberAdd struct {
	Email   string    `json:"email"`
	BoardID uuid.UUID `json:"board_id"`
	UserID  uuid.UUID `json:"user_id"`
	RoleID  uuid.UUID `json:"role_id"`
}
