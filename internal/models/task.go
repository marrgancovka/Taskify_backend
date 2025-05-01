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
	Percent     int32     `json:"percent"`
	CreatedDate string    `json:"created_date"`
}

type TaskInBoard struct {
	Board    *ListBoards        `json:"board"`
	Sections []*SectionWithTask `json:"section"`
}

type TaskCreate struct {
	ID          uuid.UUID `json:"id"`
	BoardID     uuid.UUID `json:"board_id"`
	SectionID   uuid.UUID `json:"section_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"`
	Priority    string    `json:"priority"`
	Percent     string    `json:"percent"`
	AssigneeID  uuid.UUID `json:"assignee_id"`
	CreatedDate string    `json:"created_date"`
}

type TaskData struct {
	ID               uuid.UUID `json:"id"`
	BoardID          uuid.UUID `json:"board_id"`
	SectionID        uuid.UUID `json:"section_id"`
	Name             string    `json:"name"`
	Description      string    `json:"description"`
	DueDate          string    `json:"due_date"`
	Priority         string    `json:"priority"`
	Percent          string    `json:"percent"`
	AssigneeID       uuid.UUID `json:"assignee_id"`
	AssigneeEmail    string    `json:"assignee_email"`
	AssigneeUsername string    `json:"assignee_username"`
	CreatedDate      string    `json:"created_date"`
}

type AllTask struct {
	ID                uuid.UUID   `json:"id"`
	Name              string      `json:"name"`
	DueDate           string      `json:"due_date"`
	Priority          string      `json:"priority"`
	CreatedDate       string      `json:"created_date"`
	AssigneeID        []uuid.UUID `json:"assignee_id"`
	TaskDependencies  []uuid.UUID `json:"task_dependencies"`
	DonePercent       string      `json:"done_percent"`
	PredictionDueDate string      `json:"prediction_due_date"`
}

type UpdateTask struct {
	ID          uuid.UUID `json:"id"`
	SectionID   uuid.UUID `json:"section_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     string    `json:"due_date"`
	Priority    string    `json:"priority"`
	AssigneeID  uuid.UUID `json:"assignee_id"`
	Percent     string    `json:"percent"`
}
