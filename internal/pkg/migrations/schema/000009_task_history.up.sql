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
