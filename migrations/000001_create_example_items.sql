-- +goose Up
CREATE TABLE example_items (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX idx_example_items_deleted_at ON example_items(deleted_at);

-- +goose Down
DROP TABLE IF EXISTS example_items;