-- +goose Up
-- +goose StatementBegin
CREATE TABLE retailers
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    logo_url    TEXT,
    website_url TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE stores
(
    id            INT AUTO_INCREMENT PRIMARY KEY,
    retailer_id   INT NOT NULL,
    country_code  INT NOT NULL,
    currency_code INT NOT NULL,
    base_url      TEXT,
    is_active     BOOLEAN DEFAULT TRUE,
    UNIQUE (retailer_id, country_code),
    FOREIGN KEY (retailer_id) REFERENCES retailers (id) ON DELETE CASCADE,
    FOREIGN KEY (country_code) REFERENCES countries (id) ON DELETE CASCADE,
    FOREIGN KEY (currency_code) REFERENCES currencies (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stores;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS retailers;
-- +goose StatementEnd
