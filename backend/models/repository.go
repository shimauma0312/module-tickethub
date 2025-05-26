package models

import (
	"time"

	"gorm.io/gorm"
)

// RepositoryType はリポジトリのタイプを表す型
type RepositoryType string

const (
	// PublicRepo は公開リポジトリを表す
	PublicRepo RepositoryType = "public"
	// PrivateRepo は非公開リポジトリを表す
	PrivateRepo RepositoryType = "private"
	// InternalRepo は組織内公開リポジトリを表す
	InternalRepo RepositoryType = "internal"
)

// Repository はリポジトリ情報を表す構造体
type Repository struct {
	ID          int64          `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        RepositoryType `json:"type"`
	OwnerID     int64          `json:"owner_id"`
	OwnerName   string         `json:"owner_name,omitempty"` // JOINなどで取得する場合用
	IsArchived  bool           `json:"is_archived"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// NewRepository は新しいRepositoryインスタンスを作成する
func NewRepository(name, description string, repoType RepositoryType, ownerID int64) *Repository {
	now := time.Now()
	return &Repository{
		Name:        name,
		Description: description,
		Type:        repoType,
		OwnerID:     ownerID,
		IsArchived:  false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update はリポジトリ情報を更新する
func (r *Repository) Update(name, description string, repoType RepositoryType) {
	if name != "" {
		r.Name = name
	}
	r.Description = description
	if repoType != "" {
		r.Type = repoType
	}
	r.UpdatedAt = time.Now()
}

// Archive はリポジトリをアーカイブ状態に設定する
func (r *Repository) Archive() {
	r.IsArchived = true
	r.UpdatedAt = time.Now()
}

// Unarchive はリポジトリのアーカイブを解除する
func (r *Repository) Unarchive() {
	r.IsArchived = false
	r.UpdatedAt = time.Now()
}

// IsPublic はリポジトリが公開状態かどうかを返す
func (r *Repository) IsPublic() bool {
	return r.Type == PublicRepo
}

// IsPrivate はリポジトリが非公開状態かどうかを返す
func (r *Repository) IsPrivate() bool {
	return r.Type == PrivateRepo
}

// IsInternal はリポジトリが組織内公開状態かどうかを返す
func (r *Repository) IsInternal() bool {
	return r.Type == InternalRepo
}

// AutoMigrateRepository はRepositoryテーブルのマイグレーションを実行します
func AutoMigrateRepository(db interface{}) error {
	if gormDB, ok := db.(*gorm.DB); ok {
		return gormDB.AutoMigrate(&Repository{})
	}
	return nil
}
