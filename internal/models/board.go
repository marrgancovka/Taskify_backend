package models

import (
	"github.com/google/uuid"
	"time"
)

type Board struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
