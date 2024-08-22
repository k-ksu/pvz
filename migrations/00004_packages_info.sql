-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS packages_info (
    package    TEXT,
    surcharge  INTEGER NOT NULL,
    max_weight INTEGER NOT NULL DEFAULT -1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE packages_info;
-- +goose StatementEnd
