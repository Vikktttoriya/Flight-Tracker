package postgres

import (
	"context"
	"errors"

	"github.com/Vikktttoriya/flight-tracker/internal/domain/user"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/db_errors"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/dto"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/mapper"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	sq "github.com/Masterminds/squirrel"
)

type UserRepository struct {
	db *pgxpool.Pool
	sb sq.StatementBuilderType
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
		sb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) (*user.User, error) {
	userDTO := mapper.UserToDTO(u)

	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "user"),
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
		Insert("users").
		Columns(
			"login",
			"password_hash",
			"role",
			"created_at",
		).
		Values(
			userDTO.Login,
			userDTO.PasswordHash,
			userDTO.Role,
			userDTO.CreatedAt,
		).
		Suffix("RETURNING login").
		ToSql()

	if err != nil {
		log.Error("Failed to build CREATE user query", zap.Error(err))
		return nil, err
	}

	err = tx.QueryRow(ctx, query, args...).Scan(&u.Login)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			log.Warn("Unique constraint violation",
				zap.String("pg_error_code", pgErr.Code),
				zap.String("pg_error_message", pgErr.Message),
			)
			return nil, db_errors.ErrDuplicateUser
		}

		log.Error("Failed to create user", zap.Error(err))
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction", zap.Error(err))
		return nil, err
	}

	log.Info("User created successfully")
	return u, nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*user.User, error) {
	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "user"),
		zap.String("operation", "get_by_login"),
	)

	query, args, err := r.sb.
		Select(
			"login",
			"password_hash",
			"role",
			"created_at",
		).
		From("users").
		Where(sq.Eq{"login": login}).
		ToSql()

	if err != nil {
		log.Error("Failed to build SELECT user query", zap.Error(err))
		return nil, err
	}

	var d dto.User
	row := r.db.QueryRow(ctx, query, args...)
	if err := row.Scan(
		&d.Login,
		&d.PasswordHash,
		&d.Role,
		&d.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("User not found")
			return nil, db_errors.ErrUserNotFound
		}
		log.Error("Failed to scan user", zap.Error(err))
		return nil, err
	}

	log.Info("User retrieved successfully")
	return mapper.UserFromDTO(&d), nil
}

func (r *UserRepository) List(ctx context.Context) ([]*user.User, error) {
	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "user"),
		zap.String("operation", "list"),
	)

	query, args, err := r.sb.
		Select(
			"login",
			"password_hash",
			"role",
			"created_at",
		).
		From("users").
		OrderBy("created_at DESC").
		ToSql()

	if err != nil {
		log.Error("Failed to build LIST users query", zap.Error(err))
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		log.Error("Failed to execute LIST users query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var users []*user.User
	for rows.Next() {
		var d dto.User
		if err := rows.Scan(
			&d.Login,
			&d.PasswordHash,
			&d.Role,
			&d.CreatedAt,
		); err != nil {
			log.Error("Failed to scan user row", zap.Error(err))
			return nil, err
		}
		users = append(users, mapper.UserFromDTO(&d))
	}

	if err := rows.Err(); err != nil {
		log.Error("Error iterating user rows", zap.Error(err))
		return nil, err
	}

	log.Info("Users listed successfully", zap.Int("count", len(users)))
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) (*user.User, error) {
	userDTO := mapper.UserToDTO(u)

	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "user"),
		zap.String("operation", "update"),
		zap.String("login", u.Login),
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
		Update("users").
		Set("password_hash", userDTO.PasswordHash).
		Set("role", userDTO.Role).
		Where(sq.Eq{"login": userDTO.Login}).
		ToSql()

	if err != nil {
		log.Error("Failed to build UPDATE user query", zap.Error(err))
		return nil, err
	}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		log.Error("Failed to execute UPDATE user query", zap.Error(err))
		return nil, err
	}

	if result.RowsAffected() == 0 {
		log.Warn("User not found for update")
		return nil, db_errors.ErrUserNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction", zap.Error(err))
		return nil, err
	}

	log.Info("User updated successfully")
	return u, nil
}

func (r *UserRepository) Count(ctx context.Context) (int, error) {
	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "user"),
		zap.String("operation", "count"),
	)

	query, args, err := r.sb.
		Select("COUNT(*)").
		From("users").
		ToSql()

	if err != nil {
		log.Error("Failed to build COUNT users query", zap.Error(err))
		return 0, err
	}

	var count int
	row := r.db.QueryRow(ctx, query, args...)
	if err := row.Scan(&count); err != nil {
		log.Error("Failed to scan user count", zap.Error(err))
		return 0, err
	}

	return count, nil
}

func (r *UserRepository) Delete(ctx context.Context, login string) error {
	log := zap.L().With(
		zap.String("layer", "repository"),
		zap.String("component", "user"),
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
		Delete("users").
		Where(sq.Eq{"login": login}).
		ToSql()

	if err != nil {
		log.Error("Failed to build DELETE user query", zap.Error(err))
		return err
	}

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		log.Error("Failed to execute DELETE user query", zap.Error(err))
		return err
	}

	if result.RowsAffected() == 0 {
		log.Warn("User not found for deletion")
		return db_errors.ErrUserNotFound
	}

	if err := tx.Commit(ctx); err != nil {
		log.Error("Failed to commit transaction", zap.Error(err))
		return err
	}

	log.Info("User deleted successfully")
	return nil
}
