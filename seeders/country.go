package seeders

import (
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Country struct {
	Code string
	Name string
}

func SeedCountries(db *sqlx.DB, log *zap.Logger) {
	items := []Country{
		{Code: "ITA", Name: "Italy"},
		{Code: "ESP", Name: "Spain"},
		{Code: "BAR", Name: "Barcelona"},
	}

	start := time.Now()
	log.Info("Country Seeder Started")

	db.NamedExec(`INSERT INTO countries (code, name) VALUES (:code, :name) ON DUPLICATE KEY UPDATE code = code`, items)

	log.Info("Country Seeder Finished:", zap.String("time", time.Since(start).String()))

}
