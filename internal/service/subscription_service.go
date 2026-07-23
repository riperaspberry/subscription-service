package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/riperaspberry/subscription-service/internal/model"
)

type SubscriptionService interface {
	Create(ctx context.Context, req model.CreateSubscriptionRequest) (*model.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateSubscriptionRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context) ([]model.Subscription, error)
}
