-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_client_id ON orders (client_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_client_id;
-- +goose StatementEnd
