package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByLogin(ctx context.Context, login string) (*User, error)
	List(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Count(ctx context.Context) (int, error)
	Delete(ctx context.Context, login string) error
}
