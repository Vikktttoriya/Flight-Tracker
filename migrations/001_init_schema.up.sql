CREATE TABLE IF NOT EXISTS users (
                                     login TEXT PRIMARY KEY,
                                     password_hash TEXT NOT NULL,
                                     role TEXT NOT NULL CHECK (role IN ('admin', 'dispatcher', 'passenger')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS flights (
                                       id SERIAL PRIMARY KEY,
                                       flight_number TEXT NOT NULL,
                                       airline_code TEXT NOT NULL,
                                       departure_airport TEXT NOT NULL,
                                       arrival_airport TEXT NOT NULL,
                                       scheduled_departure TIMESTAMP WITH TIME ZONE NOT NULL,
                                       scheduled_arrival TIMESTAMP WITH TIME ZONE NOT NULL,
                                       actual_departure TIMESTAMP WITH TIME ZONE,
                                       actual_arrival TIMESTAMP WITH TIME ZONE,
                                       status TEXT NOT NULL CHECK (status IN ('scheduled', 'check-in', 'boarding', 'departed', 'arrived', 'canceled')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );

CREATE TABLE IF NOT EXISTS statistics (
                                          id SERIAL PRIMARY KEY,
                                          total_users INTEGER NOT NULL DEFAULT 0,
                                          total_flights INTEGER NOT NULL DEFAULT 0,
                                          created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
    );