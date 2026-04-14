package protected

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/domain_errors"
	"github.com/Vikktttoriya/flight-tracker/internal/domain/flight"
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

func (h *FlightHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateFlightRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.RespondError(w, http.StatusBadRequest, "bad_request", "invalid json")
		return
	}

	dep, err := time.Parse(time.RFC3339, req.ScheduledDeparture)
	arr, err2 := time.Parse(time.RFC3339, req.ScheduledArrival)
	if err != nil || err2 != nil {
		dto.RespondError(w, http.StatusBadRequest, "bad_request", "invalid time format")
		return
	}

	f, err := flight.NewFlight(
		req.FlightNumber,
		req.AirlineCode,
		req.DepartureAirport,
		req.ArrivalAirport,
		dep,
		arr,
	)
	if err != nil {
		if errors.Is(err, domain_errors.ErrFlightInPast) {
			dto.RespondError(w, http.StatusConflict, "past_flight", "cannot create flight in the past")
		} else {
			error_handler.HandleServiceError(w, err)
		}
		return
	}

	created, err := h.flightService.CreateFlight(r.Context(), f)
	if err != nil {
		error_handler.HandleServiceError(w, err)
		return
	}

	dto.RespondJSON(w, http.StatusCreated, mapper.FlightToResponse(created))
}

func (h *FlightHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var req dto.UpdateFlightStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		dto.RespondError(w, http.StatusBadRequest, "bad_request", "invalid json")
		return
	}

	updated, err := h.flightService.UpdateFlightStatus(
		r.Context(),
		id,
		flight.Status(req.Status),
	)
	if err != nil {
		error_handler.HandleServiceError(w, err)
		return
	}

	dto.RespondJSON(w, http.StatusOK, mapper.FlightToResponse(updated))
}
