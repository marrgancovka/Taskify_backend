package user

import (
	"TaskTracker/internal/models"
	"context"
	"github.com/google/uuid"
)

type Usecase interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}
type Repository interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}
