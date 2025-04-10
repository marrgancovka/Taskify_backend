CREATE TABLE default.users
(
    id UUID,
    username String,
    email String,
    password String,
    created_at DateTime DEFAULT now()
)
    ENGINE = MergeTree()
        ORDER BY id;
