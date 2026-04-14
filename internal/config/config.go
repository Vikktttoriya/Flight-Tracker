package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App    AppConfig
	DB     DBConfig
	JWT    JWTConfig
	Logger LoggerConfig
	Worker WorkerConfig
}

type AppConfig struct {
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret string
	TTL    time.Duration
}

type LoggerConfig struct {
	Level    string
	FilePath string
	Format   string
	MaxAge   time.Duration
}

type WorkerConfig struct {
	StatsCollectionInterval time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	appPort := getEnv("APP_PORT", "8080")

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "flight_tracker")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	jwtSecret := getEnv("JWT_SECRET", "u8x/A?D(G+KbPeShVmYq3t6w9z+C-E)H@McQfTjWnZr4u7x!A%D*G-JaNdRgUkXp")
	jwtTTLHours := hoursToDuration(getEnvAsInt("JWT_TTL_HOURS", 24))

	logLevel := getEnv("LOG_LEVEL", "info")
	logFilePath := getEnv("LOG_FILE_PATH", "./logs/app.log")
	logMaxAge := hoursToDuration(getEnvAsInt("LOG_MAX_AGE_HOURS", 720))

	statsInterval := minutesToDuration(getEnvAsInt("WORKER_STATS_COLLECTION_INTERVAL_MINUTES", 2))

	return &Config{
		App: AppConfig{
			Port: appPort,
		},
		DB: DBConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			SSLMode:  dbSSLMode,
		},
		JWT: JWTConfig{
			Secret: jwtSecret,
			TTL:    jwtTTLHours,
		},
		Logger: LoggerConfig{
			Level:    logLevel,
			FilePath: logFilePath,
			MaxAge:   logMaxAge,
		},
		Worker: WorkerConfig{
			StatsCollectionInterval: statsInterval,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

func getEnvAsInt(key string, defaultVal int) int {
	if val := os.Getenv(key); strings.TrimSpace(val) != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func hoursToDuration(hours int) time.Duration {
	return time.Duration(hours) * time.Hour
}

func minutesToDuration(minutes int) time.Duration {
	return time.Duration(minutes) * time.Minute
}
