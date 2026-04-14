package service

import (
	"context"
	"testing"

	"github.com/Vikktttoriya/flight-tracker/internal/config"
	"github.com/Vikktttoriya/flight-tracker/internal/domain/user"
	"github.com/Vikktttoriya/flight-tracker/internal/domain/user/mocks"
	"github.com/Vikktttoriya/flight-tracker/internal/infrastructure/auth"
	"github.com/Vikktttoriya/flight-tracker/internal/service/service_errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthService_Login_Success(t *testing.T) {
	repo := new(mocks.Repository)
	jwt := auth.NewJWTManager(config.JWTConfig{
		Secret: "secret",
		TTL:    1,
	})

	password := "password"
	hash, _ := auth.HashPassword(password)

	repo.On("GetByLogin", mock.Anything, "john").
		Return(&user.User{
			Login:        "john",
			PasswordHash: hash,
			Role:         user.RolePassenger,
		}, nil)

	service := NewAuthService(repo, *jwt)

	token, err := service.Login(context.Background(), "john", password)

	require.NoError(t, err)
	require.NotEmpty(t, token)
	repo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	repo := new(mocks.Repository)
	jwt := auth.NewJWTManager(config.JWTConfig{Secret: "secret", TTL: 1})

	repo.On("GetByLogin", mock.Anything, "john").
		Return(&user.User{
			Login:        "john",
			PasswordHash: "hash",
			Role:         user.RolePassenger,
		}, nil)

	service := NewAuthService(repo, *jwt)

	_, err := service.Login(context.Background(), "john", "wrong")

	require.Error(t, err)

	svcErr := err.(*service_errors.Error)
	require.Equal(t, service_errors.CodeInvalidCredentials, svcErr.Code)
}
