package service

import (
	"context"
	"errors"
	"time"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/domain_errors"
	"github.com/Vikktttoriya/flight-tracker/internal/domain/flight"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/db_errors"
	"github.com/Vikktttoriya/flight-tracker/internal/service/service_errors"

	"go.uber.org/zap"
)

type FlightService struct {
	flightRepo flight.Repository
}

func NewFlightService(flightRepo flight.Repository) *FlightService {
	return &FlightService{
		flightRepo: flightRepo,
	}
}

func (s *FlightService) GetByID(ctx context.Context, id int64) (*flight.Flight, error) {
	log := zap.L().With(
		zap.String("layer", "service"),
		zap.String("component", "flight"),
		zap.String("operation", "get by id"),
	)

	if id <= 0 {
		log.Warn("invalid flight ID")
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInvalidArgument,
			Message: "invalid flight id",
		}
	}

	f, err := s.flightRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, db_errors.ErrFlightNotFound) {
			return nil, &service_errors.Error{
				Code:    service_errors.CodeNotFound,
				Message: db_errors.ErrFlightNotFound.Error(),
			}
		}
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "database error",
			Err:     err,
		}
	}

	return f, nil
}

func (s *FlightService) List(ctx context.Context, offset, limit int) ([]*flight.Flight, error) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	flights, err := s.flightRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "database error",
			Err:     err,
		}
	}

	return flights, nil
}

func (s *FlightService) CreateFlight(ctx context.Context, f *flight.Flight) (*flight.Flight, error) {

	createdFlight, err := s.flightRepo.Create(ctx, f)
	if err != nil {
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "database error",
			Err:     err,
		}
	}

	return createdFlight, nil
}

func (s *FlightService) UpdateFlightStatus(ctx context.Context, flightID int64, newStatus flight.Status) (*flight.Flight, error) {
	log := zap.L().With(
		zap.String("layer", "service"),
		zap.String("component", "flight"),
		zap.String("operation", "update"),
	)

	currentFlight, err := s.flightRepo.GetByID(ctx, flightID)
	if err != nil {
		if errors.Is(err, db_errors.ErrFlightNotFound) {
			return nil, &service_errors.Error{
				Code:    service_errors.CodeNotFound,
				Message: db_errors.ErrFlightNotFound.Error(),
			}
		}
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "database error",
			Err:     err,
		}
	}

	err = currentFlight.ChangeStatus(newStatus, time.Now())
	if err != nil {
		log.Error("Error changing status", zap.Error(err))
		if errors.Is(err, domain_errors.ErrInvalidStatusTransition) {
			return nil, &service_errors.Error{
				Code:    service_errors.CodeInvalidTransition,
				Message: domain_errors.ErrInvalidStatusTransition.Error(),
			}
		}
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "internal error",
			Err:     err,
		}
	}

	result, err := s.flightRepo.Update(ctx, currentFlight)
	if err != nil {
		if errors.Is(err, db_errors.ErrFlightNotFound) {
			return nil, &service_errors.Error{
				Code:    service_errors.CodeNotFound,
				Message: db_errors.ErrFlightNotFound.Error(),
			}
		}
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "database error",
			Err:     err,
		}
	}

	return result, nil
}
