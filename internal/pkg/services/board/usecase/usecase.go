package usecase

import (
	"TaskTracker/internal/models"
	"TaskTracker/internal/pkg/services/board"
	"context"
	"errors"
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
func (uc *Usecase) GetUserListBoards(ctx context.Context, id uuid.UUID) ([]*models.ListBoards, error) {
	boardList, err := uc.repo.GetUserListBoards(ctx, id)
	if err != nil {
		return nil, err
	}
	return boardList, nil
}

func (uc *Usecase) SetFavouriteBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error {
	return uc.repo.SetFavouriteBoard(ctx, boardID, userID)
}

func (uc *Usecase) SetNoFavouriteBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error {
	return uc.repo.SetNoFavouriteBoard(ctx, boardID, userID)
}

func (uc *Usecase) GetTaskInBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) (*models.TaskInBoard, error) {
	isMember, err := uc.repo.IsBoardMember(ctx, boardID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("not a member")
	}

	taskInBoard, err := uc.repo.GetTaskInBoard(ctx, boardID, userID)
	if err != nil {
		return nil, err
	}
	return taskInBoard, nil
}

func (uc *Usecase) AddMember(ctx context.Context, memberData *models.BoardMemberAdd) error {
	return uc.repo.AddMember(ctx, memberData)
}
