package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/riperaspberry/subscription-service/internal/model"
	"github.com/riperaspberry/subscription-service/internal/repository"
)

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}
func (s *subscriptionService) Create(ctx context.Context, req model.CreateSubscriptionRequest) (*model.Subscription, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, err
	}
	subscription := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
	}
	err = s.repo.Create(ctx, subscription)

	if err != nil {
		return nil, err
	}

	return subscription, nil
}
func (s *subscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *subscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
func (s *subscriptionService) List(ctx context.Context) ([]model.Subscription, error) {
	return s.repo.List(ctx)
}
func (s *subscriptionService) Update(ctx context.Context, id uuid.UUID, req model.UpdateSubscriptionRequest) error {
	subscription := &model.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
	}

	return s.repo.Update(ctx, subscription)
}
