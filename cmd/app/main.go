package main

import (
	"log/slog"
	_ "github.com/riperaspberry/subscription-service/docs"
	"github.com/riperaspberry/subscription-service/internal/config"
	"github.com/riperaspberry/subscription-service/internal/database"
	"github.com/riperaspberry/subscription-service/internal/handler"
	"github.com/riperaspberry/subscription-service/internal/logger"
	"github.com/riperaspberry/subscription-service/internal/repository"
	"github.com/riperaspberry/subscription-service/internal/router"
	"github.com/riperaspberry/subscription-service/internal/service"
)

func main() {
	logger.Init()

	cfg := config.Load()

	slog.Info("starting subscription service", "port", cfg.AppPort)

	db, err := database.New(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return
	}
	defer db.Close()

	repo := repository.NewSubscriptionRepository(db)
	subscriptionService := service.NewSubscriptionService(repo)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)
	r := router.SetupRouter(subscriptionHandler)

	slog.Info("server started", "addr", ":"+cfg.AppPort)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		slog.Error("server stopped with error", "error", err)
	}
}
