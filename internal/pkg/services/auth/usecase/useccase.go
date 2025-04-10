package usecase

import (
	"TaskTracker/internal/models"
	"TaskTracker/internal/pkg/services/auth"
	"TaskTracker/internal/pkg/tokenizer"
	"TaskTracker/pkg/hasher"
	"context"
	"errors"
	"go.uber.org/fx"
	"log/slog"
)

type Params struct {
	fx.In

	Repo      auth.Repository
	Tokenizer *tokenizer.Tokenizer
	Logger    *slog.Logger
}

type Usecase struct {
	repo      auth.Repository
	tokenizer *tokenizer.Tokenizer
	log       *slog.Logger
}

func New(params Params) *Usecase {
	return &Usecase{
		repo:      params.Repo,
		log:       params.Logger,
		tokenizer: params.Tokenizer,
	}
}

func (uc *Usecase) CreateUser(ctx context.Context, data *models.SignUpRequest) (*models.TokenResponse, error) {
	data.Password = hasher.GenerateHashString(data.Password)
	id, err := uc.repo.CreateUser(ctx, data)
	if err != nil {
		return nil, err
	}
	token, err := uc.tokenizer.GenerateJWT(&models.TokenPayload{UserID: id})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (uc *Usecase) Login(ctx context.Context, data *models.LoginRequest) (*models.TokenResponse, error) {
	user, err := uc.repo.GetUserByEmail(ctx, data.Email)
	if err != nil {
		return nil, err
	}
	if !hasher.CompareStringHash(data.Password, user.Password) {
		return nil, errors.New("invalid password")
	}
	token, err := uc.tokenizer.GenerateJWT(&models.TokenPayload{UserID: user.ID})
	if err != nil {
		return nil, err
	}

	return token, nil
}
