package main

import (
	"storePrices/internal/platform/conf"
	"storePrices/internal/platform/database"
	"storePrices/internal/platform/logger"
	"storePrices/seeders"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		// General Staff
		fx.Provide(
			conf.NewConfig,
			logger.NewFactory,
			database.NewDatabase,

			func(f logger.Factory) (*zap.Logger, error) {
				l := f.For(logger.General)
				return l, nil
			},
		),

		// Seed Database
		fx.Invoke(
			seeders.SeedCountries,
			seeders.SeedCurrency,
		),
		fx.Invoke(func(shutdown fx.Shutdowner) {
			_ = shutdown.Shutdown()
		}),
	).Run()
}
