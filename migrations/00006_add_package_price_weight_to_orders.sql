-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
ADD package    TEXT NOT NULL DEFAULT 'noPackage',
ADD price      INTEGER NOT NULL DEFAULT 0,
ADD max_weight INTEGER NOT NULL DEFAULT -1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
DROP COLUMN package,
DROP COLUMN price,
DROP COLUMN max_weight;
-- +goose StatementEnd
