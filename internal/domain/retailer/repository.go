package retailer

import (
	"storePrices/internal/domain/country"
	"storePrices/internal/platform/logger"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Repository struct {
	db            *sqlx.DB
	log           *zap.Logger
	loggerFactory logger.Factory
}

func NewRepository(db *sqlx.DB, loggerFactory logger.Factory) *Repository {
	return &Repository{
		db:            db,
		log:           loggerFactory.For(logger.Country),
		loggerFactory: loggerFactory,
	}
}
func (r *Repository) GetRetailersWithStores() ([]Model, error) {
	items := Retailers{}
	err := r.db.Select(&items, "SELECT * FROM retailers ORDER BY id DESC")
	if err != nil {
		r.log.Error("Selecting Retailers", zap.Error(err))
		return nil, err
	}

	// Attach Stores
	stores, err := r.getStoresForRetailers(items.IdList())
	if err != nil {
		return nil, err
	}

	// Map stores by retailer id
	storeMap := make(map[int64]Stores)
	for _, s := range stores {
		storeMap[s.RetailerId] = append(storeMap[s.RetailerId], s)
	}

	for i, r := range items {
		items[i].Stores = storeMap[r.Id]
	}

	return items, nil
}

func (r *Repository) getStoresForRetailers(retailerIds []int64) (Stores, error) {
	if len(retailerIds) == 0 {
		return nil, nil
	}

	items := Stores{}
	query, args, err := sqlx.In("SELECT * FROM stores WHERE retailer_id IN (?)", retailerIds)
	if err != nil {
		r.log.Error("Building IN query", zap.Error(err))
		return nil, err
	}

	query = r.db.Rebind(query)

	err = r.db.Select(&items, query, args...)
	if err != nil {
		r.log.Error("Selecting Stores", zap.Error(err))
		return nil, err
	}

	// Attach countries to stores
	countryRepo := country.NewRepository(r.db, r.loggerFactory)
	countries, _ := countryRepo.GetCountriesByIds(items.CountryIdList())

	// Map Countries by id
	countryMap := make(map[int64]*country.Model)
	for _, c := range countries {
		countryMap[int64(c.Id)] = &c
	}

	for i, v := range items {
		items[i].Country = countryMap[v.CountryId]
	}

	return items, nil
}
