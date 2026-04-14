package service

import (
	"context"
	"testing"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/stats"
	"github.com/Vikktttoriya/flight-tracker/internal/domain/stats/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestStatsService_GetLatest_Success(t *testing.T) {
	repo := new(mocks.Repository)

	repo.On("GetLatest", mock.Anything).
		Return(&stats.Stats{FlightsCount: 10}, nil)

	service := NewStatsService(repo)

	stat, err := service.GetLatest(context.Background())

	require.NoError(t, err)
	require.Equal(t, 10, stat.FlightsCount)
	repo.AssertExpectations(t)
}

func TestStatsService_GetLatest_Empty(t *testing.T) {
	repo := new(mocks.Repository)

	repo.On("GetLatest", mock.Anything).
		Return(nil, nil)

	service := NewStatsService(repo)

	stat, err := service.GetLatest(context.Background())

	require.NoError(t, err)
	require.Nil(t, stat)
	repo.AssertExpectations(t)
}
