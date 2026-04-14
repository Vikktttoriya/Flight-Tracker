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

type AuthService struct {
	userRepo   user.Repository
	jwtManager auth.JWTManager
}

func NewAuthService(userRepo user.Repository, jwtManager auth.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	log := zap.L().With(
		zap.String("layer", "service"),
		zap.String("component", "auth"),
		zap.String("operation", "login"),
	)

	u, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, db_errors.ErrUserNotFound) {
			return "", &service_errors.Error{
				Code:    service_errors.CodeNotFound,
				Message: db_errors.ErrUserNotFound.Error(),
			}
		}
		return "", &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "database error",
			Err:     err,
		}
	}

	if !auth.CheckPassword(u.PasswordHash, password) {
		log.Warn("Invalid password provided")
		return "", &service_errors.Error{
			Code:    service_errors.CodeInvalidCredentials,
			Message: "wrong password",
		}
	}

	token, err := s.jwtManager.GenerateToken(u.Login, u.Role)
	if err != nil {
		log.Error("Failed to generate token", zap.Error(err))
		return "", &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "failed to generate token",
			Err:     err,
		}
	}

	log.Info("User logged in successfully")
	return token, nil
}
