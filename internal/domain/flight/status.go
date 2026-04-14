package flight

import (
	"time"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/domain_errors"
)

type Status string

const (
	StatusScheduled Status = "scheduled"
	StatusCheckIn   Status = "check_in"
	StatusBoarding  Status = "boarding"
	StatusDeparted  Status = "departed"
	StatusArrived   Status = "arrived"
	StatusCanceled  Status = "canceled"
)

var allowedTransitions = map[Status][]Status{
	StatusScheduled: {StatusCheckIn, StatusCanceled},
	StatusCheckIn:   {StatusBoarding, StatusCanceled},
	StatusBoarding:  {StatusDeparted, StatusCanceled},
	StatusDeparted:  {StatusArrived},
	StatusArrived:   {},
	StatusCanceled:  {},
}

func (f *Flight) CanChangeStatus(to Status) bool {
	for _, allowed := range allowedTransitions[f.Status] {
		if allowed == to {
			return true
		}
	}
	return false
}

func (f *Flight) ChangeStatus(to Status, now time.Time) error {
	if !f.CanChangeStatus(to) {
		return domain_errors.ErrInvalidStatusTransition
	}

	f.Status = to

	switch to {
	case StatusDeparted:
		f.ActualDeparture = &now
	case StatusArrived:
		f.ActualArrival = &now
	}

	return nil
}
