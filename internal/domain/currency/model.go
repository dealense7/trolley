package currency

import "time"

type Code string

const (
	GeorgianLari  Code = "GEL"
	USDollar      Code = "USD"
	Euro          Code = "EUR"
	PoundSterling Code = "GBP"
)

// Currency represents the currency metadata
type Currency struct {
	Code   Code   `db:"code"`   // 'USD', 'EUR', 'GBP'
	Symbol string `db:"symbol"` // '$', '€', '£'
	Name   string `db:"name"`   // full name of the currency
}

// CurrencyValue represents the value of a currency on a given date
type Value struct {
	Id     int64     `db:"id"`     // unique ID for each record
	Code   Code      `db:"code"`   // 'USD', 'EUR', 'GBP'
	Value  float64   `db:"value"`  // currency value
	Date   time.Time `db:"date"`   // the date of the value
	Status bool      `db:"status"` // if true, this value is used
}
