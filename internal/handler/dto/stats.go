package dto

import "time"

type StatsResponse struct {
	UsersCount   int       `json:"users_count"`
	FlightsCount int       `json:"flights_count"`
	CollectedAt  time.Time `json:"collected_at"`
}
