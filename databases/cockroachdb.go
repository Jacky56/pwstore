package databases

import (
	"context"
	"fmt"
	"log"
	"pwstore/commons"
	"pwstore/data"
	"pwstore/types"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Cockroachdb struct {
	connectionString   string
	context            context.Context
	conn               *pgx.Conn
	TableUser          string
	TablePasswordStore string
}

func NewCockroachdb(ce *commons.ConfigEntry) types.Database {
	return &Cockroachdb{
		connectionString:   ce.GetConnString(),
		TableUser:          ce.GetUserTableName(),
		TablePasswordStore: ce.GetPasswordStoreTableName(),
		context:            context.Background(),
	}
}

func (db *Cockroachdb) GetConn() *pgx.Conn {
	return db.conn
}

func (db *Cockroachdb) Connect() (bool, error) {
	conn, err := pgx.Connect(db.context, db.connectionString)
	if err != nil {
		log.Printf("cannot connect\nreason:\n%s", err)
		return false, err
	}
	db.conn = conn
	return true, nil
}

func (db *Cockroachdb) Close() (bool, error) {
	err := db.conn.Close(db.context)
	if err != nil {
		log.Printf("cannot close\nreason:\n%s", err)
		return false, err
	}
	return true, nil
}

func (db *Cockroachdb) Query(sql string, args ...any) (*[]any, error) {
	rows, err := db.conn.Query(db.context, sql, args...)
	if err != nil {
		log.Printf("cannot query\nreason:\n%s", err)
		return &[]any{}, err
	}

	data := []any{}
	for rows.Next() {
		// record, err := pgx.RowToStructByName(rows)
		record, err := rows.Values()
		if err != nil {
			log.Printf("cannot parse\nreason:\n%s", err)
		} else {
			data = append(data, record)
		}
	}
	return &data, nil
}

func (db *Cockroachdb) ListUsers() (*[]data.User, error) {
	sql := fmt.Sprintf(`
	select 
		uuid,
		email,
		provider,
		'' as password 
	from %s;
	`, db.TableUser)
	var qem pgx.QueryExecMode = 4
	rows, err := db.conn.Query(db.context, sql, qem)
	if err != nil {
		log.Printf("cannot query\nreason:\n%s", err)
		return &[]data.User{}, err
	}
	users, err := pgx.CollectRows[data.User](rows, pgx.RowToStructByName[data.User])
	if err != nil {
		log.Printf("cannot collect users\nreason:\n%s", err)
		return &[]data.User{}, err
	}
	return &users, nil
}

func (db *Cockroachdb) GetUser(email string) (data.User, error) {
	sql := fmt.Sprintf(`
	select * from %s
	where lower(email) = lower($1)
	limit 1;
	`, db.TableUser)
	rows, err := db.conn.Query(db.context, sql, email)

	user, err := pgx.CollectOneRow[data.User](rows, pgx.RowToStructByName[data.User])
	if err != nil {
		log.Printf("cannot query user\nreason:\n%s", err)
		return data.User{}, err
	}
	return user, nil
}

func (db *Cockroachdb) SetUser(user *data.User) error {
	sql := fmt.Sprintf(`
	upsert into %s (email, password, provider)
	values ($1, $2, $3);
	`, db.TableUser)
	_, err := db.conn.Exec(db.context, sql, user.Email, user.Password, user.Provider)
	if err != nil {
		log.Printf("cannot upsert\nreason:\n%s", err)
		return err
	}
	return nil
}

func (db *Cockroachdb) GetPasswordStore(user_uuid uuid.UUID) (data.PasswordStore, error) {
	sql := fmt.Sprintf(`
	select * from %s
	where uuid = $1
	limit 1;
	`, db.TablePasswordStore)
	rows, err := db.conn.Query(db.context, sql, user_uuid)

	passwordStore, err := pgx.CollectOneRow[data.PasswordStore](rows, pgx.RowToStructByName[data.PasswordStore])
	if err != nil {
		log.Printf("cannot query passowrdStore\nreason:\n%s", err)
		return data.PasswordStore{}, err
	}
	return passwordStore, nil
}

func (db *Cockroachdb) SetPasswordStore(pwstore *data.PasswordStore) error {
	sql := fmt.Sprintf(`
	upsert into %s (uuid, password_store)
	values ($1, $2);
	`, db.TablePasswordStore)
	_, err := db.conn.Exec(db.context, sql, pwstore.Uuid, pwstore.PasswordStore)
	if err != nil {
		log.Printf("cannot upsert\nreason:\n%s", err)
		return err
	}
	return nil
}
