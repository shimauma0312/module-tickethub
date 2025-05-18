package repositories

import (
	"context"

	"github.com/shimauma0312/module-tickethub/backend/models"
)

// IssueRepository はIssue関連のデータベース操作を抽象化するインターフェース
type IssueRepository interface {
	// Create は新しいIssueを作成します
	Create(ctx context.Context, issue *models.Issue) error
	// GetByID はIDによってIssueを取得します
	GetByID(ctx context.Context, id int64) (*models.Issue, error)
	// List は条件に一致するIssueの一覧を取得します
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Issue, int, error)
	// Update は既存のIssueを更新します
	Update(ctx context.Context, issue *models.Issue) error
	// Delete はIssueを削除します
	Delete(ctx context.Context, id int64) error
	// Search はIssueの全文検索を行います
	Search(ctx context.Context, query string, page, limit int) ([]*models.Issue, int, error)
}

// RepositoryFactory はDBタイプに応じたリポジトリのインスタンスを生成するインターフェース
type RepositoryFactory interface {
	// NewIssueRepository はIssueRepositoryの新しいインスタンスを生成します
	NewIssueRepository() (IssueRepository, error)
	// Close はデータベース接続をクローズします
	Close() error
}
