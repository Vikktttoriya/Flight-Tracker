package service

import (
	"context"
	"testing"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/user"
	"github.com/Vikktttoriya/flight-tracker/internal/domain/user/mocks"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/db_errors"
	"github.com/Vikktttoriya/flight-tracker/internal/service/service_errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_CreateUser_Success(t *testing.T) {
	repo := new(mocks.Repository)

	repo.On("Create", mock.Anything, mock.Anything).
		Return(&user.User{
			Login: "john",
			Role:  user.RolePassenger,
		}, nil)

	service := NewUserService(repo)

	u, err := service.CreateUser(context.Background(), "john", "password")

	require.NoError(t, err)
	require.Equal(t, "john", u.Login)
	repo.AssertExpectations(t)
}

func TestUserService_CreateUser_Duplicate(t *testing.T) {
	repo := new(mocks.Repository)

	repo.On("Create", mock.Anything, mock.Anything).
		Return(nil, db_errors.ErrDuplicateUser)

	service := NewUserService(repo)

	_, err := service.CreateUser(context.Background(), "john", "password")

	require.Error(t, err)

	svcErr := err.(*service_errors.Error)
	require.Equal(t, service_errors.CodeAlreadyExists, svcErr.Code)
}

func TestUserService_ChangeRole_SelfModification(t *testing.T) {
	repo := new(mocks.Repository)
	service := NewUserService(repo)

	_, err := service.ChangeUserRole(
		context.Background(),
		"admin",
		user.RolePassenger,
		"admin",
	)

	require.Error(t, err)

	svcErr := err.(*service_errors.Error)
	require.Equal(t, service_errors.CodeSelfModification, svcErr.Code)
}
