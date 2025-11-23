package seeders

import (
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Currency struct {
	Code   string
	Symbol string
	Name   string
}

func SeedCurrency(db *sqlx.DB, log *zap.Logger) {
	items := []Currency{
		{Code: "EUR", Symbol: "€", Name: "Euro"},
		{Code: "USD", Symbol: "$", Name: "US Dollar"},
	}

	start := time.Now()
	log.Info("Currency Seeder Started")

	db.NamedExec(`INSERT INTO currencies (code, symbol, name) VALUES (:code, :symbol, :name) ON DUPLICATE KEY UPDATE code = code`, items)

	log.Info("Country Seeder Finished:", zap.String("time", time.Since(start).String()))

}
