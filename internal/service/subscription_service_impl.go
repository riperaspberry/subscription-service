package service

import (
	"time"
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


	startDate, err := time.Parse(
		"01-2006",
		req.StartDate,
	)

	if err != nil {
		return nil, err
	}


	subscription := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
	}


	err = s.repo.Create(
		ctx,
		subscription,
	)

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
func (s *subscriptionService) Calculate(ctx context.Context, req model.CalculateRequest) (*model.CalculateResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, err
	}

	from, err := time.Parse("01-2006", req.From)
	if err != nil {
		return nil, err
	}

	to, err := time.Parse("01-2006", req.To)
	if err != nil {
		return nil, err
	}
	
	now := time.Now()

if to.After(now) {
	to = time.Date(
		now.Year(),
		now.Month(),
		1,
		0, 0, 0, 0,
		time.UTC,
	)
}

if from.After(to) {
	return &model.CalculateResponse{
		Total: 0,
	}, nil
}

	subscriptions, err := s.repo.GetSubscriptionsForCalculation(
		ctx,
		userID,
		req.ServiceName,
	)

	if err != nil {
		return nil, err
	}

	total := 0

	for _, subscription := range subscriptions {

		start := subscription.StartDate

var end time.Time

if subscription.EndDate == nil {
	end = to
} else {
	end = *subscription.EndDate
}


		if start.Before(from) {
			start = from
		}

		if end.After(to) {
			end = to
		}

		if start.After(end) {
			continue
		}

		months := monthsBetween(start, end)

		total += months * subscription.Price
	}

	return &model.CalculateResponse{
		Total: total,
	}, nil
}
func monthsBetween(start, end time.Time) int {

	months := 0

	current := time.Date(
		start.Year(),
		start.Month(),
		1,
		0,0,0,0,
		time.UTC,
	)

	last := time.Date(
		end.Year(),
		end.Month(),
		1,
		0,0,0,0,
		time.UTC,
	)

	for !current.After(last) {
		months++

		current = current.AddDate(0,1,0)
	}

	return months
}