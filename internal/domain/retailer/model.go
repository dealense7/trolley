package retailer

import (
	"storePrices/internal/domain/country"
	"time"
)

type Model struct {
	Id        int64     `db:"id"`
	Name      string    `db:"name"`
	LogoURL   *string   `db:"logo_url"`
	CreatedAt time.Time `db:"created_at"`
	Stores    []Store
}

type Store struct {
	Id         int64 `db:"id"`
	RetailerId int64 `db:"retailer_id"`
	Country    *country.Model
	CountryId  int64  `db:"country_id"`
	CurrencyId int64  `db:"currency_id"`
	City       string `db:"city"` // CHAR(3)
	Url        string `db:"base_url"`
	IsActive   bool   `db:"is_active"`
}

type Retailers []Model

func (r Retailers) IdList() []int64 {
	var list []int64
	for _, retailer := range r {
		list = append(list, retailer.Id)
	}
	return list
}

type Stores []Store

func (s Stores) CountryIdList() []int64 {
	var list []int64
	for _, st := range s {
		list = append(list, st.CountryId)
	}
	return list
}
