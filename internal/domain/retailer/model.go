package retailer

import (
	"time"
)

type Retailer struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	LogoURL   string    `db:"logo_url"`
	CreatedAt time.Time `db:"created_at"`
}

type Store struct {
	Id           int64  `db:"id"`
	RetailerId   int64  `db:"retailer_id"`
	Country      string `db:"country_code"` // CHAR(2)
	City         string `db:"country_code"` // CHAR(3)
	CurrencyCode string `db:"currency_code"`
	BaseURL      string `db:"base_url"`
	IsActive     bool   `db:"is_active"`
}
