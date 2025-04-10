CREATE TABLE default.boards
(
    id UUID,
    owner_id UUID,
    name String,
    color String,
    created_at DateTime DEFAULT now()
)
    ENGINE = MergeTree()
ORDER BY id;
