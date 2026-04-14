package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Vikktttoriya/flight-tracker/internal/config"
	rout "github.com/Vikktttoriya/flight-tracker/internal/handler/http"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/http/protected"
	"github.com/Vikktttoriya/flight-tracker/internal/handler/http/public"
	"github.com/Vikktttoriya/flight-tracker/internal/infrastructure/auth"
	"github.com/Vikktttoriya/flight-tracker/internal/infrastructure/database"
	"github.com/Vikktttoriya/flight-tracker/internal/infrastructure/logger"
	"github.com/Vikktttoriya/flight-tracker/internal/repository/postgres"
	"github.com/Vikktttoriya/flight-tracker/internal/service"
	"github.com/Vikktttoriya/flight-tracker/internal/worker"

	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()
	zap.ReplaceGlobals(log)

	log.Info("starting application")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	db, err := database.NewPostgres(cfg.DB)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	userRepo := postgres.NewUserRepository(db)
	flightRepo := postgres.NewFlightRepository(db)
	statsRepo := postgres.NewStatsRepository(db)

	jwtManager := auth.NewJWTManager(cfg.JWT)

	authService := service.NewAuthService(userRepo, *jwtManager)
	userService := service.NewUserService(userRepo)
	flightService := service.NewFlightService(flightRepo)
	statsService := service.NewStatsService(statsRepo)

	authHandler := public.NewAuthHandler(authService, userService)
	flightHandler := public.NewFlightHandler(flightService)
	protectedFlightHandler := protected.NewFlightHandler(flightService)
	statsHandler := public.NewStatsHandler(statsService)
	userHandler := protected.NewUserHandler(userService)

	statsWorker := worker.NewStatsCollector(
		userRepo,
		flightRepo,
		statsRepo,
		cfg.Worker,
	)
	statsWorker.Start(ctx)

	router := rout.NewRouter(rout.Params{
		JWTManager:             *jwtManager,
		AuthHandler:            authHandler,
		FlightHandler:          flightHandler,
		ProtectedFlightHandler: protectedFlightHandler,
		StatsHandler:           statsHandler,
		UserHandler:            userHandler,
	})

	server := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("http server started", zap.String("port", cfg.App.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("http server failed", zap.Error(err))
		}
	}()

	<-sigCh
	log.Info("shutdown signal received")

	cancel()
	statsWorker.Stop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("http server shutdown failed", zap.Error(err))
	}

	log.Info("application stopped gracefully")
}
