package usecase

import (
	"TaskTracker/internal/models"
	"TaskTracker/internal/pkg/services/user"
	"context"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"log/slog"
)

type Params struct {
	fx.In

	Repo   user.Repository
	Logger *slog.Logger
}

type Usecase struct {
	repo user.Repository
	log  *slog.Logger
}

func New(params Params) *Usecase {
	return &Usecase{
		repo: params.Repo,
		log:  params.Logger,
	}
}

func (uc *Usecase) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	userData, err := uc.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func (uc *Usecase) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	userData, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return userData, nil
}
