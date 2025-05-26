package gorm

import (
	"context"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
	"gorm.io/gorm"
)

// repositoryRepository はRepositoryRepositoryのGORM実装
type repositoryRepository struct {
	db *gorm.DB
}

// NewRepositoryRepository は新しいrepositoryRepositoryインスタンスを作成します
func NewRepositoryRepository(db *gorm.DB) repositories.RepositoryRepository {
	return &repositoryRepository{
		db: db,
	}
}

// Create は新しいリポジトリを作成します
func (r *repositoryRepository) Create(ctx context.Context, repo *models.Repository) error {
	return r.db.WithContext(ctx).Create(repo).Error
}

// Update は既存のリポジトリを更新します
func (r *repositoryRepository) Update(ctx context.Context, repo *models.Repository) error {
	return r.db.WithContext(ctx).Save(repo).Error
}

// Delete はリポジトリを削除します
func (r *repositoryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Repository{}, id).Error
}

// GetByID はIDによってリポジトリを取得します
func (r *repositoryRepository) GetByID(ctx context.Context, id int64) (*models.Repository, error) {
	var repo models.Repository
	if err := r.db.WithContext(ctx).First(&repo, id).Error; err != nil {
		return nil, err
	}
	return &repo, nil
}

// GetByName は名前によってリポジトリを取得します
func (r *repositoryRepository) GetByName(ctx context.Context, name string) (*models.Repository, error) {
	var repo models.Repository
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&repo).Error; err != nil {
		return nil, err
	}
	return &repo, nil
}

// List は条件に一致するリポジトリの一覧を取得します
func (r *repositoryRepository) List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Repository, int64, error) {
	var repos []*models.Repository
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Repository{})

	// フィルタリング条件を適用
	if filter != nil {
		for key, value := range filter {
			switch key {
			case "name":
				query = query.Where("name LIKE ?", "%"+value.(string)+"%")
			case "owner_id":
				query = query.Where("owner_id = ?", value)
			case "type":
				query = query.Where("type = ?", value)
			case "is_archived":
				query = query.Where("is_archived = ?", value)
			}
		}
	}

	// 総数を取得
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// ページネーション
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&repos).Error; err != nil {
		return nil, 0, err
	}

	return repos, total, nil
}
