package seeders

import (
	"storePrices/internal/domain/country"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func SeedCountries(db *sqlx.DB, log *zap.Logger) {
	items := []country.Model{
		{Id: country.Italy, Code: "ITA", Name: "Italy"},
		{Id: country.Spain, Code: "ESP", Name: "Spain"},
	}

	start := time.Now()
	log.Info("Country Seeder Started")

	db.NamedExec(`INSERT INTO countries (id, code, name) VALUES (:id, :code, :name) ON DUPLICATE KEY UPDATE id = id`, items)

	log.Info("Country Seeder Finished:", zap.String("time", time.Since(start).String()))

}
