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
	"strings"
	"time"
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
	repo.log.Debug("insert board", "boardData", board)

	_, err = repo.db.ExecContext(ctx, queryCreateBoardMembers,
		board.ID,
		board.OwnerID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create board members: %w", err)
	}
	repo.log.Debug("insert board_members")

	section := &models.Section{
		ID:       uuid.New(),
		BoardID:  board.ID,
		Name:     "Без раздела",
		Position: 0,
	}
	_, err = repo.db.ExecContext(ctx, `
		INSERT INTO sections (
			id, board_id, name, position
		) VALUES (?, ?, ?, ?)`,
		section.ID, section.BoardID, section.Name, section.Position)
	if err != nil {
		repo.log.Error("failed to insert section: %w", err)
		return nil, fmt.Errorf("failed to insert section: %w", err)
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
        t.description, t.due_date, t.priority, t.created_at, t.percent,
        ta.user_id as assignee_id,
        u.username as assignee_name,
        u.email as assignee_email
    FROM default.tasks t
    LEFT JOIN task_assignees ta ON t.id = ta.task_id
    LEFT JOIN users u ON ta.user_id = u.id
    WHERE t.board_id = ?
)
        SELECT 
            b.id, b.name, b.color, b.is_favourite, b.is_owner,
            s.id, s.board_id, s.name, s.position,
            t.id, t.board_id, t.section_id, t.name, 
            t.description, t.due_date, t.priority, t.created_at, t.percent,
            t.assignee_id, t.assignee_name, t.assignee_email
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
		task := &models.TaskData{}
		//var assigneeID uuid.NullUUID // Используем NullUUID для обработки NULL значений

		err = rows.Scan(
			&board.ID, &board.Name, &board.Color, &board.IsFav, &board.IsOwner,
			&section.ID, &section.BoardID, &section.Name, &section.Position,
			&task.ID, &task.BoardID, &task.SectionID, &task.Name,
			&task.Description, &task.DueDate, &task.Priority, &task.CreatedDate, &task.Percent,
			&task.AssigneeID, &task.AssigneeUsername, &task.AssigneeEmail,
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

	countQuery := `SELECT COUNT(t.id) FROM tasks t WHERE t.board_id = ?`
	row := repo.db.QueryRowContext(ctx, countQuery, result.Board.ID)
	err = row.Scan(&result.Board.TaskCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get task count: %w", err)
	}

	return result, nil
}

// AddTask создает новую задачу в ClickHouse
func (repo *Repository) AddTask(ctx context.Context, task *models.TaskCreate, createdBy uuid.UUID) (*models.TaskCreate, error) {
	// 1. Проверяем существование доски и секции
	if err := repo.validateBoardAndSection(ctx, task.BoardID, task.SectionID, createdBy); err != nil {
		return nil, err
	}

	createdTask, err := repo.insertTask(ctx, task)
	if err != nil {
		return nil, err
	}

	if err = repo.addAssignees(ctx, task.ID, task.AssigneeID); err != nil {
		return nil, err
	}

	if err = repo.addTaskHistory(ctx, task.ID, createdBy, "create", "", "", ""); err != nil {
		return nil, err
	}

	if err = repo.createTaskStatus(ctx, task.ID); err != nil {
		return nil, err
	}

	return createdTask, nil
}

func (repo *Repository) validateBoardAndSection(ctx context.Context, boardID, sectionID, userID uuid.UUID) error {
	var exists bool
	err := repo.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM board_members 
			WHERE board_id = ? AND user_id = ?
		)`, boardID, userID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check board access: %w", err)
	}
	if !exists {
		return errors.New("user doesn't have access to this board")
	}

	// Проверяем, что секция принадлежит доске
	err = repo.db.QueryRowContext(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM sections 
			WHERE id = ? AND board_id = ?
		)`, sectionID, boardID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to validate section: %w", err)
	}
	if !exists {
		return errors.New("section doesn't belong to the board")
	}

	return nil
}

