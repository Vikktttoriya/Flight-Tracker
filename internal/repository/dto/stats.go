package dto

import "time"

type Stats struct {
	ID           int64     `db:"id"`
	UsersCount   int       `db:"total_users"`
	FlightsCount int       `db:"total_flights"`
	CollectedAt  time.Time `db:"collected_at"`
}
