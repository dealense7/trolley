package main

import (
	"storePrices/internal/domain/retailer"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func GetTargets(db *sqlx.DB, log *zap.Logger) []retailer.Model {
	items, _ := retailer.GetRetailersWithStores(db, log)

	return items
}
