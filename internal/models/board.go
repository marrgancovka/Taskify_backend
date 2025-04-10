package models

import (
	"github.com/google/uuid"
	"time"
)

type Board struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"-"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"-"`
}

type ListBoards struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	IsFav     bool      `json:"is_fav"`
	Color     string    `json:"color"`
	IsOwner   bool      `json:"is_owner"`
	TaskCount int64     `json:"task_count"`
}
