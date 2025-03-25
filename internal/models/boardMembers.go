package models

import "github.com/google/uuid"

type BoardMember struct {
	BoardID uuid.UUID `json:"board_id"`
	UserID  uuid.UUID `json:"user_id"`
	RoleID  uuid.UUID `json:"role_id"`
}
