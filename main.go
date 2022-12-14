package main

import (
	"awesome-api/api"
	"awesome-api/config"
	"awesome-api/jwt"
	"awesome-api/logger"
	mailer "awesome-api/mail"
	"awesome-api/store"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/rs/zerolog"
)

func main() {
	var migrateCMD string

	flag.StringVar(&migrateCMD, "migrate", "up", "migration")
	flag.Parse()

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
		return
	}
	ctx := context.Background()

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
	apiDB := api.DB{
		ElibraryPostgres: db,
	}
	mailer := setupMail(config)
	jwt := setupJWT(config)
	tokenVerification := api.TokenVerificationConfig{
		Expiry: time.Duration(config.TokenVerificationExpirationMinute),
	}
	apiLogger := zlog.With().
		Str("component", "api").
		Logger()

	if migrateCMD != "" {
		err = store.Migrate(db, migrateCMD)
		if err != nil {
			zlog.Fatal().Err(err).Msg("failed migration")
		}
	}

	srv := api.NewServer(
		fmt.Sprintf("%s:%d", config.AppHost, config.AppPort),
		apiLogger,
		apiDB,
		tokenVerification,
		mailer,
		jwt,
	)
	srv.Run(ctx)
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

func setupMail(cfg config.Config) mailer.EmailSender {
	mailerConfig := &mailer.Config{
		AppUrl:       fmt.Sprintf("%s://%s:%d", cfg.AppProtocol, cfg.AppHost, cfg.AppPort),
		MailHost:     cfg.MailHost,
		MailPort:     cfg.MailPort,
		MailUsername: cfg.MailUsername,
		MailPassword: cfg.MailPassword,
	}
	return mailer.NewMail(mailerConfig)
}

func setupJWT(cfg config.Config) jwt.JWT {
	jwtCfg := jwt.JWTConfig{
		TokenAccessExpiration:  time.Duration(cfg.TokenAccessExpirationMinute),
		TokenRefreshExpiration: time.Duration(cfg.TokenRefreshExpirationMinute),
	}
	return jwt.NewJWT(jwtCfg)
}
