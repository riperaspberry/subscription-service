package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/riperaspberry/subscription-service/internal/model"
	"github.com/riperaspberry/subscription-service/internal/service"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
}

func NewSubscriptionHandler(
	service service.SubscriptionService,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
	}
}
// Create subscription
// @Summary Create subscription
// @Description Creates a new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.CreateSubscriptionRequest true "Subscription data"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	var req model.CreateSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.WarnContext(c.Request.Context(), "invalid create subscription request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	subscription, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to create subscription", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.InfoContext(c.Request.Context(), "subscription created via http", "id", subscription.ID)
	c.JSON(http.StatusCreated, subscription)
}
// List subscriptions
// @Summary Get all subscriptions
// @Description Returns all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(c *gin.Context) {
	subscriptions, err := h.service.List(c.Request.Context())
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to list subscriptions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.InfoContext(c.Request.Context(), "subscriptions listed via http", "count", len(subscriptions))
	c.JSON(http.StatusOK, subscriptions)
}
// Get subscription by ID
// @Summary Get subscription
// @Description Returns subscription by UUID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription UUID"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		slog.WarnContext(c.Request.Context(), "invalid subscription id", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	subscription, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		slog.WarnContext(c.Request.Context(), "subscription not found", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.InfoContext(c.Request.Context(), "subscription fetched via http", "id", id)
	c.JSON(http.StatusOK, subscription)
}
// Delete subscription
// @Summary Delete subscription
// @Description Deletes subscription by UUID
// @Tags subscriptions
// @Param id path string true "Subscription UUID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		slog.WarnContext(c.Request.Context(), "invalid subscription id", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to delete subscription", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.InfoContext(c.Request.Context(), "subscription deleted via http", "id", id)
	c.Status(http.StatusNoContent)
}
// Update subscription
// @Summary Update subscription
// @Description Updates subscription information
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription UUID"
// @Param subscription body model.UpdateSubscriptionRequest true "Update data"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		slog.WarnContext(c.Request.Context(), "invalid subscription id", "id", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	var req model.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.WarnContext(c.Request.Context(), "invalid update subscription request", "id", id, "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, req); err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to update subscription", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.InfoContext(c.Request.Context(), "subscription updated via http", "id", id)
	c.Status(http.StatusOK)
}
// Calculate subscriptions cost
// @Summary Calculate total subscription cost
// @Description Calculates total spending for selected period
// @Tags subscriptions
// @Produce json
// @Param user_id query string true "User UUID"
// @Param service_name query string true "Service name"
// @Param from query string true "Start month MM-YYYY"
// @Param to query string true "End month MM-YYYY"
// @Success 200 {object} model.CalculateResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/calculate [get]
func (h *SubscriptionHandler) Calculate(c *gin.Context) {
	var req model.CalculateRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		slog.WarnContext(c.Request.Context(), "invalid calculate request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := h.service.Calculate(c.Request.Context(), req)
	if err != nil {
		slog.ErrorContext(c.Request.Context(), "failed to calculate subscriptions total",
			"user_id", req.UserID,
			"service_name", req.ServiceName,
			"error", err,
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	slog.InfoContext(c.Request.Context(), "subscription total calculated via http",
		"user_id", req.UserID,
		"service_name", req.ServiceName,
		"from", req.From,
		"to", req.To,
		"total", response.Total,
	)
	c.JSON(http.StatusOK, response)
}
