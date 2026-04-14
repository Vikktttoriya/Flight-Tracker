package service

import (
	"context"
	"testing"
	"time"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/flight"
	"github.com/Vikktttoriya/flight-tracker/internal/domain/flight/mocks"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/db_errors"
	"github.com/Vikktttoriya/flight-tracker/internal/service/service_errors"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFlightService_GetByID_Success(t *testing.T) {
	repo := new(mocks.Repository)

	repo.On("GetByID", mock.Anything, int64(1)).
		Return(&flight.Flight{ID: 1}, nil)

	service := NewFlightService(repo)

	f, err := service.GetByID(context.Background(), 1)

	require.NoError(t, err)
	require.Equal(t, int64(1), f.ID)
}

func TestFlightService_GetByID_NotFound(t *testing.T) {
	repo := new(mocks.Repository)

	repo.On("GetByID", mock.Anything, int64(1)).
		Return(nil, db_errors.ErrFlightNotFound)

	service := NewFlightService(repo)

	_, err := service.GetByID(context.Background(), 1)

	require.Error(t, err)

	svcErr := err.(*service_errors.Error)
	require.Equal(t, service_errors.CodeNotFound, svcErr.Code)
}

func TestFlightService_UpdateStatus_InvalidTransition(t *testing.T) {
	repo := new(mocks.Repository)

	f := &flight.Flight{
		ID:     1,
		Status: flight.StatusArrived,
	}

	repo.On("GetByID", mock.Anything, int64(1)).
		Return(f, nil)

	service := NewFlightService(repo)

	_, err := service.UpdateFlightStatus(
		context.Background(),
		1,
		flight.StatusBoarding,
	)

	require.Error(t, err)

	svcErr := err.(*service_errors.Error)

	require.Equal(t, service_errors.CodeInvalidTransition, svcErr.Code)
}

func TestFlightService_CreateFlight_Success(t *testing.T) {
	repo := new(mocks.Repository)

	now := time.Now()
	flightToCreate, _ := flight.NewFlight("SU100", "010", "SVO", "UFA", now.Add(2*time.Hour), now.Add(5*time.Hour))

	createdFlight := &flight.Flight{
		ID:                 1,
		FlightNumber:       "SU100",
		AirlineCode:        "010",
		DepartureAirport:   "SVO",
		ArrivalAirport:     "UFA",
		ScheduledDeparture: now.Add(2 * time.Hour),
		ScheduledArrival:   now.Add(5 * time.Hour),
		Status:             flight.StatusScheduled,
		CreatedAt:          now,
	}

	repo.On("Create", mock.Anything, mock.MatchedBy(func(f *flight.Flight) bool {
		return f.FlightNumber == "SU100" &&
			f.DepartureAirport == "SVO" &&
			f.ArrivalAirport == "UFA" &&
			f.Status == flight.StatusScheduled
	})).Return(createdFlight, nil).Once()

	service := NewFlightService(repo)

	f, err := service.CreateFlight(context.Background(), flightToCreate)

	require.NoError(t, err)
	require.Equal(t, int64(1), f.ID)
	require.Equal(t, "SU100", f.FlightNumber)
	require.Equal(t, flight.StatusScheduled, f.Status)
	repo.AssertExpectations(t)
}
