package storage

import (
	"avito-pvz-service/config"
	logger "avito-pvz-service/pkg"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Log.Fatalf("Error connecting to DB: %v", err)
	}

	err = db.Ping()
	if err != nil {
		logger.Log.Fatalf("DB is not reachable: %v", err)
	}
	logger.Log.Info("Connected to the database")
	return db
}
