CREATE TABLE default.task_statuses
(
    task_id UUID,
    is_completed Bool,
    updated_at DateTime DEFAULT now()
)
    ENGINE = ReplacingMergeTree()
ORDER BY task_id;
