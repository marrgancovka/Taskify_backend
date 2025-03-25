package board

import (
	"TaskTracker/internal/models"
	"context"
	"github.com/google/uuid"
)

type Usecase interface {
	CreateBoard(ctx context.Context, board *models.Board) (*models.Board, error)
	GetUserListBoards(ctx context.Context, id uuid.UUID) ([]*models.Board, error)
}
type Repository interface {
	CreateBoard(ctx context.Context, board *models.Board) (*models.Board, error)
	GetUserListBoards(ctx context.Context, userId uuid.UUID) ([]*models.Board, error)
}
