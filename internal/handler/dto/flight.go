package dto

import "time"

type FlightResponse struct {
	ID                 int64      `json:"id"`
	FlightNumber       string     `json:"flight_number"`
	AirlineCode        string     `json:"airline_code"`
	DepartureAirport   string     `json:"departure_airport"`
	ArrivalAirport     string     `json:"arrival_airport"`
	ScheduledDeparture time.Time  `json:"scheduled_departure"`
	ScheduledArrival   time.Time  `json:"scheduled_arrival"`
	ActualDeparture    *time.Time `json:"actual_departure,omitempty"`
	ActualArrival      *time.Time `json:"actual_arrival,omitempty"`
	Status             string     `json:"status"`
}

type CreateFlightRequest struct {
	FlightNumber       string `json:"flight_number"`
	AirlineCode        string `json:"airline_code"`
	DepartureAirport   string `json:"departure_airport"`
	ArrivalAirport     string `json:"arrival_airport"`
	ScheduledDeparture string `json:"scheduled_departure"`
	ScheduledArrival   string `json:"scheduled_arrival"`
}

type UpdateFlightStatusRequest struct {
	Status string `json:"status"`
}
