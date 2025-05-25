package gorm

import (
	"context"
	"errors"
	"fmt"

	"github.com/shimauma0312/module-tickethub/models"
	"gorm.io/gorm"
)

// LabelRepository はGORMベースのLabelリポジトリ実装
type LabelRepository struct {
	db *gorm.DB
}

// NewLabelRepository は新しいLabelRepositoryを作成します
func NewLabelRepository(db *gorm.DB) (*LabelRepository, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &LabelRepository{db: db}, nil
}

// Create は新しいLabelを作成します
func (r *LabelRepository) Create(ctx context.Context, label *models.Label) error {
	if err := r.db.WithContext(ctx).Create(label).Error; err != nil {
		return fmt.Errorf("failed to create label: %w", err)
	}
	return nil
}

// GetByID はIDによってLabelを取得します
func (r *LabelRepository) GetByID(ctx context.Context, id int64) (*models.Label, error) {
	var label models.Label
	if err := r.db.WithContext(ctx).First(&label, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("label with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get label by id: %w", err)
	}
	return &label, nil
}

// GetByName は名前とタイプによってLabelを取得します
func (r *LabelRepository) GetByName(ctx context.Context, name, labelType string) (*models.Label, error) {
	var label models.Label
	if err := r.db.WithContext(ctx).Where("name = ? AND type = ?", name, labelType).First(&label).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("label with name '%s' and type '%s' not found", name, labelType)
		}
		return nil, fmt.Errorf("failed to get label by name and type: %w", err)
	}
	return &label, nil
}

// List は条件に一致するLabelの一覧を取得します
func (r *LabelRepository) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Label, int, error) {
	var labels []*models.Label
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Label{})

	for key, value := range filter {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count labels: %w", err)
	}

	offset := (page - 1) * limit
	if err := query.Limit(limit).Offset(offset).Order("name ASC").Find(&labels).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list labels: %w", err)
	}

	return labels, int(total), nil
}

// Update は既存のLabelを更新します
func (r *LabelRepository) Update(ctx context.Context, label *models.Label) error {
	if err := r.db.WithContext(ctx).Save(label).Error; err != nil {
		return fmt.Errorf("failed to update label: %w", err)
	}
	return nil
}

// Delete はLabelを削除します
func (r *LabelRepository) Delete(ctx context.Context, id int64) error {
	if err := r.db.WithContext(ctx).Delete(&models.Label{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete label: %w", err)
	}
	return nil
}

// Search はLabelの全文検索を行います (簡易的なLIKE検索の例)
func (r *LabelRepository) Search(ctx context.Context, query string, page, limit int) ([]*models.Label, int, error) {
	var labels []*models.Label
	var total int64

	searchQuery := "%" + query + "%"
	dbQuery := r.db.WithContext(ctx).Model(&models.Label{}).
		Where("name LIKE ? OR description LIKE ?", searchQuery, searchQuery)

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count labels for search: %w", err)
	}

	offset := (page - 1) * limit
	if err := dbQuery.Limit(limit).Offset(offset).Order("name ASC").Find(&labels).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to search labels: %w", err)
	}

	return labels, int(total), nil
}
