package postgres

import (
	"context"
	"errors"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/flight"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/db_errors"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/dto"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/mapper"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	sq "github.com/Masterminds/squirrel"
)

type FlightRepository struct {
	db *pgxpool.Pool
	sb sq.StatementBuilderType
}

func NewFlightRepository(db *pgxpool.Pool) *FlightRepository {
	return &FlightRepository{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *FlightRepository) Create(ctx context.Context, f *flight.Flight) (*flight.Flight, error) {
	flightDTO := mapper.FlightToDTO(f)

	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "flight"),
		zap.String("operation", "create"),
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
		Insert("flights").
		Columns(
			"flight_number",
			"airline_code",
			"departure_airport",
			"arrival_airport",
			"scheduled_departure",
			"scheduled_arrival",
			"status",
			"created_at",
		).
		Values(
			flightDTO.FlightNumber,
			flightDTO.AirlineCode,
			flightDTO.DepartureAirport,
			flightDTO.ArrivalAirport,
			flightDTO.ScheduledDeparture,
			flightDTO.ScheduledArrival,
			flightDTO.Status,
			flightDTO.CreatedAt,
		).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		log.Error("Failed to build CREATE flight query", zap.Error(err))
		return nil, err
	}

	err = tx.QueryRow(ctx, query, args...).Scan(&f.ID)
	if err != nil {
		log.Error("Failed to create flight", zap.Error(err))
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction", zap.Error(err))
		return nil, err
	}

	log.Info("Flight created successfully",
		zap.Int64("flight_id", f.ID),
	)

	return f, nil
}

func (r *FlightRepository) GetByID(ctx context.Context, id int64) (*flight.Flight, error) {
	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "flight"),
		zap.String("operation", "get by id"),
	)

	query, args, err := r.sb.
		Select(
			"id",
			"flight_number",
			"airline_code",
			"departure_airport",
			"arrival_airport",
			"scheduled_departure",
			"scheduled_arrival",
			"actual_departure",
			"actual_arrival",
			"status",
			"created_at",
		).
		From("flights").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		log.Error("Failed to build SELECT flight query", zap.Error(err))
		return nil, err
	}

	var d dto.Flight
	row := r.db.QueryRow(ctx, query, args...)
	if err := row.Scan(
		&d.ID,
		&d.FlightNumber,
		&d.AirlineCode,
		&d.DepartureAirport,
		&d.ArrivalAirport,
		&d.ScheduledDeparture,
		&d.ScheduledArrival,
		&d.ActualDeparture,
		&d.ActualArrival,
		&d.Status,
		&d.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("Flight not found")
			return nil, db_errors.ErrFlightNotFound
		}
		log.Error("Failed to scan flight", zap.Error(err))
		return nil, err
	}

	log.Info("Flight retrieved successfully")
	return mapper.FlightFromDTO(&d), nil
}

func (r *FlightRepository) List(ctx context.Context, offset, limit int) ([]*flight.Flight, error) {
	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "flight"),
		zap.String("operation", "list"),
	)

	query, args, err := r.sb.
		Select(
			"id",
			"flight_number",
			"airline_code",
			"departure_airport",
			"arrival_airport",
			"scheduled_departure",
			"scheduled_arrival",
			"actual_departure",
			"actual_arrival",
			"status",
			"created_at",
		).
		From("flights").
		OrderBy("scheduled_departure").
		Offset(uint64(offset)).
		Limit(uint64(limit)).
		ToSql()

	if err != nil {
		log.Error("Failed to build LIST flights query", zap.Error(err))
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("Failed to execute LIST flights query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var flights []*flight.Flight
	for rows.Next() {
		var d dto.Flight
		if err := rows.Scan(
			&d.ID,
			&d.FlightNumber,
			&d.AirlineCode,
			&d.DepartureAirport,
			&d.ArrivalAirport,
			&d.ScheduledDeparture,
			&d.ScheduledArrival,
			&d.ActualDeparture,
			&d.ActualArrival,
			&d.Status,
			&d.CreatedAt,
		); err != nil {
			log.Error("Failed to scan flight row", zap.Error(err))
			return nil, err
		}
		flights = append(flights, mapper.FlightFromDTO(&d))
	}

	if err := rows.Err(); err != nil {
		log.Error("Error iterating flight rows", zap.Error(err))
		return nil, err
	}

	log.Info("Flights listed successfully", zap.Int("count", len(flights)))
	return flights, nil
}

func (r *FlightRepository) Update(ctx context.Context, flight *flight.Flight) (*flight.Flight, error) {
	flightDTO := mapper.FlightToDTO(flight)

	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "flight"),
		zap.String("operation", "update"),
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
		Update("flights").
		Set("status", flightDTO.Status).
		Set("actual_departure", flightDTO.ActualDeparture).
		Set("actual_arrival", flightDTO.ActualArrival).
		Where(sq.Eq{"id": flightDTO.ID}).
		ToSql()

	if err != nil {
		log.Error("Failed to build UPDATE flight query", zap.Error(err))
		return nil, err
	}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		log.Error("Failed to execute UPDATE flight query", zap.Error(err))
		return nil, err
	}

	if result.RowsAffected() == 0 {
		log.Warn("Flight not found for update")
		return nil, db_errors.ErrFlightNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction", zap.Error(err))
		return nil, err
	}

	log.Info("Flight updated successfully")
	return flight, nil
}

func (r *FlightRepository) Count(ctx context.Context) (int, error) {
	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "flight"),
		zap.String("operation", "count"),
	)
	log.Debug("Counting total flights")

	query, args, err := r.sb.
		Select("COUNT(*)").
		From("flights").
		ToSql()

	if err != nil {
		log.Error("Failed to build COUNT flights query", zap.Error(err))
		return 0, err
	}

	var count int
	row := r.db.QueryRow(ctx, query, args...)
	if err := row.Scan(&count); err != nil {
		log.Error("Failed to scan flight count", zap.Error(err))
		return 0, err
	}

	return count, nil
}

func (r *FlightRepository) Delete(ctx context.Context, id int64) error {
	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "flight"),
		zap.String("operation", "delete"),
	)

	tx, err := r.db.Begin(ctx)
	if err != nil {
		log.Error("Failed to begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Error("Transaction rollback failed", zap.Error(err))
		}
	}()

	query, args, err := r.sb.
		Delete("flights").
		Where(sq.Eq{"id": id}).
		ToSql()

	if err != nil {
		log.Error("Failed to build DELETE flight query", zap.Error(err))
		return err
	}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		log.Error("Failed to execute DELETE flight query", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		log.Warn("Flight not found for delete")
		return db_errors.ErrFlightNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction", zap.Error(err))
		return err
	}

	log.Info("Flight deleted successfully")
	return nil
}
