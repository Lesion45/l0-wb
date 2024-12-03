package app

import (
	"context"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	// Context
	ctx, cancel := context.WithCancel(context.Background())

	// Database init
	log.Info("Database initialization...")
	pg := database.NewPostgres(ctx, log, cfg.PgDSN)
	log.Info("Database initialization: OK.")

	// Cache init
	log.Info("Cache initialization...")
	memoryCache := cache.NewMemoryCache()
	log.Info("Cache initialization: OK.")

	// Repositories init
	log.Info("Repository initialization...")
	repositories := repository.NewRepositories(pg)
	log.Info("Repository initialization: OK.")

	// Services init
	log.Info("Services initialization...")
	deps := service.ServicesDependencies{
		Log:   log,
		Cache: memoryCache,
		Repos: repositories,
	}
	services := service.NewServices(deps)
	log.Info("Services initialization: OK.")

	// Restore cache
	log.Info("Restoring cache...")
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
	log.Info("Kafka reader initialization...")
	kafka := broker.NewKafkaConsumer(log, services, []string{cfg.Kafka.Host}, cfg.Kafka.Topic)
	go func() {
		if err := kafka.Listen(ctx); err != nil {
			log.Error("Kafka consumer stopped with error",
				zap.Error(err),
			)
		}
	}()

	// Router init
	log.Info("Router initialization...")
	app := fiber.New(fiber.Config{
		AppName: "WB-INTERNSHIP-L0",
	})
	app.Use(recover.New())
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: log,
	}))
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

	err = kafka.Shutdown()
	if err != nil {
		log.Error("Failed to close Kafka")
	}

	if err := app.Shutdown(); err != nil {
		log.Error("Error shutting down Fiber",
			zap.Error(err),
		)
	} else {
		log.Info("Fiber server stopped")
	}

	log.Info("Gracefully stopped")

}
