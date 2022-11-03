package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID                int
	Email             string
	Password          sql.NullString
	Fullname          string
	IsVerified        bool
	TokenID           sql.NullString
	TokenVerification sql.NullString
	TokenExpiration   sql.NullString
}

type UserRegister struct {
	Email             string
	Password          string
	Fullname          string
	IsVerified        bool
	TokenVerification string
	TokenExpiration   string
}

type UserStore interface {
	Insert(ctx context.Context, usr *UserRegister) error
	FindOneById(ctx context.Context, id int) (*User, error)
	FindOneByEmail(ctx context.Context, email string) (*User, error)
	FindOneCredentialByEmail(ctx context.Context, email string) (*User, error)
	UpdateTokenIdById(ctx context.Context, token string, id int) error
	DeleteTokenIdById(ctx context.Context, id int) error
}
