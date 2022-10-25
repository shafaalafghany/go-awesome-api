package api

import (
	"awesome-api/api/handler/auth"
	mailer "awesome-api/mail"
	"awesome-api/store"
	"awesome-api/store/postgresql"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Server struct {
	Addr              string
	logger            zerolog.Logger
	stores            *stores
	tokenVerification TokenVerificationConfig
	mailer            mailer.EmailSender
}

type DB struct {
	ElibraryPostgres *sql.DB
}

type stores struct {
	userStore store.UserStore
}

type TokenVerificationConfig struct {
	Expiry time.Duration
}

func NewServer(
	addr string,
	logger zerolog.Logger,
	db DB,
	tokenVerification TokenVerificationConfig,
	mailer mailer.EmailSender,
) *Server {
	s := &Server{
		Addr:              addr,
		logger:            logger,
		tokenVerification: tokenVerification,
		mailer:            mailer,
	}
	var err error
	s.stores, err = initStores(s, db)
	if err != nil {
		logger.Fatal().Err(err)
	}
	return s
}

func initStores(s *Server, db DB) (*stores, error) {
	stores := &stores{}
	var err error
	if stores.userStore, err = postgresql.NewUserStore(
		s.logger.With().Str("store", "user_store").Logger(),
		db.ElibraryPostgres,
	); err != nil {
		return nil, err
	}
	return stores, nil
}

func (s *Server) Run(ctx context.Context) {
	handler := chi.NewMux()
	handler.Mount("/", handlers(s))

	srv := &http.Server{
		Addr:    s.Addr,
		Handler: handler,
	}
	s.logger.Info().Msgf("serving http on : %s\n", s.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Error().Err(fmt.Errorf("failed to serve http: %w", err))
	}
}

func handlers(s *Server) http.Handler {
	h := chi.NewMux()

	h.Post("/auth", auth.Signup(
		s.logger,
		s.stores.userStore,
		s.tokenVerification.Expiry,
		s.mailer,
	))
	return h
}
