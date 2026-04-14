package middleware

import (
	"net/http"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/user"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/dto"
)

func RequireRole(roles ...user.Role) func(http.Handler) http.Handler {
	allowed := make(map[user.Role]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			role, ok := UserRole(r.Context())
			if !ok {
				dto.RespondError(w, http.StatusForbidden, "forbidden", "role not found")
				return
			}

			if _, ok := allowed[role]; !ok {
				dto.RespondError(w, http.StatusForbidden, "forbidden", "access denied")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
