package public

import (
	"net/http"

	"github.com/Vikktttoriya/flight-tracker/internal/handler/dto"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/http/error_handler"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/mapper"
	"github.com/Vikktttoriya/flight-tracker/internal/service"
)

type StatsHandler struct {
	statsService *service.StatsService
}

func NewStatsHandler(service *service.StatsService) *StatsHandler {
	return &StatsHandler{statsService: service}
}

func (h *StatsHandler) GetLatest(w http.ResponseWriter, r *http.Request) {
	stats, err := h.statsService.GetLatest(r.Context())
	if err != nil {
		error_handler.HandleServiceError(w, err)
		return
	}

	if stats == nil {
		dto.RespondJSON(w, http.StatusOK, map[string]interface{}{})
		return
	}

	dto.RespondJSON(w, http.StatusOK, mapper.StatsToResponse(stats))
}
