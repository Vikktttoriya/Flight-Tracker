package middleware

import (
	"context"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/user"
)

type ctxKey string

const (
	ctxUserLogin ctxKey = "user_login"
	ctxUserRole  ctxKey = "user_role"
)

func WithUser(ctx context.Context, login string, role user.Role) context.Context {
	ctx = context.WithValue(ctx, ctxUserLogin, login)
	ctx = context.WithValue(ctx, ctxUserRole, role)
	return ctx
}

func UserLogin(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(ctxUserLogin).(string)
	return login, ok
}

func UserRole(ctx context.Context) (user.Role, bool) {
	role, ok := ctx.Value(ctxUserRole).(user.Role)
	return role, ok
}
