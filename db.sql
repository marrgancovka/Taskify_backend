CREATE TABLE default.users
(
    id UUID,
    username String,
    email String,
    password String,
    created_at DateTime DEFAULT now()
)
    ENGINE = MergeTree()
        ORDER BY id;

CREATE TABLE default.boards
(
    id UUID,
    owner_id UUID,
    name String,
    color String,
    created_at DateTime DEFAULT now()
)
    ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE default.board_roles
(
    id UUID,
    board_id UUID,
    name String
)
    ENGINE = MergeTree()
ORDER BY id;

CREATE TABLE default.board_members
(
    board_id UUID,
    user_id UUID,
    role_id UUID,
    is_favourite Bool,
)
    ENGINE = ReplacingMergeTree()
ORDER BY (board_id, user_id);

CREATE TABLE default.sections
(
    id UUID,
    board_id UUID,
    name String,
    position Int32
)
    ENGINE = MergeTree()
ORDER BY (board_id, position);

CREATE TABLE default.tasks
(
    id UUID,
    board_id UUID,
    section_id UUID,
    name String,
    description String,
    due_date DateTime,
    priority Int32,
    percent Int32,
    created_at DateTime DEFAULT now()
)
    ENGINE = MergeTree()
ORDER BY (board_id, section_id, due_date, priority);

CREATE TABLE default.task_assignees
(
    task_id UUID,
    user_id UUID
)
    ENGINE = MergeTree()
ORDER BY (task_id, user_id);

CREATE TABLE default.task_statuses
(
    task_id UUID,
    is_completed Bool,
    updated_at DateTime DEFAULT now()
)
    ENGINE = ReplacingMergeTree()
ORDER BY task_id;

CREATE TABLE default.task_history
(
    id UUID,
    task_id UUID,
    user_id UUID,
    action String, -- "create", "update", "status_change", "move", "assign", "dependency"
    field_name String, -- например, "name", "section", "status", "due_date"
    old_value String, -- старое значение в текстовом формате
    new_value String, -- новое значение в текстовом формате
    created_at DateTime DEFAULT now()
)
    ENGINE = MergeTree()
ORDER BY (task_id, created_at);

CREATE TABLE default.task_dependencies
(
    parent_task_id UUID,
    child_task_id UUID
)
    ENGINE = MergeTree()
ORDER BY (parent_task_id, child_task_id);
