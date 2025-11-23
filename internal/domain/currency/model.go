package currency

import "time"

type Id int64

const (
	USDollar Id = 1
	Euro     Id = 2
)

type Model struct {
	Id     Id     `db:"id"`
	Code   string `db:"code"`   // 'USD', 'EUR', 'GBP'
	Symbol string `db:"symbol"` // '$', '€', '£'
	Name   string `db:"name"`   // full name of the currency
}

type ModelValue struct {
	Id         int64     `db:"id"`          // unique ID for each record
	CurrencyId int64     `db:"currency_id"` // 'USD', 'EUR', 'GBP'
	Value      float64   `db:"value"`       // currency value
	Date       time.Time `db:"date"`        // the date of the value
	Status     bool      `db:"status"`      // if true, this value can be used
}
