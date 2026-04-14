package public

import (
	"encoding/json"
	"net/http"

	"github.com/Vikktttoriya/flight-tracker/internal/handler/dto"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/http/error_handler"
	"github.com/Vikktttoriya/flight-tracker/internal/service"

	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthHandler(auth *service.AuthService, user *service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: auth,
		userService: user,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.RespondError(w, http.StatusBadRequest, "bad_request", "invalid json")
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.Login, req.Password)
	if err != nil {
		zap.L().Warn("register failed", zap.Error(err))
		error_handler.HandleServiceError(w, err)
		return
	}

	dto.RespondJSON(w, http.StatusCreated, map[string]string{
		"login": user.Login,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.RespondError(w, http.StatusBadRequest, "bad_request", "invalid json")
		return
	}

	token, err := h.authService.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		error_handler.HandleServiceError(w, err)
		return
	}

	dto.RespondJSON(w, http.StatusOK, dto.AuthResponse{Token: token})
}
