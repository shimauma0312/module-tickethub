package postgres

import (
	"context"
	"database/sql"

	"github.com/shimauma0312/module-tickethub/backend/models"
)

// PushSubscriptionRepository はPostgreSQL固有のPushSubscriptionリポジトリ実装です
type PushSubscriptionRepository struct {
	db *sql.DB
}

// NewPushSubscriptionRepository は新しいPostgreSQL PushSubscriptionRepositoryを生成します
func NewPushSubscriptionRepository(db *sql.DB) *PushSubscriptionRepository {
	return &PushSubscriptionRepository{
		db: db,
	}
}

// Create は新しいPushSubscriptionを作成します
func (r *PushSubscriptionRepository) Create(ctx context.Context, subscription *models.PushSubscription) error {
	query := `
		INSERT INTO push_subscriptions (user_id, endpoint, p256dh, auth, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		subscription.UserID,
		subscription.Endpoint,
		subscription.P256dh,
		subscription.Auth,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	).Scan(&subscription.ID)

	return err
}

// GetByUserID はユーザーIDによってPushSubscriptionを取得します
func (r *PushSubscriptionRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.PushSubscription, error) {
	query := `
		SELECT id, user_id, endpoint, p256dh, auth, created_at, updated_at
		FROM push_subscriptions
		WHERE user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*models.PushSubscription
	for rows.Next() {
		subscription := &models.PushSubscription{}
		err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.Endpoint,
			&subscription.P256dh,
			&subscription.Auth,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, subscription)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

// DeleteByEndpoint はEndpointによってPushSubscriptionを削除します
func (r *PushSubscriptionRepository) DeleteByEndpoint(ctx context.Context, endpoint string) error {
	query := `DELETE FROM push_subscriptions WHERE endpoint = $1`
	_, err := r.db.ExecContext(ctx, query, endpoint)
	return err
}

// DeleteByUserID はユーザーIDによってPushSubscriptionを削除します
func (r *PushSubscriptionRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	query := `DELETE FROM push_subscriptions WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
