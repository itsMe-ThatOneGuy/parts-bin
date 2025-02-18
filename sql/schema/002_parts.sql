-- +goose Up
CREATE TABLE parts (
    part_id SERIAL PRIMARY KEY,
    id UUID UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    sku TEXT,
    parent_id UUID NOT NULL 
    REFERENCES bins(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE parts;
