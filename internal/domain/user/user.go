package user

import "time"

type User struct {
	Login        string
	PasswordHash string
	Role         Role
	CreatedAt    time.Time
}

func NewUser(login string, passwordHash string) *User {
	return &User{
		Login:        login,
		PasswordHash: passwordHash,
		Role:         RolePassenger,
		CreatedAt:    time.Now(),
	}
}
