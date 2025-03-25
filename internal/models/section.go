package models

import "github.com/google/uuid"

type Section struct {
	ID       uuid.UUID `json:"id"`
	BoardID  uuid.UUID `json:"board_id"`
	Name     string    `json:"name"`
	Position int32     `json:"position"`
}
