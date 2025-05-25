package gorm

import (
	"context"
	"fmt"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"gorm.io/gorm"
)

type discussionRepository struct {
	db *gorm.DB
}

func NewDiscussionRepository(db *gorm.DB) *discussionRepository {
	return &discussionRepository{db: db}
}

func (r *discussionRepository) Create(ctx context.Context, discussion *models.Discussion) error {
	return r.db.WithContext(ctx).Create(discussion).Error
}

func (r *discussionRepository) GetByID(ctx context.Context, id int64) (*models.Discussion, error) {
	var discussion models.Discussion
	err := r.db.WithContext(ctx).First(&discussion, id).Error
	return &discussion, err
}

func (r *discussionRepository) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Discussion, int, error) {
	var discussions []*models.Discussion
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Discussion{})
	if filter != nil {
		query = query.Where(filter)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&discussions).Error
	return discussions, int(total), err
}

func (r *discussionRepository) Update(ctx context.Context, discussion *models.Discussion) error {
	return r.db.WithContext(ctx).Save(discussion).Error
}

func (r *discussionRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Discussion{}, id).Error
}

func (r *discussionRepository) Search(ctx context.Context, queryString string, page, limit int) ([]*models.Discussion, int, error) {
	// Implement search logic if necessary, or return not implemented
	// This is a placeholder
	var discussions []*models.Discussion
	var total int64
	// Example: r.db.WithContext(ctx).Where("title LIKE ?", "%"+queryString+"%").Find(&discussions)
	// Count total records for pagination
	// r.db.WithContext(ctx).Model(&models.Discussion{}).Where("title LIKE ?", "%"+queryString+"%").Count(&total)
	return discussions, int(total), gorm.ErrNotImplemented // Placeholder
}

// CountDiscussions は総Discussion数を取得します
func (r *discussionRepository) CountDiscussions(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Discussion{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count discussions: %w", err)
	}
	return count, nil
}

// CountOpenDiscussions はオープンなDiscussion数を取得します
func (r *discussionRepository) CountOpenDiscussions(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Discussion{}).Where("status = ?", "open").Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count open discussions: %w", err)
	}
	return count, nil
}
