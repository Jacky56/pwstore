package types

import (
	"github.com/jackc/pgx/v5"
)

type Database[T any] interface {
	Connect() (bool, error)
	Query(string, ...any) (*[]T, error)
	Close() (bool, error)
	GetConn() *pgx.Conn
}
