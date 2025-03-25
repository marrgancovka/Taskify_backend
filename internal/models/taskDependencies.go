package models

import "github.com/google/uuid"

type TaskDependencies struct {
	ParentTaskID uuid.UUID `json:"parent_task_id"`
	ChildTaskID  uuid.UUID `json:"child_task_id"`
}
