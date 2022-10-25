package postgresql

import (
	"database/sql"
	"fmt"
)

func prepareStatement(db *sql.DB, storeName, queryName, sql string) (*sql.Stmt, error) {
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, fmt.Errorf("%s failed to prepare %s: %w", storeName, queryName, err)
	}
	return stmt, nil
}
