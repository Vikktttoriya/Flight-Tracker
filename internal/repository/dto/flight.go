package dto

import "time"

type Flight struct {
	ID                 int64      `db:"id"`
	FlightNumber       string     `db:"flight_number"`
	AirlineCode        string     `db:"airline_code"`
	DepartureAirport   string     `db:"departure_airport"`
	ArrivalAirport     string     `db:"arrival_airport"`
	ScheduledDeparture time.Time  `db:"scheduled_departure"`
	ScheduledArrival   time.Time  `db:"scheduled_arrival"`
	ActualDeparture    *time.Time `db:"actual_departure"`
	ActualArrival      *time.Time `db:"actual_arrival"`
	Status             string     `db:"status"`
	CreatedAt          time.Time  `db:"created_at"`
}
