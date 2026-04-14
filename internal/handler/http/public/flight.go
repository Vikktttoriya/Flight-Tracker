package public

import (
	"net/http"
	"strconv"

	"github.com/Vikktttoriya/flight-tracker/internal/handler/dto"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/http/error_handler"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/mapper"
	"github.com/Vikktttoriya/flight-tracker/internal/service"
	"github.com/go-chi/chi/v5"
)

type FlightHandler struct {
	flightService *service.FlightService
}

func NewFlightHandler(service *service.FlightService) *FlightHandler {
	return &FlightHandler{flightService: service}
}

func (h *FlightHandler) List(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	flights, err := h.flightService.List(r.Context(), offset, limit)
	if err != nil {
		error_handler.HandleServiceError(w, err)
		return
	}

	dto.RespondJSON(w, http.StatusOK, mapper.FlightsToResponse(flights))
}

func (h *FlightHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		dto.RespondError(w, http.StatusBadRequest, "invalid_id", "invalid flight id")
		return
	}

	flight, err := h.flightService.GetByID(r.Context(), id)
	if err != nil {
		error_handler.HandleServiceError(w, err)
		return
	}

	dto.RespondJSON(w, http.StatusOK, mapper.FlightToResponse(flight))
}
