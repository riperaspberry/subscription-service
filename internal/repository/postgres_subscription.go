package repository

import (
	"context"
	"log/slog"

	"github.com/riperaspberry/subscription-service/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type subscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(ctx context.Context, subscription *model.Subscription) error {
	query := `
	INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at
	`
	err := r.db.QueryRow(ctx, query,
		subscription.ServiceName,
		subscription.Price,
		subscription.UserID,
		subscription.StartDate,
		subscription.EndDate,
	).Scan(&subscription.ID, &subscription.CreatedAt)
	if err != nil {
		slog.ErrorContext(ctx, "db insert subscription failed",
			"user_id", subscription.UserID,
			"service_name", subscription.ServiceName,
			"error", err,
		)
		return err
	}

	slog.InfoContext(ctx, "db insert subscription",
		"id", subscription.ID,
		"user_id", subscription.UserID,
	)
	return nil
}

func (r *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at FROM subscriptions WHERE id = $1`
	var subscription model.Subscription
	err := r.db.QueryRow(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.ServiceName,
		&subscription.Price,
		&subscription.UserID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.CreatedAt,
	)
	if err != nil {
		slog.WarnContext(ctx, "db select subscription by id failed", "id", id, "error", err)
		return nil, err
	}

	slog.InfoContext(ctx, "db select subscription by id", "id", id)
	return &subscription, nil
}

func (r *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	if _, err := r.db.Exec(ctx, query, id); err != nil {
		slog.ErrorContext(ctx, "db delete subscription failed", "id", id, "error", err)
		return err
	}

	slog.InfoContext(ctx, "db delete subscription", "id", id)
	return nil
}

func (r *subscriptionRepository) List(ctx context.Context) ([]model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at FROM subscriptions`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		slog.ErrorContext(ctx, "db list subscriptions failed", "error", err)
		return nil, err
	}
	defer rows.Close()

	var subscriptions []model.Subscription
	for rows.Next() {
		var subscription model.Subscription
		if err := rows.Scan(
			&subscription.ID,
			&subscription.ServiceName,
			&subscription.Price,
			&subscription.UserID,
			&subscription.StartDate,
			&subscription.EndDate,
			&subscription.CreatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "db scan subscription failed", "error", err)
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	slog.InfoContext(ctx, "db list subscriptions", "count", len(subscriptions))
	return subscriptions, nil
}

func (r *subscriptionRepository) Update(ctx context.Context, subscription *model.Subscription) error {
	query := `UPDATE subscriptions SET service_name = $1, price = $2 WHERE id = $3`
	if _, err := r.db.Exec(ctx, query,
		subscription.ServiceName,
		subscription.Price,
		subscription.ID,
	); err != nil {
		slog.ErrorContext(ctx, "db update subscription failed", "id", subscription.ID, "error", err)
		return err
	}

	slog.InfoContext(ctx, "db update subscription", "id", subscription.ID)
	return nil
}

func (r *subscriptionRepository) GetSubscriptionsForCalculation(ctx context.Context, userID uuid.UUID, serviceName string) ([]model.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date, created_at FROM subscriptions WHERE user_id = $1 AND service_name = $2`
	rows, err := r.db.Query(ctx, query, userID, serviceName)
	if err != nil {
		slog.ErrorContext(ctx, "db select subscriptions for calculation failed",
			"user_id", userID,
			"service_name", serviceName,
			"error", err,
		)
		return nil, err
	}
	defer rows.Close()

	var subscriptions []model.Subscription
	for rows.Next() {
		var subscription model.Subscription
		if err := rows.Scan(
			&subscription.ID,
			&subscription.ServiceName,
			&subscription.Price,
			&subscription.UserID,
			&subscription.StartDate,
			&subscription.EndDate,
			&subscription.CreatedAt,
		); err != nil {
			slog.ErrorContext(ctx, "db scan subscription for calculation failed", "error", err)
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	slog.InfoContext(ctx, "db select subscriptions for calculation",
		"user_id", userID,
		"service_name", serviceName,
		"count", len(subscriptions),
	)
	return subscriptions, nil
}
