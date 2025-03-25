package repo

import (
	"TaskTracker/internal/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"log/slog"
)

type Params struct {
	fx.In

	DB     *sql.DB
	Logger *slog.Logger
}

type Repository struct {
	db  *sql.DB
	log *slog.Logger
}

func New(params Params) *Repository {
	return &Repository{
		db:  params.DB,
		log: params.Logger,
	}
}

func (repo *Repository) CreateBoard(ctx context.Context, board *models.Board) (*models.Board, error) {
	query := `INSERT INTO default.boards (id, owner_id, name)
		VALUES (?, ?, ?);`

	board.ID = uuid.New()

	_, err := repo.db.ExecContext(ctx, query,
		board.ID,
		board.OwnerID,
		board.Name,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}
	return board, nil
}

func (repo *Repository) GetUserListBoards(ctx context.Context, userId uuid.UUID) ([]*models.Board, error) {
	query := `SELECT b.id, b.name
FROM boards b
WHERE b.owner_id = ?
ORDER BY b.created_at DESC;
`
	rows, err := repo.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query boards: %w", err)
	}
	defer rows.Close()

	var boards []*models.Board
	for rows.Next() {
		board := &models.Board{}
		if err = rows.Scan(&board.ID, &board.Name); err != nil {
			return nil, fmt.Errorf("failed to scan board: %w", err)
		}
		boards = append(boards, board)
	}
	return boards, nil
}
