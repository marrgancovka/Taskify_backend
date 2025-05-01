CREATE TABLE default.tasks
(
    id UUID,
    board_id UUID,
    section_id UUID,
    name String,
    description String,
    due_date DateTime,
    priority Int32,
    percent Int32 default 0,
    created_at DateTime DEFAULT now()
)
    ENGINE = MergeTree
ORDER BY (board_id);
