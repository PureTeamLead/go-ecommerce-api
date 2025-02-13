-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    price FLOAT NOT NULL,
    name TEXT NOT NULL,
    category TEXT DEFAULT 'N/A',
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    seller users
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
