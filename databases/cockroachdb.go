package databases

import (
	"context"
	"log"
	"pwstore/types"

	"github.com/jackc/pgx/v5"
)

type Cockroachdb[T any] struct {
	connectionString string
	context          context.Context
	conn             *pgx.Conn
	fullTableName    string
}

func NewCockroachdb[T any](connString string, fullTableName string) types.Database[T] {
	return &Cockroachdb[T]{
		connectionString: connString,
		context:          context.Background(),
		fullTableName:    fullTableName,
	}
}

func (db *Cockroachdb[T]) GetConn() *pgx.Conn {
	return db.conn
}

func (db *Cockroachdb[T]) Connect() (bool, error) {
	conn, err := pgx.Connect(db.context, db.connectionString)
	if err != nil {
		log.Printf("cannot connect\nreason:\n%s", err)
		return false, err
	}
	db.conn = conn
	return true, nil
}

func (db *Cockroachdb[T]) Close() (bool, error) {
	err := db.conn.Close(db.context)
	if err != nil {
		log.Printf("cannot close\nreason:\n%s", err)
		return false, err
	}
	return true, nil
}

func (db *Cockroachdb[T]) Query(sql string, args ...any) (*[]T, error) {
	rows, err := db.conn.Query(db.context, sql, args...)
	if err != nil {
		log.Printf("cannot query\nreason:\n%s", err)
		return &[]T{}, err
	}

	data := []T{}
	for rows.Next() {
		record, err := pgx.RowToStructByPos[T](rows)
		if err != nil {
			log.Printf("cannot parse\nreason:\n%s", err)
		} else {
			data = append(data, record)
		}
	}
	return &data, nil
}
