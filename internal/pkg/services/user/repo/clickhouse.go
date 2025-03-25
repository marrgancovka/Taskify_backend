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

func (repo *Repository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, username, password
		FROM users
		WHERE id = ?
		LIMIT 1
	`
	row := repo.db.QueryRowContext(ctx, query, id)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (repo *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, username, password
		FROM users
		WHERE email = ?
		LIMIT 1
	`
	row := repo.db.QueryRowContext(ctx, query, email)
	user := &models.User{}
	err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}
