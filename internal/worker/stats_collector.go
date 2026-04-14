package worker

import (
	"context"
	"sync"
	"time"

	"github.com/Vikktttoriya/flight-tracker/internal/config"
	"github.com/Vikktttoriya/flight-tracker/internal/domain/flight"
	"github.com/Vikktttoriya/flight-tracker/internal/domain/stats"
	"github.com/Vikktttoriya/flight-tracker/internal/domain/user"

	"go.uber.org/zap"
)

type StatsCollector struct {
	userRepo   user.Repository
	flightRepo flight.Repository
	statsRepo  stats.Repository
	interval   time.Duration

	wg     sync.WaitGroup
	stopCh chan struct{}
}

func NewStatsCollector(userRepo user.Repository, flightRepo flight.Repository, statsRepo stats.Repository, cfg config.WorkerConfig) *StatsCollector {
	return &StatsCollector{
		userRepo:   userRepo,
		flightRepo: flightRepo,
		statsRepo:  statsRepo,
		interval:   cfg.StatsCollectionInterval,
		stopCh:     make(chan struct{}),
	}
}

func (w *StatsCollector) Start(ctx context.Context) {
	log := zap.L().With(zap.String("worker", "stats_collector"))
	log.Info("stats worker started", zap.Duration("interval", w.interval))

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.run(ctx)
	}()
}

func (w *StatsCollector) run(ctx context.Context) {
	log := zap.L().With(zap.String("worker", "stats_collector"))

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("stats worker stopped by context")
			return

		case <-w.stopCh:
			log.Info("stats worker stopped by signal")
			return

		case <-ticker.C:
			w.collect(ctx)
		}
	}
}

func (w *StatsCollector) collect(ctx context.Context) {
	zap.L().Info("stats collection started")

	usersCount, err := w.userRepo.Count(ctx)
	if err != nil {
		zap.L().Error("failed to count users", zap.Error(err))
		return
	}

	flightsCount, err := w.flightRepo.Count(ctx)
	if err != nil {
		zap.L().Error("failed to count flights", zap.Error(err))
		return
	}

	s := stats.New(usersCount, flightsCount)

	_, err = w.statsRepo.Save(ctx, s)
	if err != nil {
		zap.L().Error("failed to save stats", zap.Error(err))
		return
	}

	zap.L().Info("stats collected successfully", zap.Int("users", usersCount), zap.Int("flights", flightsCount))
}

func (w *StatsCollector) Stop() {
	close(w.stopCh)

	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		zap.L().Info("stats worker stopped gracefully")
	case <-time.After(5 * time.Second):
		zap.L().Warn("stats worker stop timeout")
	}
}
