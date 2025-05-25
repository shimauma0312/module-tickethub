package gorm

import (
	"context"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"gorm.io/gorm"
)

type pushSubscriptionRepository struct {
	db *gorm.DB
}

func NewPushSubscriptionRepository(db *gorm.DB) *pushSubscriptionRepository {
	return &pushSubscriptionRepository{db: db}
}

func (r *pushSubscriptionRepository) Create(ctx context.Context, subscription *models.PushSubscription) error {
	return r.db.WithContext(ctx).Create(subscription).Error
}

func (r *pushSubscriptionRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.PushSubscription, error) {
	var subscriptions []*models.PushSubscription
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&subscriptions).Error
	return subscriptions, err
}

func (r *pushSubscriptionRepository) DeleteByEndpoint(ctx context.Context, endpoint string) error {
	return r.db.WithContext(ctx).Where("endpoint = ?", endpoint).Delete(&models.PushSubscription{}).Error
}

func (r *pushSubscriptionRepository) DeleteByUserID(ctx context.Context, userID int64) error {
	return r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.PushSubscription{}).Error
}
