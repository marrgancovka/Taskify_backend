package repo

import (
	"TaskTracker/internal/models"
	"context"
	"database/sql"
	"errors"
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
	queryCreateBoard := `INSERT INTO default.boards (id, owner_id, name, color)
		VALUES (?, ?, ?, ?);`
	queryCreateBoardMembers := `INSERT INTO default.board_members (board_id, user_id) VALUES (?, ?);`

	board.ID = uuid.New()
	_, err := repo.db.ExecContext(ctx, queryCreateBoard,
		board.ID,
		board.OwnerID,
		board.Name,
		board.Color,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}

	_, err = repo.db.ExecContext(ctx, queryCreateBoardMembers,
		board.ID,
		board.OwnerID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create board members: %w", err)
	}

	return board, nil
}

func (repo *Repository) GetUserListBoards(ctx context.Context, userId uuid.UUID) ([]*models.ListBoards, error) {
	query := `SELECT 
    b.id, 
    b.name, 
    b.color,
    COALESCE(bm.is_favourite, false) AS is_favourite,
    b.owner_id = ? AS is_owner
FROM 
    boards b
JOIN 
    board_members bm ON b.id = bm.board_id
WHERE 
    bm.user_id = ?
ORDER BY 
    b.name DESC;`
	countQuery := `select count(t.id) from tasks t where t.board_id=?`
	rows, err := repo.db.QueryContext(ctx, query, userId, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query boards: %w", err)
	}
	defer rows.Close()

	var boards []*models.ListBoards
	for rows.Next() {
		board := &models.ListBoards{}
		if err = rows.Scan(&board.ID, &board.Name, &board.Color, &board.IsFav, &board.IsOwner); err != nil {
			return nil, fmt.Errorf("failed to scan board: %w", err)
		}
		row := repo.db.QueryRowContext(ctx, countQuery, board.ID)
		err = row.Scan(&board.TaskCount)
		if err != nil {
			return nil, fmt.Errorf("failed to get task count: %w", err)
		}
		boards = append(boards, board)
	}
	return boards, nil
}

func (repo *Repository) SetFavouriteBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error {
	query := `ALTER TABLE default.board_members 
              UPDATE is_favourite = true 
              WHERE board_id = ? AND user_id = ?`

	_, err := repo.db.ExecContext(ctx, query, boardID, userID)
	if err != nil {
		return fmt.Errorf("failed to set favourite board: %w", err)
	}
	return nil
}
func (repo *Repository) SetNoFavouriteBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) error {
	query := `ALTER TABLE default.board_members 
              UPDATE is_favourite = false 
              WHERE board_id = ? AND user_id = ?`

	_, err := repo.db.ExecContext(ctx, query, boardID, userID)
	if err != nil {
		return fmt.Errorf("failed to set no-favourite board: %w", err)
	}
	return nil
}

func (repo *Repository) IsBoardMember(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 
            FROM default.board_members 
            WHERE board_id = ? AND user_id = ?
        )`

	var exists uint8 // ClickHouse возвращает EXISTS как UInt8 (0 или 1)
	err := repo.db.QueryRowContext(ctx, query, boardID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check board membership: %w", err)
	}

	return exists == 1, nil
}

func (repo *Repository) IsBoardOwner(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) (bool, error) {
	query := `
        SELECT EXISTS(
            SELECT 1 
            FROM default.boards 
            WHERE id = ? AND owner_id = ?
        )`

	var exists uint8 // ClickHouse возвращает EXISTS как UInt8 (0 или 1)
	err := repo.db.QueryRowContext(ctx, query, boardID, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check board owner: %w", err)
	}

	return exists == 1, nil
}

func (repo *Repository) GetTaskInBoard(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) (*models.TaskInBoard, error) {
	query := `
        WITH board_data AS (
            SELECT b.id, b.name, b.color, COALESCE(bm.is_favourite, false) AS is_favourite, b.owner_id = ? AS is_owner
			FROM default.boards b
			LEFT JOIN board_members bm ON b.id = bm.board_id AND bm.user_id = ?
			WHERE b.id = ?
        ),
        sections_data AS (
            SELECT id, board_id, name, position
            FROM default.sections
            WHERE board_id = ?
            ORDER BY position ASC
        ),
        tasks_data AS (
            SELECT 
                t.id, t.board_id, t.section_id, t.name, 
                t.description, t.due_date, t.priority, t.created_at
            FROM default.tasks t
            WHERE t.board_id = ?
            ORDER BY t.created_at DESC
        )
        SELECT 
            b.id, b.name, b.color, b.is_favourite, b.is_owner,
            s.id, s.board_id, s.name, s.position,
            t.id, t.board_id, t.section_id, t.name, 
            t.description, t.due_date, t.priority, t.created_at
        FROM board_data b
        LEFT JOIN sections_data s ON s.board_id = b.id
        LEFT JOIN tasks_data t ON t.section_id = s.id
    `

	rows, err := repo.db.QueryContext(ctx, query, userID, userID, boardID, boardID, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to query board data: %w", err)
	}
	defer rows.Close()

	result := &models.TaskInBoard{
		Board:    &models.ListBoards{},
		Sections: []*models.SectionWithTask{},
	}
	sectionsMap := make(map[uuid.UUID]*models.SectionWithTask)

	for rows.Next() {
		board := &models.ListBoards{}
		section := &models.SectionWithTask{}
		task := &models.Task{}

		err = rows.Scan(
			&board.ID, &board.Name, &board.Color, &board.IsFav, &board.IsOwner,
			&section.ID, &section.BoardID, &section.Name, &section.Position,
			&task.ID, &task.BoardID, &task.SectionID, &task.Name,
			&task.Description, &task.DueDate, &task.Priority, &task.CreatedDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Если это первая итерация, сохраняем данные доски
		if result.Board.ID == uuid.Nil {
			result.Board = board
		}

		// Обрабатываем секции
		if _, exists := sectionsMap[section.ID]; !exists && section.ID != uuid.Nil {
			sectionsMap[section.ID] = section
			result.Sections = append(result.Sections, section)
		}

		// Добавляем задачу в соответствующую секцию
		if task.ID != uuid.Nil {
			for i, sec := range result.Sections {
				if sec.ID == task.SectionID {
					result.Sections[i].Tasks = append(sec.Tasks, task)
					break
				}
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	// Если доска не найдена
	if result.Board.ID == uuid.Nil {
		return nil, fmt.Errorf("board not found")
	}

	countQuery := `select count(t.id) from tasks t where t.board_id=?`
	row := repo.db.QueryRowContext(ctx, countQuery, result.Board.ID)
	err = row.Scan(&result.Board.TaskCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get task count: %w", err)
	}

	return result, nil
}

func (repo *Repository) AddMember(ctx context.Context, memberData *models.BoardMemberAdd) error {
	queryUser := `
        SELECT 
            id
        FROM 
            default.users 
        WHERE 
            email = ?
        LIMIT 1`

	var userID uuid.UUID
	err := repo.db.QueryRowContext(ctx, queryUser, memberData.Email).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("пользователь не найден") // Пользователь не найден
		}
		return fmt.Errorf("failed to get user ID: %w", err)
	}

	memberData.UserID = userID

	query := `INSERT INTO default.board_members (board_id, user_id, role_id) VALUES(?, ?, ?)`
	_, err = repo.db.ExecContext(ctx, query, memberData.BoardID, memberData.UserID, memberData.RoleID)
	if err != nil {
		return fmt.Errorf("failed to add member to board: %w", err)
	}
	return nil
}
