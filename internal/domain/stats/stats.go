package stats

import "time"

type Stats struct {
	ID           int64
	UsersCount   int
	FlightsCount int
	CollectedAt  time.Time
}

func New(users, flights int) *Stats {
	return &Stats{
		UsersCount:   users,
		FlightsCount: flights,
		CollectedAt:  time.Now(),
	}
}
