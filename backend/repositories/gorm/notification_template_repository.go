package gorm

import (
	"context"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"gorm.io/gorm"
)

type notificationTemplateRepository struct {
	db *gorm.DB
}

func NewNotificationTemplateRepository(db *gorm.DB) *notificationTemplateRepository {
	return &notificationTemplateRepository{db: db}
}

func (r *notificationTemplateRepository) Create(ctx context.Context, template *models.NotificationTemplate) error {
	return r.db.WithContext(ctx).Create(template).Error
}

func (r *notificationTemplateRepository) GetByType(ctx context.Context, templateType string) (*models.NotificationTemplate, error) {
	var template models.NotificationTemplate
	err := r.db.WithContext(ctx).Where("template_type = ?", templateType).First(&template).Error
	return &template, err
}

func (r *notificationTemplateRepository) GetAll(ctx context.Context) ([]*models.NotificationTemplate, error) {
	var templates []*models.NotificationTemplate
	err := r.db.WithContext(ctx).Find(&templates).Error
	return templates, err
}

func (r *notificationTemplateRepository) Update(ctx context.Context, template *models.NotificationTemplate) error {
	return r.db.WithContext(ctx).Save(template).Error
}

func (r *notificationTemplateRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.NotificationTemplate{}, id).Error
}
