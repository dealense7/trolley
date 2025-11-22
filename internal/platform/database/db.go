package database

import (
	"context"
	"storePrices/internal/platform/conf"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL Driver
	"github.com/jmoiron/sqlx"
)

func NewDatabase(cfg *conf.Config) (*sqlx.DB, error) {
	connectionString := cfg.DB.DSN()

	// "root:pass@tcp(127.0.0.1:3306)/dbName?parseTime=true"
	db, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Set strict timeouts and pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
