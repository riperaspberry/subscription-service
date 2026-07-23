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
		subscriptions.GET("/:id", subscriptionHandler.GetByID)
		subscriptions.DELETE("/:id", subscriptionHandler.Delete)
		subscriptions.PUT("/:id", subscriptionHandler.Update)
	}

	return r
}