package dto

import "time"

type User struct {
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
}
