-- +goose Up

-- +goose StatementBegin
CREATE TABLE currencies
(
    id     INT AUTO_INCREMENT PRIMARY KEY,
    code   CHAR(3) UNIQUE,    -- 'USD', 'EUR', 'GBP'
    symbol VARCHAR(5), -- '$', '€', '£'
    name   VARCHAR(50)
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE currency_values
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    currency_id INT            NOT NULL,
    value       NUMERIC(18, 6) NOT NULL, -- currency value
    date        DATE           NOT NULL, -- the date of the value
    status      BOOLEAN DEFAULT FALSE,   -- if true we use it
    FOREIGN KEY (currency_id) REFERENCES currencies (id)
);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS currency_values;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS currencies;
-- +goose StatementEnd
