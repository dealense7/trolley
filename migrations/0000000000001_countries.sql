-- +goose Up
-- +goose StatementBegin
CREATE TABLE countries
(
    id   INT AUTO_INCREMENT PRIMARY KEY,
    code CHAR(3) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS countries;
-- +goose StatementEnd
