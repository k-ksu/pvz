-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id            TEXT,
    client_id      TEXT,
    condition      TEXT,
    arrived_at     TIMESTAMPTZ,
    received_at    TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
