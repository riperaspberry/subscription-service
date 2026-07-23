package main

import (
	"log"

	"github.com/riperaspberry/subscription-service/internal/config"
	"github.com/riperaspberry/subscription-service/internal/database"
	"github.com/riperaspberry/subscription-service/internal/handler"
	"github.com/riperaspberry/subscription-service/internal/repository"
	"github.com/riperaspberry/subscription-service/internal/router"
	"github.com/riperaspberry/subscription-service/internal/service"
)

func main() {

	cfg := config.Load()

	db, err := database.New(cfg)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	repo := repository.NewSubscriptionRepository(db)

	subscriptionService := service.NewSubscriptionService(repo)

	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionService)

	r := router.SetupRouter(subscriptionHandler)

	err = r.Run(":" + cfg.AppPort)

	if err != nil {
		log.Fatal(err)
	}
}
