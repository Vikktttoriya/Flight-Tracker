package middleware

import (
	"net/http"
	"strings"

	"github.com/Vikktttoriya/flight-tracker/internal/handler/dto"
	"github.com/Vikktttoriya/flight-tracker/internal/infrastructure/auth"

	"go.uber.org/zap"
)

func JWT(jwtManager auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				dto.RespondError(w, http.StatusUnauthorized, "unauthorized", "missing authorization header")
				return
			}

			parts := strings.Split(header, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				dto.RespondError(w, http.StatusUnauthorized, "unauthorized", "invalid authorization header")
				return
			}

			claims, err := jwtManager.ParseToken(parts[1])
			if err != nil {
				zap.L().Warn("invalid jwt", zap.Error(err))
				dto.RespondError(w, http.StatusUnauthorized, "unauthorized", "invalid token")
				return
			}

			ctx := WithUser(r.Context(), claims.Login, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
