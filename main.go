package main

import (
	"awesome-api/config"
	"awesome-api/logger"
	"awesome-api/store"
	"database/sql"
	"fmt"
	"log"

	"github.com/rs/zerolog"
)

func main() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
		return
	}
	var zlog zerolog.Logger
	logConfig := logger.Config{
		Level:  config.LoggerLevel,
		Output: config.LoggerOutput,
	}
	zlog, err = logger.NewZerolog(logConfig)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create zerolog console: %w", err))
	}
	db := setupPg(config, zlog)
	fmt.Println(db)
}

func setupPg(cfg config.Config, logger zerolog.Logger) *sql.DB {
	db, err := store.OpenPg(cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open postgres db")
		return nil
	}
	err = db.Ping()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to ping postgres db")
		return nil
	}
	return db
}
