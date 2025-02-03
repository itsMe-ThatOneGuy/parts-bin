-- +goose Up
CREATE TABLE bins (
    id UUID PRIMARY KEY, 
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    parent_bin UUID DEFAULT NULL,
    FOREIGN KEY (parent_bin) 
    REFERENCES bins(id) ON DELETE SET NULL,
    UNIQUE (name, parent_bin)
);

-- +goose Down
DROP TABLE bins;
