package app

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"wb-internship-l0/config"
	"wb-internship-l0/internal/broker"
	v1 "wb-internship-l0/internal/controller/http/v1"
	database "wb-internship-l0/internal/lib/pg"
	"wb-internship-l0/internal/repository"
	"wb-internship-l0/internal/service"
	"wb-internship-l0/pkg/cache"
	"wb-internship-l0/pkg/logger"
)

func Run() {
	// Config init
	cfg := config.MustLoad()

	// Logger init
	log := logger.NewZap(cfg.Env)
	defer log.Sync()

	// Context
	ctx, cancel := context.WithCancel(context.Background())

	// Database init
	pg := database.NewPostgres(ctx, log, cfg.PgDSN)

	// Cache init
	memoryCache := cache.NewMemoryCache()

	// Repositories init
	repositories := repository.NewRepositories(pg)

	// Services init
	deps := service.ServicesDependencies{
		Log:   log,
		Cache: memoryCache,
		Repos: repositories,
	}
	services := service.NewServices(deps)

	// Restore cache
	err := services.Order.LoadOrdersToCache(ctx)
	if err != nil {
		log.Warn("Failed to restore cache",
			zap.Error(err),
		)
	}

	// Channel for signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Broker init
	kafka := broker.NewKafkaConsumer(log, services, []string{cfg.Kafka.Host}, cfg.Kafka.Topic)
	go func() {
		if err := kafka.Listen(ctx); err != nil {
			log.Error("Kafka consumer stopped with error",
				zap.Error(err),
			)
		}
	}()

	// Router init
	app := fiber.New(fiber.Config{
		AppName: "WB-INTERNSHIP-L0",
	})
	v1.InitRouter(log, app, services)
	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Error("Fiber server error",
				zap.Error(err),
			)
		}
	}()

	// Graceful shutdown
	<-quit
	log.Info("Shutdown signal received")
	cancel()

	kafka.Shutdown()

	if err := app.Shutdown(); err != nil {
		log.Error("Error shutting down Fiber",
			zap.Error(err),
		)
	} else {
		log.Info("Fiber server stopped")
	}

	log.Info("Gracefully stopped")

}
