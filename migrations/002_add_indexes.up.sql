CREATE INDEX idx_flights_flight_number ON flights(flight_number);
CREATE INDEX idx_flights_status ON flights(status);

CREATE INDEX idx_statistics_created_at_desc ON statistics(created_at DESC);