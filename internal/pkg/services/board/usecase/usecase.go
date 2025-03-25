package usecase

import (
	"TaskTracker/internal/models"
	"TaskTracker/internal/pkg/services/board"
	"context"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"log/slog"
)

type Params struct {
	fx.In

	Repo   board.Repository
	Logger *slog.Logger
}

type Usecase struct {
	repo board.Repository
	log  *slog.Logger
}

func New(params Params) *Usecase {
	return &Usecase{
		repo: params.Repo,
		log:  params.Logger,
	}
}

func (uc *Usecase) CreateBoard(ctx context.Context, board *models.Board) (*models.Board, error) {
	newBoard, err := uc.repo.CreateBoard(ctx, board)
	if err != nil {
		return nil, err
	}
	return newBoard, nil
}
func (uc *Usecase) GetUserListBoards(ctx context.Context, id uuid.UUID) ([]*models.Board, error) {
	boardList, err := uc.repo.GetUserListBoards(ctx, id)
	if err != nil {
		return nil, err
	}
	return boardList, nil
}
