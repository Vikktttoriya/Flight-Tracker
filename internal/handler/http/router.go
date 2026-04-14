package http

import (
	"net/http"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/user"
	mwr "github.com/Vikktttoriya/flight-tracker/internal/handler/http/middleware"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/http/protected"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/http/public"
	"github.com/Vikktttoriya/flight-tracker/internal/infrastructure/auth"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Params struct {
	JWTManager             auth.JWTManager
	AuthHandler            *public.AuthHandler
	FlightHandler          *public.FlightHandler
	ProtectedFlightHandler *protected.FlightHandler
	StatsHandler           *public.StatsHandler
	UserHandler            *protected.UserHandler
}

func NewRouter(params Params) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", public.Health)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", params.AuthHandler.Register)
		r.Post("/login", params.AuthHandler.Login)
	})

	r.Get("/flights", params.FlightHandler.List)
	r.Get("/flights/{id}", params.FlightHandler.GetByID)

	r.Get("/stats", params.StatsHandler.GetLatest)

	r.Group(func(r chi.Router) {
		r.Use(mwr.JWT(params.JWTManager))

		r.Group(func(r chi.Router) {
			r.Use(mwr.RequireRole(user.RoleAdmin, user.RoleDispatcher))
			r.Post("/flights", params.ProtectedFlightHandler.Create)
			r.Patch("/flights/{id}/status", params.ProtectedFlightHandler.UpdateStatus)
		})

		r.Group(func(r chi.Router) {
			r.Use(mwr.RequireRole(user.RoleAdmin))
			r.Get("/admin/users", params.UserHandler.ListUsers)
			r.Patch("/admin/users/{login}/role", params.UserHandler.ChangeRole)
		})
	})

	return r
}
