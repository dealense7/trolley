package seeders

import (
	"storePrices/internal/domain/country"
	"storePrices/internal/domain/currency"
	"storePrices/internal/domain/retailer"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func SeedStores(db *sqlx.DB, log *zap.Logger) {
	items := []retailer.Model{
		{
			Name: "Dia",
			Stores: []retailer.Store{
				{
					CountryId:  int64(country.Spain),
					CurrencyId: int64(currency.Euro),
					City:       "MAD",
					Url:        "https://glovoapp.com/es/es/madrid/stores/supermercado-dia-2",
				},
			},
		},
	}

	start := time.Now()
	log.Info("Store Seeder Started")

	for _, r := range items {

		// Insert retailer OR get existing
		_, err := db.NamedExec(`INSERT INTO retailers (name) VALUES (:name) ON DUPLICATE KEY UPDATE name = name`, r)
		if err != nil {
			log.Error("failed inserting retailer", zap.Error(err))
			continue
		}

		var retailerID int64
		err = db.Get(&retailerID, `SELECT id FROM retailers WHERE name = ?`, r.Name)
		if err != nil {
			log.Error("failed fetching retailer id", zap.Error(err))
			continue
		}

		// Insert stores
		for _, s := range r.Stores {
			s.RetailerId = retailerID

			_, err = db.NamedExec(`
				INSERT INTO stores (retailer_id, country_id, currency_id, city, base_url)
				VALUES (:retailer_id, :country_id, :currency_id, :city, :base_url)
				ON DUPLICATE KEY UPDATE base_url = VALUES(base_url)
			`, s)
			if err != nil {
				log.Error("failed inserting store", zap.Error(err))
				continue
			}
		}
	}

	log.Info("Store Seeder Finished", zap.String("time", time.Since(start).String()))

}
