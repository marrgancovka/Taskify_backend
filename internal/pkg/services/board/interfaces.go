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
	AddTask(ctx context.Context, task *models.TaskCreate, createdBy uuid.UUID) (*models.TaskCreate, error)
	CreateSection(ctx context.Context, section *models.Section, userID uuid.UUID) (*models.Section, error)
	AddMember(ctx context.Context, memberData *models.BoardMemberAdd, inviterID uuid.UUID) (*models.BoardMember, error)
	GetBoardMemberList(ctx context.Context, boardID uuid.UUID) ([]*models.BoardMemberList, error)
	GetAllTasks(ctx context.Context, boardID uuid.UUID) ([]*models.AllTask, error)
	UpdateTask(ctx context.Context, task *models.UpdateTask, updatedBy uuid.UUID) (*models.UpdateTask, error)
}
type Repository interface {
	CreateBoard(ctx context.Context, board *models.Board) (*models.Board, error)
	GetUserListBoards(ctx context.Context, userId uuid.UUID) ([]*models.ListBoards, error)
	SetFavouriteBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error
	SetNoFavouriteBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error
	GetTaskInBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) (*models.TaskInBoard, error)
	IsBoardMember(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) (bool, error)
	IsBoardOwner(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) (bool, error)
	AddBoardMember(ctx context.Context, boardMember *models.BoardMember, inviterID uuid.UUID) (*models.BoardMember, error)
	AddTask(ctx context.Context, task *models.TaskCreate, createdBy uuid.UUID) (*models.TaskCreate, error)
	CreateSection(ctx context.Context, section *models.Section, userID uuid.UUID) (*models.Section, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetBoardMembers(ctx context.Context, boardID uuid.UUID) ([]*models.BoardMemberList, error)
	GetAllTasks(ctx context.Context, boardID uuid.UUID) ([]*models.AllTask, error)
	UpdateTask(ctx context.Context, task *models.UpdateTask, updatedBy uuid.UUID) (*models.UpdateTask, error)
}
