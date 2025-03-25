package auth

import (
	"TaskTracker/internal/models"
	"context"
	"github.com/google/uuid"
)

type Usecase interface {
	CreateUser(ctx context.Context, data *models.SignUpRequest) (*models.TokenResponse, error)
	Login(ctx context.Context, data *models.LoginRequest) (*models.TokenResponse, error)
}

type Repository interface {
	CreateUser(ctx context.Context, data *models.SignUpRequest) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}
