package router

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/riperaspberry/subscription-service/internal/handler"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		slog.Info("http request",
			"method", c.Request.Method,
			"path", path,
			"status", c.Writer.Status(),
			"duration_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}

func SetupRouter(
	subscriptionHandler *handler.SubscriptionHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery(), requestLogger())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	subscriptions := r.Group("/subscriptions")
	{
		subscriptions.POST("", subscriptionHandler.Create)
		subscriptions.GET("", subscriptionHandler.List)
		subscriptions.GET("/calculate", subscriptionHandler.Calculate)
		subscriptions.GET("/:id", subscriptionHandler.GetByID)
		subscriptions.PUT("/:id", subscriptionHandler.Update)
		subscriptions.DELETE("/:id", subscriptionHandler.Delete)
	}

	return r
}
