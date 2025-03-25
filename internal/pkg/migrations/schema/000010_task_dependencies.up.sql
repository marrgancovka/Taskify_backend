CREATE TABLE default.task_dependencies
(
    parent_task_id UUID,
    child_task_id UUID
)
    ENGINE = MergeTree()
ORDER BY (parent_task_id, child_task_id);
