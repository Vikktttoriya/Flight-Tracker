package service

import (
	"context"
	"errors"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/user"
	"github.com/Vikktttoriya/flight-tracker/internal/infrastructure/auth"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/db_errors"
	"github.com/Vikktttoriya/flight-tracker/internal/service/service_errors"

	"go.uber.org/zap"
)

type UserService struct {
	userRepo user.Repository
}

func NewUserService(userRepo user.Repository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx context.Context, login string, password string) (*user.User, error) {
	log := zap.L().With(
		zap.String("layer", "service"),
		zap.String("component", "user"),
		zap.String("operation", "create"),
	)

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		log.Error("Failed to hash password", zap.Error(err))
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "database error",
			Err:     err,
		}
	}

	u := user.NewUser(login, passwordHash)

	createdUser, err := s.userRepo.Create(ctx, u)
	if err != nil {
		if errors.Is(err, db_errors.ErrDuplicateUser) {
			log.Warn("User with this login already exists")
			return nil, &service_errors.Error{
				Code:    service_errors.CodeAlreadyExists,
				Message: db_errors.ErrDuplicateUser.Error(),
			}
		}
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "database error",
			Err:     err,
		}
	}

	return createdUser, nil
}

func (s *UserService) GetUserByLogin(ctx context.Context, login string) (*user.User, error) {
	u, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, db_errors.ErrUserNotFound) {
			return nil, &service_errors.Error{
				Code:    service_errors.CodeNotFound,
				Message: db_errors.ErrUserNotFound.Error(),
			}
		}
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "database error",
			Err:     err,
		}
	}

	return u, nil
}

func (s *UserService) ListUsers(ctx context.Context) ([]*user.User, error) {
	users, err := s.userRepo.List(ctx)
	if err != nil {
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "database error",
			Err:     err,
		}
	}

	return users, nil
}

func (s *UserService) ChangeUserRole(ctx context.Context, login string, role user.Role, currentUserLogin string) (*user.User, error) {
	log := zap.L().With(
		zap.String("layer", "service"),
		zap.String("component", "user"),
		zap.String("operation", "change role"),
	)

	if login == currentUserLogin {
		log.Warn("Attempt to change own role")
		return nil, &service_errors.Error{
			Code:    service_errors.CodeSelfModification,
			Message: "cannot change your own role",
		}
	}

	u, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, db_errors.ErrUserNotFound) {
			return nil, &service_errors.Error{
				Code:    service_errors.CodeNotFound,
				Message: db_errors.ErrUserNotFound.Error(),
			}
		}
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "database error",
			Err:     err,
		}
	}

	updatedUser := &user.User{
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
		Role:         role,
	}

	updatedUser, err = s.userRepo.Update(ctx, updatedUser)
	if err != nil {
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "database error",
			Err:     err,
		}
	}

	return updatedUser, nil
}
