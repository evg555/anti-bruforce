-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS blacklist (
    id SERIAL PRIMARY KEY,
    subnet CIDR
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blacklist;
-- +goose StatementEnd
