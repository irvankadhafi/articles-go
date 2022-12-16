-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS "articles" (
    id SERIAL PRIMARY KEY,
    author TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT 'now()'
);


-- +migrate Down
DROP TABLE IF EXISTS "articles";