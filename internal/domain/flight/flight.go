package flight

import (
	"time"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/domain_errors"
)

type Flight struct {
	ID                 int64
	FlightNumber       string
	AirlineCode        string
	DepartureAirport   string
	ArrivalAirport     string
	ScheduledDeparture time.Time
	ScheduledArrival   time.Time
	ActualDeparture    *time.Time
	ActualArrival      *time.Time
	Status             Status
	CreatedAt          time.Time
}

func NewFlight(
	flightNumber string,
	airlineCode string,
	departureAirport string,
	arrivalAirport string,
	scheduledDeparture time.Time,
	scheduledArrival time.Time,
) (*Flight, error) {

	if scheduledDeparture.Before(time.Now()) {
		return nil, domain_errors.ErrFlightInPast
	}

	return &Flight{
		FlightNumber:       flightNumber,
		AirlineCode:        airlineCode,
		DepartureAirport:   departureAirport,
		ArrivalAirport:     arrivalAirport,
		ScheduledDeparture: scheduledDeparture,
		ScheduledArrival:   scheduledArrival,
		Status:             StatusScheduled,
		CreatedAt:          time.Now(),
	}, nil
}
