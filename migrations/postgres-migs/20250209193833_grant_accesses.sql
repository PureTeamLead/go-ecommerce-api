-- +goose Up
-- +goose StatementBegin
GRANT ALL PRIVILEGES ON DATABASE e_commerce TO themaxs;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
REVOKE ALL PRIVILEGES ON DATABASE e_commerce FROM themaxs;
-- +goose StatementEnd
