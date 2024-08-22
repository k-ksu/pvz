-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_order_id ON orders (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_order_id;
-- +goose StatementEnd
