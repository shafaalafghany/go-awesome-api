package postgresql

import (
	"awesome-api/store"
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"
)

type UserStore struct {
	log zerolog.Logger
	db  *sql.DB
	ps  *userPrepareStatement
}

type userPrepareStatement struct {
	Insert                   *sql.Stmt
	FindOneById              *sql.Stmt
	FindOneByEmail           *sql.Stmt
	FindOneCredentialByEmail *sql.Stmt
	UpdateTokenIdById        *sql.Stmt
	DeleteTokenIdById        *sql.Stmt
}

func (us *UserStore) prepareStatement() error {
	storeName := "UserStore"
	var err error
	if us.ps.Insert, err = prepareStatement(us.db, storeName, "Insert", userInsert); err != nil {
		return err
	}
	if us.ps.FindOneByEmail, err = prepareStatement(us.db, storeName, "FindOneByEmail", userFindOneByEmail); err != nil {
		return err
	}
	if us.ps.FindOneCredentialByEmail, err = prepareStatement(us.db, storeName, "FindOneCredentialByEmail", userFindOneCredentialByEmail); err != nil {
		return err
	}
	if us.ps.UpdateTokenIdById, err = prepareStatement(us.db, storeName, "UpdateTokenIdById", userUpdateTokenId); err != nil {
		return err
	}
	if us.ps.FindOneById, err = prepareStatement(us.db, storeName, "FindOneById", userFindOneById); err != nil {
		return err
	}
	if us.ps.DeleteTokenIdById, err = prepareStatement(us.db, storeName, "DeleteTokenIdById", userDeleteTokenIdById); err != nil {
		return err
	}
	return nil
}

func NewUserStore(log zerolog.Logger, db *sql.DB) (*UserStore, error) {
	us := &UserStore{
		db:  db,
		log: log,
		ps:  &userPrepareStatement{},
	}
	err := us.prepareStatement()
	if err != nil {
		return nil, err
	}
	return us, nil
}

const userFindOneBase = `
SELECT id, email, fullname, is_verified,
token_id, token_verification, token_expiration
FROM "users"
`

const userFindOneByEmail = userFindOneBase + "WHERE email = $1"

func (us *UserStore) FindOneByEmail(ctx context.Context, email string) (*store.User, error) {
	row := us.ps.FindOneByEmail.QueryRowContext(ctx, email)
	return us.scanRow(row)
}

const userFindOneById = userFindOneBase + "WHERE id = $1"

func (us *UserStore) FindOneById(ctx context.Context, id int) (*store.User, error) {
	row := us.ps.FindOneById.QueryRowContext(ctx, id)
	return us.scanRow(row)
}

const userInsert = `
INSERT INTO "users" (
	email, password, fullname, is_verified,
	token_verification, token_expiration
) VALUES (
	$1, $2, $3, $4, $5, $6
)
`

func (us *UserStore) Insert(ctx context.Context, usr *store.UserRegister) error {
	_, err := us.ps.Insert.ExecContext(ctx,
		usr.Email, usr.Password, usr.Fullname,
		usr.IsVerified, usr.TokenVerification,
		usr.TokenExpiration,
	)
	if err != nil {
		return fmt.Errorf("failed to Insert: %w", err)
	}
	return nil
}

const userFindOneCredentialByEmail = `
SELECT id, email, password, is_verified
FROM "users"
WHERE email = $1
`

func (us *UserStore) FindOneCredentialByEmail(ctx context.Context, email string) (*store.User, error) {
	row := us.ps.FindOneCredentialByEmail.QueryRowContext(ctx, email)
	user := &store.User{}
	err := row.Scan(
		&user.ID, &user.Email, &user.Password,
		&user.IsVerified,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}
	return user, nil
}

const userUpdateTokenId = `
UPDATE "users" SET
token_id = $1
WHERE id = $2
`

func (us *UserStore) UpdateTokenIdById(ctx context.Context, token string, id int) error {
	_, err := us.ps.UpdateTokenIdById.ExecContext(ctx, token, id)
	if err != nil {
		return fmt.Errorf("failed to UpdateTokenIdById: %w", err)
	}
	return nil
}

const userDeleteTokenIdById = `
UPDATE "users" SET
token_id = NULL
WHERE id = $1
`

func (us *UserStore) DeleteTokenIdById(ctx context.Context, id int) error {
	_, err := us.ps.DeleteTokenIdById.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to DeleteTokenIdById: %w", err)
	}
	return nil
}

func (us *UserStore) scanRow(row *sql.Row) (*store.User, error) {
	user := &store.User{}
	err := row.Scan(
		&user.ID, &user.Email, &user.Fullname,
		&user.IsVerified, &user.TokenID, &user.TokenVerification,
		&user.TokenExpiration,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scanRow: %w", err)
	}
	return user, nil
}
