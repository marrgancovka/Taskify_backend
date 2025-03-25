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

func (r *Repository) CreateUser(ctx context.Context, data *models.SignUpRequest) (uuid.UUID, error) {
	query := `
		INSERT INTO users (id, username, email, password)
		VALUES (?, ?, ?, ?)
	`

	id := uuid.New()

	_, err := r.db.ExecContext(ctx, query,
		id.String(),
		data.Username,
		data.Email,
		data.Password,
	)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, email, username, password
		FROM users
		WHERE email = ?
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, email)
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
