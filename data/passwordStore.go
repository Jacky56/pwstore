package data

import "github.com/google/uuid"

type PasswordStore struct {
	Uuid          uuid.UUID         `json:"uuid" db:"uuid"`
	PasswordStore map[string]string `json:"password_store" db:"password_store"`
}
