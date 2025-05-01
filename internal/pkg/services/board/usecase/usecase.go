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

func (uc *Usecase) AddMember(ctx context.Context, memberData *models.BoardMemberAdd, inviterID uuid.UUID) (*models.BoardMember, error) {
	userData, err := uc.repo.GetUserByEmail(ctx, memberData.Email)
	if err != nil {
		return nil, err
	}
	boardMember := &models.BoardMember{
		BoardID: memberData.BoardID,
		UserID:  userData.ID,
		RoleID:  memberData.RoleID,
		IsFav:   false,
	}

	boardMember, err = uc.repo.AddBoardMember(ctx, boardMember, inviterID)
	if err != nil {
		return nil, err
	}
	return boardMember, nil
}

func (uc *Usecase) AddTask(ctx context.Context, task *models.TaskCreate, createdBy uuid.UUID) (*models.TaskCreate, error) {
	task.ID = uuid.New()
	task, err := uc.repo.AddTask(ctx, task, createdBy)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (uc *Usecase) CreateSection(ctx context.Context, section *models.Section, userID uuid.UUID) (*models.Section, error) {
	section.ID = uuid.New()
	section, err := uc.repo.CreateSection(ctx, section, userID)
	if err != nil {
		return nil, err
	}
	return section, nil
}

func (uc *Usecase) GetBoardMemberList(ctx context.Context, boardID uuid.UUID) ([]*models.BoardMemberList, error) {
	return uc.repo.GetBoardMembers(ctx, boardID)
}

func (uc *Usecase) GetAllTasks(ctx context.Context, boardID uuid.UUID) ([]*models.AllTask, error) {
	return uc.repo.GetAllTasks(ctx, boardID)
}

func (uc *Usecase) UpdateTask(ctx context.Context, task *models.UpdateTask, updatedBy uuid.UUID) (*models.UpdateTask, error) {
	return uc.repo.UpdateTask(ctx, task, updatedBy)
}
