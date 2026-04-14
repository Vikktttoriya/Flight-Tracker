package service

import (
	"context"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/stats"
	"github.com/Vikktttoriya/flight-tracker/internal/service/service_errors"

	"go.uber.org/zap"
)

type StatsService struct {
	statsRepo stats.Repository
}

func NewStatsService(statsRepo stats.Repository) *StatsService {
	return &StatsService{
		statsRepo: statsRepo,
	}
}

func (s *StatsService) GetLatest(ctx context.Context) (*stats.Stats, error) {
	log := zap.L().With(
		zap.String("layer", "service"),
		zap.String("component", "stats"),
		zap.String("operation", "get latest"),
	)

	stat, err := s.statsRepo.GetLatest(ctx)
	if err != nil {
		return nil, &service_errors.Error{
			Code:    service_errors.CodeInternal,
			Message: "database error",
			Err:     err,
		}
	}

	if stat == nil {
		log.Info("No statistics found")
		return nil, nil
	}

	return stat, nil
}
