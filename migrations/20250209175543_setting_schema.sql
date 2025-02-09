-- +goose Up
-- +goose StatementBegin
SET search_path TO example;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO public;
-- +goose StatementEnd
