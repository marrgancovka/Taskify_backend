CREATE TABLE default.board_roles
(
    id UUID,
    board_id UUID,
    name String
)
    ENGINE = MergeTree()
ORDER BY id;
