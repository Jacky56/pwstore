package types

import (
	"pwstore/data"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Database interface {
	Connect() (bool, error)
	Query(string, ...any) (*[]any, error)
	ListUsers() (*[]data.User, error)
	GetUser(string) (data.User, error)
	SetUser(*data.User) error
	GetPasswordStore(uuid.UUID) (data.PasswordStore, error)
	SetPasswordStore(*data.PasswordStore) error
	Close() (bool, error)
	GetConn() *pgx.Conn
}
