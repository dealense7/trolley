package main

import (
	"math/rand"
	"storePrices/internal/domain/parser"
	"storePrices/internal/domain/parser/strategies"
	"storePrices/internal/platform/conf"
	"storePrices/internal/platform/database"
	"storePrices/internal/platform/logger"
	"time"

	"github.com/jmoiron/sqlx"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	app := fx.New(
		// 1. Infrastructure
		fx.Provide(
			conf.NewConfig,
			logger.NewFactory,
			database.NewDatabase,
		),

		// 2. Provide Service and Strategies separately
		fx.Provide(
			parser.NewParserService,     // The Manager
			strategies.NewGlovoStrategy, // The Worker
		),

		// 3. Register Strategies (The Wiring)
		// This tells Fx: "Get the Service and the GlovoStrategy, and put the Strategy inside the Service"
		fx.Invoke(func(s *parser.ParserService, glovo *strategies.GlovoStrategy) {
			s.AddStrategy(glovo)
		}),

		// 4. Run the actual job
		fx.Invoke(runWorker),
	)

	app.Run()
}

func runWorker(s *parser.ParserService, db *sqlx.DB, logFactory logger.Factory, shutdown fx.Shutdowner) {
	retailers := GetTargets(db, logFactory)
	log := logFactory.For("worker")

	go func() {
		log.Info("worker started", zap.Int("queue_size", len(retailers)))

		for _, retailer := range retailers {
			log.Info("processing job for retailer", zap.String("retailer", retailer.Name))

			for _, store := range retailer.Stores {
				log.Info("processing job for store", zap.String("city", store.City))

				err := s.ScrapeAndPrint(store)
				if err != nil {
					log.Error("scrape failed", zap.String("store", store.Url), zap.Error(err))
				}
			}
		}
		// Just to be police
		time.Sleep(time.Duration(2000+rand.Intn(3000)) * time.Millisecond)

		log.Info("all jobs completed")
		shutdown.Shutdown()
	}()
}
