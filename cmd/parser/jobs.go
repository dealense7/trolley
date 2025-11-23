package main

import (
	"storePrices/internal/domain/retailer"
	"storePrices/internal/platform/logger"

	"github.com/jmoiron/sqlx"
)

func GetTargets(db *sqlx.DB, log logger.Factory) []retailer.Model {
	r := retailer.NewRepository(db, log)
	items, _ := r.GetRetailersWithStores()

	return items
}
