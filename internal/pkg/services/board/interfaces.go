package board

import (
	"TaskTracker/internal/models"
	"context"
	"github.com/google/uuid"
)

type Usecase interface {
	CreateBoard(ctx context.Context, board *models.Board) (*models.Board, error)
	GetUserListBoards(ctx context.Context, id uuid.UUID) ([]*models.ListBoards, error)
	SetFavouriteBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error
	SetNoFavouriteBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error
	GetTaskInBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) (*models.TaskInBoard, error)
	AddMember(ctx context.Context, memberData *models.BoardMemberAdd) error
}
type Repository interface {
	CreateBoard(ctx context.Context, board *models.Board) (*models.Board, error)
	GetUserListBoards(ctx context.Context, userId uuid.UUID) ([]*models.ListBoards, error)
	SetFavouriteBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error
	SetNoFavouriteBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error
	GetTaskInBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) (*models.TaskInBoard, error)
	IsBoardMember(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) (bool, error)
	IsBoardOwner(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) (bool, error)
	AddMember(ctx context.Context, memberData *models.BoardMemberAdd) error
}
