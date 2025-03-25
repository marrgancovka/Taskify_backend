package models

import "github.com/google/uuid"

type TaskStatuses struct {
	TaskID      uuid.UUID `json:"task_id"`
	IsCompleted bool      `json:"is_completed"`
	UpdatedAt   string    `json:"-"`
}
