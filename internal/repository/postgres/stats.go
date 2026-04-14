package postgres

import (
	"context"
	"errors"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/stats"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/dto"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/mapper"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	sq "github.com/Masterminds/squirrel"
)

type StatsRepository struct {
	db *pgxpool.Pool
	sb sq.StatementBuilderType
}

func NewStatsRepository(db *pgxpool.Pool) *StatsRepository {
	return &StatsRepository{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *StatsRepository) Save(ctx context.Context, s *stats.Stats) (*stats.Stats, error) {
	statsDTO := mapper.StatsToDTO(s)

	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "stats"),
		zap.String("operation", "save"),
	)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Error("Failed to begin transaction", zap.Error(err))
		return nil, err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("Transaction rollback failed", zap.Error(err))
		}
	}()

	query, args, err := r.sb.
		Insert("statistics").
		Columns(
			"total_users",
			"total_flights",
			"collected_at",
		).
		Values(
			statsDTO.UsersCount,
			statsDTO.FlightsCount,
			statsDTO.CollectedAt,
		).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		log.Error("Failed to build SAVE stats query", zap.Error(err))
		return nil, err
	}

	err = tx.QueryRow(ctx, query, args...).Scan(&s.ID)
	if err != nil {
		log.Error("Failed to save stats", zap.Error(err))
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction", zap.Error(err))
		return nil, err
	}

	log.Info("Stats saved successfully")

	return s, nil
}

func (r *StatsRepository) GetLatest(ctx context.Context) (*stats.Stats, error) {
	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "stats"),
		zap.String("operation", "get latest"),
	)

	query, args, err := r.sb.
		Select(
			"id",
			"total_users",
			"total_flights",
			"collected_at",
		).
		From("statistics").
		OrderBy("collected_at DESC").
		Limit(1).
		ToSql()

	if err != nil {
		log.Error("Failed to build GET latest stats query", zap.Error(err))
		return nil, err
	}

	var d dto.Stats
	row := r.db.QueryRow(ctx, query, args...)
	if err := row.Scan(
		&d.ID,
		&d.UsersCount,
		&d.FlightsCount,
		&d.CollectedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("No statistics found")
			return nil, nil
		}
		log.Error("Failed to scan stats", zap.Error(err))
		return nil, err
	}

	log.Info("Latest stats retrieved successfully", zap.Int("users_count", d.UsersCount), zap.Int("flights_count", d.FlightsCount), zap.Time("collected_at", d.CollectedAt))

	return mapper.StatsFromDTO(&d), nil
}
