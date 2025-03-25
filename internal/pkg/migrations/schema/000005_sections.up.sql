CREATE TABLE default.sections
(
    id UUID,
    board_id UUID,
    name String,
    position Int32
)
    ENGINE = MergeTree()
ORDER BY (board_id, position);
