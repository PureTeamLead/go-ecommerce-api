-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA IF NOT EXISTS ecom_scheme;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA IF EXISTS ecom_scheme CASCADE;
-- +goose StatementEnd
