package country

import (
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

func (r *Repository) GetCountriesByIds(ids []int64) ([]Model, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// Fetch countries
	items := []Model{}
	query, args, err := sqlx.In("SELECT * FROM countries WHERE id IN (?)", ids)
	if err != nil {
		r.log.Error("Building IN query for countries", zap.Error(err))
		return nil, err
	}

	query = r.db.Rebind(query)
	if err := r.db.Select(&items, query, args...); err != nil {
		r.log.Error("Selecting Countries", zap.Error(err))
		return nil, err
	}

	return items, nil
}
