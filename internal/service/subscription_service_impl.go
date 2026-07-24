package service

import (
	"context"
	"log/slog"
	"time"

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

const monthYearLayout = "01-2006"

func parseMonthYear(value string) (time.Time, error) {
	return time.Parse(monthYearLayout, value)
}

func parseOptionalMonthYear(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	t, err := parseMonthYear(value)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (s *subscriptionService) Create(ctx context.Context, req model.CreateSubscriptionRequest) (*model.Subscription, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		slog.WarnContext(ctx, "invalid user id", "user_id", req.UserID, "error", err)
		return nil, err
	}

	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		slog.WarnContext(ctx, "invalid start date", "start_date", req.StartDate, "error", err)
		return nil, err
	}

	endDate, err := parseOptionalMonthYear(req.EndDate)
	if err != nil {
		slog.WarnContext(ctx, "invalid end date", "end_date", req.EndDate, "error", err)
		return nil, err
	}

	subscription := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := s.repo.Create(ctx, subscription); err != nil {
		slog.ErrorContext(ctx, "repository create failed", "user_id", userID, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "subscription created",
		"id", subscription.ID,
		"user_id", userID,
		"service_name", req.ServiceName,
	)
	return subscription, nil
}

func (s *subscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	subscription, err := s.repo.GetByID(ctx, id)
	if err != nil {
		slog.WarnContext(ctx, "subscription not found in repository", "id", id, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "subscription fetched", "id", id)
	return subscription, nil
}

func (s *subscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		slog.ErrorContext(ctx, "repository delete failed", "id", id, "error", err)
		return err
	}

	slog.InfoContext(ctx, "subscription deleted", "id", id)
	return nil
}

func (s *subscriptionService) List(ctx context.Context) ([]model.Subscription, error) {
	subscriptions, err := s.repo.List(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "repository list failed", "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "subscriptions listed", "count", len(subscriptions))
	return subscriptions, nil
}

func (s *subscriptionService) Update(ctx context.Context, id uuid.UUID, req model.UpdateSubscriptionRequest) error {
	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		slog.WarnContext(ctx, "invalid start date", "id", id, "start_date", req.StartDate, "error", err)
		return err
	}

	endDate, err := parseOptionalMonthYear(req.EndDate)
	if err != nil {
		slog.WarnContext(ctx, "invalid end date", "id", id, "end_date", req.EndDate, "error", err)
		return err
	}

	subscription := &model.Subscription{
		ID:          id,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := s.repo.Update(ctx, subscription); err != nil {
		slog.ErrorContext(ctx, "repository update failed", "id", id, "error", err)
		return err
	}

	slog.InfoContext(ctx, "subscription updated", "id", id)
	return nil
}

func (s *subscriptionService) Calculate(ctx context.Context, req model.CalculateRequest) (*model.CalculateResponse, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		slog.WarnContext(ctx, "invalid user id for calculation", "user_id", req.UserID, "error", err)
		return nil, err
	}

	from, err := time.Parse(monthYearLayout, req.From)
	if err != nil {
		slog.WarnContext(ctx, "invalid from date", "from", req.From, "error", err)
		return nil, err
	}

	to, err := time.Parse(monthYearLayout, req.To)
	if err != nil {
		slog.WarnContext(ctx, "invalid to date", "to", req.To, "error", err)
		return nil, err
	}

	now := time.Now()
	if to.After(now) {
		to = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	}

	if from.After(to) {
		slog.InfoContext(ctx, "calculation skipped: from after to",
			"user_id", userID,
			"service_name", req.ServiceName,
		)
		return &model.CalculateResponse{Total: 0}, nil
	}

	subscriptions, err := s.repo.GetSubscriptionsForCalculation(ctx, userID, req.ServiceName)
	if err != nil {
		slog.ErrorContext(ctx, "repository calculation query failed",
			"user_id", userID,
			"service_name", req.ServiceName,
			"error", err,
		)
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

		total += monthsBetween(start, end) * subscription.Price
	}

	slog.InfoContext(ctx, "subscription total calculated",
		"user_id", userID,
		"service_name", req.ServiceName,
		"from", req.From,
		"to", req.To,
		"subscriptions_count", len(subscriptions),
		"total", total,
	)

	return &model.CalculateResponse{Total: total}, nil
}

func monthsBetween(start, end time.Time) int {
	months := 0

	current := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC)
	last := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, time.UTC)

	for !current.After(last) {
		months++
		current = current.AddDate(0, 1, 0)
	}

	return months
}
