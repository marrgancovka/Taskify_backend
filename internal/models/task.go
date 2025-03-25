package models

import "github.com/google/uuid"

type Task struct {
	ID          uuid.UUID `json:"id"`
	BoardID     uuid.UUID `json:"board_id"`
	SectionID   uuid.UUID `json:"section_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"`
	Priority    string    `json:"priority"`
	CreatedDate string    `json:"created_date"`
}
