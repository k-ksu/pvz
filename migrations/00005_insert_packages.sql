-- +goose Up
-- +goose StatementBegin
INSERT INTO packages_info (package, surcharge, max_weight)
VALUES
('noPackage', 0, -1),
('plasticBag', 5, 10000),
('box', 20, 30000),
('film', 1, -1)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM packages_info
WHERE package = 'noPackage' OR package = 'plasticBag' OR package = 'box' OR package = 'film';

-- +goose StatementEnd