func (repo *Repository) insertTask(ctx context.Context, task *models.TaskCreate) (*models.TaskCreate, error) {
	_, err := repo.db.ExecContext(ctx, `
            INSERT INTO tasks (
                id, board_id, section_id, name, 
                description, due_date, priority, created_at
            ) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		task.ID, task.BoardID, task.SectionID, task.Name,
		task.Description, task.DueDate, task.Priority, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to insert task: %w", err)
	}
	return task, nil
}

func (repo *Repository) addAssignees(ctx context.Context, taskID uuid.UUID, assigneeID uuid.UUID) error {
	// Проверка на нулевой UUID
	if assigneeID == uuid.Nil {
		return nil
	}

	// Используем ExecContext напрямую без Prepare для ClickHouse
	// (подготовленные запросы могут работать неоптимально в ClickHouse)
	_, err := repo.db.ExecContext(ctx, `
        INSERT INTO task_assignees (task_id, user_id) 
        VALUES (?, ?)`,
		taskID.String(), // Явное преобразование UUID в строку
		assigneeID.String(),
	)

	if err != nil {
		// Более информативное сообщение об ошибке
		return fmt.Errorf("failed to add assignee %s to task %s: %w",
			assigneeID, taskID, err)
	}

	return nil
}

func (repo *Repository) addTaskHistory(ctx context.Context, taskID, userID uuid.UUID, action, fieldName, oldValue, newValue string) error {
	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO default.task_history (
			id, task_id, user_id, action, 
			field_name, old_value, new_value, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		uuid.New(), taskID, userID, action,
		fieldName, oldValue, newValue, time.Now())
	if err != nil {
		return fmt.Errorf("failed to insert task history: %w", err)
	}
	return nil
}

func (repo *Repository) createTaskStatus(ctx context.Context, taskID uuid.UUID) error {
	_, err := repo.db.ExecContext(ctx, `
		INSERT INTO default.task_statuses (task_id, is_completed, updated_at) 
		VALUES (?, ?, ?)`,
		taskID.String(), false, time.Now())
	if err != nil {
		return fmt.Errorf("failed to create task status: %w", err)
	}
	return nil
}

func (repo *Repository) CreateSection(ctx context.Context, section *models.Section, userID uuid.UUID) (*models.Section, error) {
	// 1. Проверяем, что пользователь имеет доступ к доске
	isMember, err := repo.IsBoardMember(ctx, section.BoardID, userID)
	if err != nil {
		repo.log.Error("failed to check if board member is a member: %w", "error", err)
		return nil, err
	}
	if !isMember {
		repo.log.Error("failed to check if board member is a member: %w", err)
		return nil, err
	}

	var maxPosition int32
	err = repo.db.QueryRowContext(ctx, `
		SELECT COALESCE(MAX(position), 0) 
		FROM sections 
		WHERE board_id = ?`, section.BoardID).Scan(&maxPosition)
	if err != nil {
		repo.log.Error("failed to get max section position: %w", err)
		return nil, fmt.Errorf("failed to get max section position: %w", err)
	}

	section.Position = maxPosition + 1

	// 4. Вставляем секцию в базу данных
	_, err = repo.db.ExecContext(ctx, `
		INSERT INTO sections (
			id, board_id, name, position
		) VALUES (?, ?, ?, ?)`,
		section.ID, section.BoardID, section.Name, section.Position, time.Now())
	if err != nil {
		repo.log.Error("failed to insert section: %w", err)
		return nil, fmt.Errorf("failed to insert section: %w", err)
	}

	return section, nil
}

// AddBoardMember добавляет участника в доску
func (repo *Repository) AddBoardMember(ctx context.Context, boardMember *models.BoardMember, inviterID uuid.UUID) (*models.BoardMember, error) {
	// 1. Проверяем, что приглашающий имеет права на добавление участников
	isOwner, err := repo.IsBoardOwner(ctx, boardMember.BoardID, inviterID)
	if err != nil {
		repo.log.Error("failed to check if board member is a owner: %w", err)
		return nil, err
	}
	if !isOwner {
		repo.log.Error("failed to check if board member is a owner: %w", err)
		return nil, err
	}

	// 2. Проверяем, что пользователь еще не является участником доски
	alreadyMember, err := repo.IsBoardMember(ctx, boardMember.BoardID, boardMember.UserID)
	if err != nil {
		repo.log.Error("failed to check if board member is a owner: %w", err)
		return nil, err
	}
	if alreadyMember {
		repo.log.Error("board member already exists")
		return nil, errors.New("board member already exists")
	}

	// 3. Проверяем существование роли
	//var roleExists bool
	//err = r.db.QueryRowContext(ctx, `
	//	SELECT EXISTS(
	//		SELECT 1 FROM board_roles
	//		WHERE id = $1 AND board_id = $2
	//	)`, roleID, boardID).Scan(&roleExists)
	//if err != nil {
	//	return fmt.Errorf("failed to check role existence: %w", err)
	//}
	//if !roleExists {
	//	return errors.New("specified role does not exist for this board")
	//}

	// 4. Добавляем участника
	_, err = repo.db.ExecContext(ctx, `
		INSERT INTO board_members (
			board_id, user_id, role_id, is_favourite
		) VALUES ($1, $2, $3, $4)`,
		boardMember.BoardID, boardMember.UserID, boardMember.RoleID, boardMember.IsFav)
	if err != nil {
		return nil, fmt.Errorf("failed to insert board member: %w", err)
	}

	return boardMember, nil
}

// GetUserByEmail возвращает пользователя по email
func (repo *Repository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := repo.db.QueryRowContext(ctx, `
		SELECT id, username, email
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// GetBoardMembers возвращает список участников доски с их ролями и информацией
func (repo *Repository) GetBoardMembers(ctx context.Context, boardID uuid.UUID) ([]*models.BoardMemberList, error) {
	query := `
		SELECT 
			u.id,
			u.username,
			u.email
		FROM board_members bm
		JOIN users u ON bm.user_id = u.id
		WHERE bm.board_id = $1
		ORDER BY u.username
	`

	rows, err := repo.db.QueryContext(ctx, query, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to query board members: %w", err)
	}
	defer rows.Close()

	var members []*models.BoardMemberList
	for rows.Next() {
		member := &models.BoardMemberList{}
		err := rows.Scan(
			&member.UserID,
			&member.Username,
			&member.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan board member: %w", err)
		}
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating board members: %w", err)
	}

	return members, nil
}

func (repo *Repository) GetAllTasks(ctx context.Context, boardID uuid.UUID) ([]*models.AllTask, error) {
	query := `SELECT 
			id,
			name,
			due_date,
			priority,
			created_at,
			percent
		FROM default.tasks
		WHERE board_id = ?
		ORDER BY due_date DESC`

	rows, err := repo.db.QueryContext(ctx, query, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to query all tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.AllTask
	for rows.Next() {
		task := &models.AllTask{}
		err = rows.Scan(
			&task.ID,
			&task.Name,
			&task.DueDate,
			&task.Priority,
			&task.CreatedDate,
			&task.DonePercent)
		if err != nil {
			return nil, fmt.Errorf("failed to scan all tasks: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating all tasks: %w", err)
	}
	return tasks, nil

}

func (repo *Repository) UpdateTask(ctx context.Context, task *models.UpdateTask, updatedBy uuid.UUID) (*models.UpdateTask, error) {
	// Получаем текущее состояние задачи
	existingTask := &models.UpdateTask{}
	query := `
		SELECT id, section_id, name, description, due_date, priority, percent
		FROM default.tasks
		WHERE id = ?
	`
	err := repo.db.QueryRowContext(ctx, query, task.ID).Scan(
		&existingTask.ID,
		&existingTask.SectionID,
		&existingTask.Name,
		&existingTask.Description,
		&existingTask.DueDate,
		&existingTask.Priority,
		&existingTask.Percent,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch existing task: %w", err)
	}

	updateFields := []string{}
	updateArgs := []interface{}{}
	historyInserts := []string{}
	historyArgs := []interface{}{}

	// Обработка секции отдельно, чтобы подтянуть названия
	if task.SectionID != existingTask.SectionID {
		oldSectionName, err := repo.getSectionNameByID(ctx, existingTask.SectionID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch old section name: %w", err)
		}
		newSectionName, err := repo.getSectionNameByID(ctx, task.SectionID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch new section name: %w", err)
		}

		updateFields = append(updateFields, "section_id = ?")
		updateArgs = append(updateArgs, task.SectionID)

		historyInserts = append(historyInserts, "(generateUUIDv4(), ?, ?, 'update', 'section', ?, ?, now())")
		historyArgs = append(historyArgs, task.ID, updatedBy, oldSectionName, newSectionName)
	}

	// Мапа остальных полей
	updates := map[string]struct {
		NewValue interface{}
		OldValue interface{}
	}{
		"name":        {NewValue: task.Name, OldValue: existingTask.Name},
		"description": {NewValue: task.Description, OldValue: existingTask.Description},
		"due_date":    {NewValue: task.DueDate, OldValue: existingTask.DueDate},
		"priority":    {NewValue: task.Priority, OldValue: existingTask.Priority},
		"percent":     {NewValue: fmt.Sprintf("%v", task.Percent), OldValue: fmt.Sprintf("%v", existingTask.Percent)},
	}

	for field, values := range updates {
		if values.NewValue != values.OldValue {
			updateFields = append(updateFields, fmt.Sprintf("%v = ?", field))
			updateArgs = append(updateArgs, values.NewValue)

			historyInserts = append(historyInserts, "(generateUUIDv4(), ?, ?, 'update', ?, ?, ?, now())")
			historyArgs = append(historyArgs, task.ID, updatedBy, field, values.OldValue, values.NewValue)
		}
	}

	if len(updateFields) == 0 {
		return existingTask, nil
	}

	updateArgs = append(updateArgs, task.ID)

	// Выполняем UPDATE
	updateQuery := fmt.Sprintf(`
		ALTER TABLE default.tasks UPDATE %v WHERE id = ?
	`, strings.Join(updateFields, ", "))

	if _, err = repo.db.ExecContext(ctx, updateQuery, updateArgs...); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	// Добавляем записи в историю изменений
	if len(historyInserts) > 0 {
		historyQuery := fmt.Sprintf(`
			INSERT INTO default.task_history
			(id, task_id, user_id, action, field_name, old_value, new_value, created_at)
			VALUES %s
		`, strings.Join(historyInserts, ", "))

		if _, err = repo.db.ExecContext(ctx, historyQuery, historyArgs...); err != nil {
			return nil, fmt.Errorf("failed to insert task history: %w", err)
		}
	}

	return task, nil
}

func (repo *Repository) getSectionNameByID(ctx context.Context, sectionID uuid.UUID) (string, error) {
	var name string
	query := `
		SELECT name
		FROM default.sections
		WHERE id = ?
		LIMIT 1
	`
	err := repo.db.QueryRowContext(ctx, query, sectionID).Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}
