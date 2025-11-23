package seeders

import (
	"storePrices/internal/domain/currency"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func SeedCurrency(db *sqlx.DB, log *zap.Logger) {
	items := []currency.Model{
		{Id: currency.Euro, Code: "EUR", Symbol: "€", Name: "Euro"},
		{Id: currency.USDollar, Code: "USD", Symbol: "$", Name: "US Dollar"},
	}

	start := time.Now()
	log.Info("Currency Seeder Started")

	db.NamedExec(`INSERT INTO currencies (id, code, symbol, name) VALUES (:id, :code, :symbol, :name) ON DUPLICATE KEY UPDATE id = id`, items)

	log.Info("Country Seeder Finished:", zap.String("time", time.Since(start).String()))

}
