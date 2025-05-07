-- +goose Up
CREATE TABLE bins (
    id UUID PRIMARY KEY,
    serial_number SERIAl UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT NOT NULL,
    sku TEXT UNIQUE,
    parent_id UUID DEFAULT NULL 
    REFERENCES bins(id) ON DELETE SET NULL,
    parent_id_or_null UUID GENERATED ALWAYS AS (COALESCE(parent_id, '00000000-0000-0000-0000-000000000000')) STORED
);

ALTER TABLE bins
ADD CONSTRAINT bins_name_parent_bin_key UNIQUE (name, parent_id_or_null);

-- +goose Down
DROP TABLE bins;
