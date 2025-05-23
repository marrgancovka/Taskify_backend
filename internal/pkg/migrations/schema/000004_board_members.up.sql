CREATE TABLE default.board_members
(
    board_id UUID,
    user_id UUID,
    role_id UUID,
    is_favourite Bool,
)
    ENGINE = ReplacingMergeTree()
ORDER BY (board_id, user_id);
