package retailer

import (
	"storePrices/internal/domain/country"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func GetRetailersWithStores(db *sqlx.DB, log *zap.Logger) ([]Model, error) {
	items := Retailers{}
	err := db.Select(&items, "SELECT * FROM retailers ORDER BY id DESC")
	if err != nil {
		log.Error("Selecting Retailers", zap.Error(err))
		return nil, err
	}

	_ = getStoresForRetailers(&items, db, log)

	return items, nil
}

func getStoresForRetailers(items *Retailers, db *sqlx.DB, log *zap.Logger) error {
	if len(*items) == 0 {
		return nil
	}

	stores := []Store{}
	query, args, err := sqlx.In("SELECT * FROM stores WHERE retailer_id IN (?)", items.IdList())
	if err != nil {
		log.Error("Building IN query", zap.Error(err))
		return err
	}

	query = db.Rebind(query)

	err = db.Select(&stores, query, args...)
	if err != nil {
		log.Error("Selecting Stores", zap.Error(err))
		return err
	}

	_ = getCountriesForRetailers(&stores, db, log)

	// Map stores so will be able to attach to retailer easily
	storeMap := make(map[int64][]Store)
	for _, s := range stores {
		storeMap[s.RetailerId] = append(storeMap[s.RetailerId], s)
	}

	for i := range *items {
		r := &(*items)[i]
		r.Stores = storeMap[r.Id]
	}

	return nil
}

func getCountriesForRetailers(items *[]Store, db *sqlx.DB, log *zap.Logger) error {
	if len(*items) == 0 {
		return nil
	}

	// Collect unique country IDs
	countryIDs := make([]int64, 0, len(*items))
	countrySet := make(map[int64]struct{})
	for _, r := range *items {
		if _, exists := countrySet[r.CountryId]; !exists {
			countryIDs = append(countryIDs, r.CountryId)
			countrySet[r.CountryId] = struct{}{}
		}
	}

	if len(countryIDs) == 0 {
		return nil
	}

	// Fetch countries
	countries := []country.Model{}
	query, args, err := sqlx.In("SELECT * FROM countries WHERE id IN (?)", countryIDs)
	if err != nil {
		log.Error("Building IN query for countries", zap.Error(err))
		return err
	}

	query = db.Rebind(query)
	if err := db.Select(&countries, query, args...); err != nil {
		log.Error("Selecting Countries", zap.Error(err))
		return err
	}

	// Map countries by ID
	countryMap := make(map[int64]*country.Model)
	for _, c := range countries {
		countryMap[int64(c.Id)] = &c
	}

	// Attach countries to retailers
	for i := range *items {
		r := &(*items)[i]
		r.Country = countryMap[r.CountryId]
	}

	return nil
}
