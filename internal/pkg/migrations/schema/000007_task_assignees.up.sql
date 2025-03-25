CREATE TABLE default.task_assignees
(
    task_id UUID,
    user_id UUID
)
    ENGINE = MergeTree()
ORDER BY (task_id, user_id);
