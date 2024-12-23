-- +goose Up
CREATE TABLE bins (
    id UUID PRIMARY KEY, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL UNIQUE,
    parent_bin UUID,
);

-- +goose Down
DROP TABLE bins;
