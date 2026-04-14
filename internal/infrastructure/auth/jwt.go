package auth

import (
	"errors"
	"time"

	"github.com/Vikktttoriya/flight-tracker/internal/config"
	"github.com/golang-jwt/jwt/v5"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/user"
)

type Claims struct {
	Login string    `json:"login"`
	Role  user.Role `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(cfg config.JWTConfig) *JWTManager {
	return &JWTManager{
		secret: []byte(cfg.Secret),
		ttl:    cfg.TTL,
	}
}

func (j *JWTManager) GenerateToken(login string, role user.Role) (string, error) {
	now := time.Now()

	claims := Claims{
		Login: login,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(j.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *JWTManager) ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return j.secret, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
