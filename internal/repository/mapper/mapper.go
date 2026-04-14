package mapper

import (
	"github.com/Vikktttoriya/flight-tracker/internal/domain/flight"
	"github.com/Vikktttoriya/flight-tracker/internal/domain/stats"
	"github.com/Vikktttoriya/flight-tracker/internal/domain/user"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/dto"
)

func FlightToDTO(f *flight.Flight) *dto.Flight {
	return &dto.Flight{
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
		CreatedAt:          f.CreatedAt,
	}
}

func FlightFromDTO(d *dto.Flight) *flight.Flight {
	return &flight.Flight{
		ID:                 d.ID,
		FlightNumber:       d.FlightNumber,
		AirlineCode:        d.AirlineCode,
		DepartureAirport:   d.DepartureAirport,
		ArrivalAirport:     d.ArrivalAirport,
		ScheduledDeparture: d.ScheduledDeparture,
		ScheduledArrival:   d.ScheduledArrival,
		ActualDeparture:    d.ActualDeparture,
		ActualArrival:      d.ActualArrival,
		Status:             flight.Status(d.Status),
		CreatedAt:          d.CreatedAt,
	}
}

func UserToDTO(u *user.User) *dto.User {
	return &dto.User{
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
		Role:         string(u.Role),
		CreatedAt:    u.CreatedAt,
	}
}

func UserFromDTO(d *dto.User) *user.User {
	return &user.User{
		Login:        d.Login,
		PasswordHash: d.PasswordHash,
		Role:         user.Role(d.Role),
		CreatedAt:    d.CreatedAt,
	}
}

func StatsToDTO(s *stats.Stats) *dto.Stats {
	return &dto.Stats{
		ID:           s.ID,
		UsersCount:   s.UsersCount,
		FlightsCount: s.FlightsCount,
		CollectedAt:  s.CollectedAt,
	}
}

func StatsFromDTO(d *dto.Stats) *stats.Stats {
	return &stats.Stats{
		ID:           d.ID,
		UsersCount:   d.UsersCount,
		FlightsCount: d.FlightsCount,
		CollectedAt:  d.CollectedAt,
	}
}
