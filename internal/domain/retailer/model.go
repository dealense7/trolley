package retailer

import (
	"time"
)

type Model struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	LogoURL   string    `db:"logo_url"`
	CreatedAt time.Time `db:"created_at"`
	Stores    []Store
}

type Store struct {
	Id         int64  `db:"id"`
	RetailerId int64  `db:"retailer_id"`
	CountryId  int64  `db:"country_id"`
	CurrencyId int64  `db:"currency_id"`
	City       string `db:"city"` // CHAR(3)
	Url        string `db:"base_url"`
	IsActive   bool   `db:"is_active"`
}
