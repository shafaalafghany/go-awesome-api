package store

import (
	"awesome-api/config"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func OpenPg(cfg config.Config) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		cfg.DbUsername,
		cfg.DbPassword,
		cfg.DbHost,
		cfg.DbPort,
		cfg.DbName,
	)
	parsedConn, err := pgx.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	conn := stdlib.RegisterConnConfig(parsedConn)
	db, err := sql.Open("pgx", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
