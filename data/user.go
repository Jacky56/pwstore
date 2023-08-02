package data

import "github.com/google/uuid"

type User struct {
	Uuid     uuid.UUID `json:"uuid" db:"uuid"`
	Email    string    `json:"email" db:"email"`
	Password string    `json:"password" db:"password"`
	Provider string    `json:"provider" db:"provider"`
}
