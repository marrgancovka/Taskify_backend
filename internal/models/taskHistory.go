package models

import (
	"github.com/google/uuid"
	"time"
)

type TaskHistory struct {
	ID        uuid.UUID `json:"id"`
	TaskID    uuid.UUID `json:"task_id"`
	UserID    uuid.UUID `json:"user_id"`
	Action    string    `json:"action"`
	FieldName string    `json:"field_name"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	CreatedAt time.Time `json:"created_at"`
}
