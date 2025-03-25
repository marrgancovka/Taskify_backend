CREATE TABLE default.tasks
(
    id UUID,
    board_id UUID,
    section_id UUID,
    name String,
    description String,
    due_date DateTime,
    priority Int32, -- 1 = низкий, 2 = средний, 3 = высокий
    created_at DateTime DEFAULT now()
)
    ENGINE = MergeTree()
ORDER BY (board_id, section_id, due_date, priority);
