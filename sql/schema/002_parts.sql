-- +goose Up
CREATE TABLE parts (
    id UUID PRIMARY KEY,
    serial_number SERIAL UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    sku TEXT UNIQUE,
    parent_id UUID NOT NULL 
    REFERENCES bins(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE parts;
