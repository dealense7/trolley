package main

import (
	"go.uber.org/fx"
	"storePrices/internal/domain/product"
	"storePrices/internal/platform/conf"
	"storePrices/internal/platform/database"
	"storePrices/internal/platform/logger"
	"storePrices/internal/platform/server"
)

func main() {
	fx.New(
		// General Staff
		fx.Provide(
			conf.NewConfig,
			logger.NewFactory,
			database.NewDatabase,
			server.New,
		),

		// Load Domain Modules
		product.Module,

		// Start the Server
		fx.Invoke(server.Start),
	).Run()
}
