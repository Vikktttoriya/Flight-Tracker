package protected

import (
	"encoding/json"
	"net/http"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/user"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/dto"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/http/error_handler"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/http/middleware"
	"github.com/Vikktttoriya/flight-tracker/internal/service"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.ListUsers(r.Context())
	if err != nil {
		error_handler.HandleServiceError(w, err)
		return
	}

	dto.RespondJSON(w, http.StatusOK, users)
}

func (h *UserHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {
	login := chi.URLParam(r, "login")

	var req dto.ChangeUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.RespondError(w, http.StatusBadRequest, "bad_request", "invalid json")
		return
	}

	currentUser, _ := middleware.UserLogin(r.Context())

	updated, err := h.userService.ChangeUserRole(
		r.Context(),
		login,
		user.Role(req.Role),
		currentUser,
	)
	if err != nil {
		error_handler.HandleServiceError(w, err)
		return
	}

	dto.RespondJSON(w, http.StatusOK, updated)
}
