-- +goose Up
CREATE TABLE IF NOT EXISTS hubs (
    id BIGSERIAL PRIMARY KEY,
    author TEXT NOT NULL,
    title TEXT,
    structure JSONB NOT NULL,
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS hubs;
