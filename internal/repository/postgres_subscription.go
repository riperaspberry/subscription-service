package repository

import (
	"context"

	"github.com/riperaspberry/subscription-service/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type subscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool,) SubscriptionRepository {
	return &subscriptionRepository{
		db: db,
	}
}
func (r *subscriptionRepository) Create(ctx context.Context, subscription *model.Subscription) error {
	query := `
	INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at
	`
	err := r.db.QueryRow(ctx, query, subscription.ServiceName, subscription.Price, subscription.UserID, subscription.StartDate, subscription.EndDate).Scan(&subscription.ID, &subscription.CreatedAt)
	return err
}
func (r *subscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error) {
	query := ` SELECT id, service_name, price, user_id, start_date, end_date, created_at FROM subscriptions WHERE id = $1`
	var subscription model.Subscription
	err := r.db.QueryRow(ctx, query, id).Scan(&subscription.ID, &subscription.ServiceName, &subscription.Price, &subscription.UserID, &subscription.StartDate, &subscription.EndDate, &subscription.CreatedAt)
	return &subscription, err
}
func (r *subscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := ` DELETE FROM subscriptions WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
func (r *subscriptionRepository) List(ctx context.Context) ([]model.Subscription, error) {
	query := ` SELECT id, service_name, price, user_id, start_date, end_date, created_at FROM subscriptions`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var subscriptions []model.Subscription
	for rows.Next() {
		var subscription model.Subscription
		err := rows.Scan(&subscription.ID, &subscription.ServiceName, &subscription.Price, &subscription.UserID, &subscription.StartDate, &subscription.EndDate, &subscription.CreatedAt)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}
func (r *subscriptionRepository) Update(ctx context.Context, subscription *model.Subscription) error {
	query := ` UPDATE subscriptions SET service_name = $1, price = $2, start_date = $3, end_date = $4 WHERE id = $5`
	_, err := r.db.Exec(ctx, query, subscription.ServiceName, subscription.Price, subscription.StartDate, subscription.EndDate, subscription.ID)
	return err
}