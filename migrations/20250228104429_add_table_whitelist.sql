-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS whitelist (
    id SERIAL PRIMARY KEY,
    subnet CIDR
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS whitelist;
-- +goose StatementEnd
