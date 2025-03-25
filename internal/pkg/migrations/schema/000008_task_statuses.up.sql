CREATE TABLE default.task_statuses
(
    task_id UUID,
    is_completed UInt8, -- 0 = не выполнено, 1 = выполнено
    updated_at DateTime DEFAULT now()
)
    ENGINE = ReplacingMergeTree()
ORDER BY task_id;
