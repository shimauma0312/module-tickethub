package gorm

import (
	"context"
	"errors"
	"fmt"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"gorm.io/gorm"
)

// MilestoneRepository はGORMベースのMilestoneリポジトリ実装
type MilestoneRepository struct {
	db *gorm.DB
}

// NewMilestoneRepository は新しいMilestoneRepositoryを作成します
func NewMilestoneRepository(db *gorm.DB) (*MilestoneRepository, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &MilestoneRepository{db: db}, nil
}

// Create は新しいMilestoneを作成します
func (r *MilestoneRepository) Create(ctx context.Context, milestone *models.Milestone) error {
	if err := r.db.WithContext(ctx).Create(milestone).Error; err != nil {
		return fmt.Errorf("failed to create milestone: %w", err)
	}
	return nil
}

// GetByID はIDによってMilestoneを取得します
func (r *MilestoneRepository) GetByID(ctx context.Context, id int64) (*models.Milestone, error) {
	var milestone models.Milestone
	if err := r.db.WithContext(ctx).First(&milestone, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("milestone with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get milestone by id: %w", err)
	}
	return &milestone, nil
}

// List は条件に一致するMilestoneの一覧を取得します
func (r *MilestoneRepository) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Milestone, int, error) {
	var milestones []*models.Milestone
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Milestone{})

	for key, value := range filter {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count milestones: %w", err)
	}

	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Order("due_date ASC").Find(&milestones).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list milestones: %w", err)
	}

	return milestones, int(total), nil
}

// Update は既存のMilestoneを更新します
func (r *MilestoneRepository) Update(ctx context.Context, milestone *models.Milestone) error {
	if err := r.db.WithContext(ctx).Save(milestone).Error; err != nil {
		return fmt.Errorf("failed to update milestone: %w", err)
	}
	return nil
}

// Delete はMilestoneを削除します
func (r *MilestoneRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&models.Milestone{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete milestone: %w", err)
	}
	return nil
}

// Search はMilestoneの全文検索を行います (簡易的なLIKE検索の例)
func (r *MilestoneRepository) Search(ctx context.Context, query string, page, limit int) ([]*models.Milestone, int, error) {
	var milestones []*models.Milestone
	var total int64

	searchQuery := "%" + query + "%"
	dbQuery := r.db.WithContext(ctx).Model(&models.Milestone{}).
		Where("title LIKE ? OR description LIKE ?", searchQuery, searchQuery)

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count milestones for search: %w", err)
	}

	offset := (page - 1) * limit
	if err := dbQuery.Limit(limit).Offset(offset).Order("due_date ASC").Find(&milestones).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search milestones: %w", err)
	}

	return milestones, int(total), nil
}
