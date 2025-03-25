package models

import "github.com/google/uuid"

type TaskAssignee struct {
	TaskID uuid.UUID `json:"task_id"`
	UserID uuid.UUID `json:"user_id"`
}
