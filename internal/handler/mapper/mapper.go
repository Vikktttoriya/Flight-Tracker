package mapper

import (
	"github.com/Vikktttoriya/flight-tracker/internal/domain/flight"
	"github.com/Vikktttoriya/flight-tracker/internal/domain/stats"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/dto"
)

func FlightToResponse(f *flight.Flight) *dto.FlightResponse {
	return &dto.FlightResponse{
		ID:                 f.ID,
		FlightNumber:       f.FlightNumber,
		AirlineCode:        f.AirlineCode,
		DepartureAirport:   f.DepartureAirport,
		ArrivalAirport:     f.ArrivalAirport,
		ScheduledDeparture: f.ScheduledDeparture,
		ScheduledArrival:   f.ScheduledArrival,
		ActualDeparture:    f.ActualDeparture,
		ActualArrival:      f.ActualArrival,
		Status:             string(f.Status),
	}
}

func FlightsToResponse(flights []*flight.Flight) []*dto.FlightResponse {
	result := make([]*dto.FlightResponse, 0, len(flights))
	for _, f := range flights {
		result = append(result, FlightToResponse(f))
	}
	return result
}

func StatsToResponse(s *stats.Stats) *dto.StatsResponse {
	return &dto.StatsResponse{
		UsersCount:   s.UsersCount,
		FlightsCount: s.FlightsCount,
		CollectedAt:  s.CollectedAt,
	}
}
