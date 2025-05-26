package repositories

import (
	"context"

	"github.com/shimauma0312/module-tickethub/backend/models"
)

// RepositoryRepository はリポジトリの管理を行うインターフェース
type RepositoryRepository interface {
	// Create は新しいリポジトリを作成します
	Create(ctx context.Context, repo *models.Repository) error
	// Update は既存のリポジトリを更新します
	Update(ctx context.Context, repo *models.Repository) error
	// Delete はリポジトリを削除します
	Delete(ctx context.Context, id int64) error
	// GetByID はIDによってリポジトリを取得します
	GetByID(ctx context.Context, id int64) (*models.Repository, error)
	// GetByName は名前によってリポジトリを取得します
	GetByName(ctx context.Context, name string) (*models.Repository, error)
	// List は条件に一致するリポジトリの一覧を取得します
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Repository, int64, error)
}
