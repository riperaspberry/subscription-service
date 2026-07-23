package repository

import (
	"context"

	"github.com/riperaspberry/subscription-service/internal/model"

	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *model.Subscription) error

	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)

	Update(ctx context.Context, subscription *model.Subscription) error

	Delete(ctx context.Context, id uuid.UUID) error

	List(ctx context.Context) ([]model.Subscription, error)
}
