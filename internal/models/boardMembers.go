package models

import "github.com/google/uuid"

type BoardMember struct {
	BoardID uuid.UUID `json:"board_id"`
	UserID  uuid.UUID `json:"user_id"`
	RoleID  uuid.UUID `json:"role_id"`
	IsFav   bool      `json:"is_fav"`
}
type BoardMemberAdd struct {
	Email   string    `json:"email"`
	BoardID uuid.UUID `json:"board_id"`
	RoleID  uuid.UUID `json:"role_id"`
}

type BoardMemberList struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	//Role     string    `json:"role"`
}
