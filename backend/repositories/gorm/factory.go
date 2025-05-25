package gorm

import (
	"github.com/shimauma0312/module-tickethub/backend/repositories"
	"gorm.io/gorm"
)

// RepositoryFactory はGORMを使用したリポジトリファクトリの実装
type RepositoryFactory struct {
	db *gorm.DB
}

// NewRepositoryFactory は新しいGORM用RepositoryFactoryを作成します
func NewRepositoryFactory(db *gorm.DB) *RepositoryFactory {
	return &RepositoryFactory{
		db: db,
	}
}

// NewIssueRepository はGORM用IssueRepositoryを作成します
func (f *RepositoryFactory) NewIssueRepository() (repositories.IssueRepository, error) {
	return NewIssueRepository(f.db), nil
}

// 他のリポジトリメソッドは必要に応じて追加
