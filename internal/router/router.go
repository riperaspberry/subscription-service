package router
import (
	"github.com/gin-gonic/gin"
	"github.com/riperaspberry/subscription-service/internal/handler"
)
func SetupRouter(
	subscriptionHandler *handler.SubscriptionHandler,
) *gin.Engine {

	r := gin.Default()

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