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

// NewUserRepository はGORM用UserRepositoryを作成します
func (f *RepositoryFactory) NewUserRepository() (repositories.UserRepository, error) {
	repo, err := NewUserRepository(f.db)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// NewUserSettingsRepository はGORM用UserSettingsRepositoryを作成します
func (f *RepositoryFactory) NewUserSettingsRepository() (repositories.UserSettingsRepository, error) {
	repo, err := NewUserSettingsRepository(f.db)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// NewLabelRepository はGORM用LabelRepositoryを作成します
func (f *RepositoryFactory) NewLabelRepository() (repositories.LabelRepository, error) {
	repo, err := NewLabelRepository(f.db)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// NewMilestoneRepository はGORM用MilestoneRepositoryを作成します
func (f *RepositoryFactory) NewMilestoneRepository() (repositories.MilestoneRepository, error) {
	repo, err := NewMilestoneRepository(f.db)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// NewDiscussionRepository はGORM用DiscussionRepositoryを作成します
func (f *RepositoryFactory) NewDiscussionRepository() (repositories.DiscussionRepository, error) {
	return NewDiscussionRepository(f.db), nil
}

// NewCommentRepository はGORM用CommentRepositoryを作成します
func (f *RepositoryFactory) NewCommentRepository() (repositories.CommentRepository, error) {
	return NewCommentRepository(f.db), nil
}

// NewReactionRepository はGORM用ReactionRepositoryを作成します
func (f *RepositoryFactory) NewReactionRepository() (repositories.ReactionRepository, error) {
	return NewReactionRepository(f.db), nil
}

// NewNotificationRepository はGORM用NotificationRepositoryを作成します
func (f *RepositoryFactory) NewNotificationRepository() (repositories.NotificationRepository, error) {
	return NewNotificationRepository(f.db), nil
}

// NewPushSubscriptionRepository はGORM用PushSubscriptionRepositoryを作成します
func (f *RepositoryFactory) NewPushSubscriptionRepository() (repositories.PushSubscriptionRepository, error) {
	return NewPushSubscriptionRepository(f.db), nil
}

// NewNotificationTemplateRepository はGORM用NotificationTemplateRepositoryを作成します
func (f *RepositoryFactory) NewNotificationTemplateRepository() (repositories.NotificationTemplateRepository, error) {
	return NewNotificationTemplateRepository(f.db), nil
}

// NewSystemSettingsRepository はGORM用SystemSettingsRepositoryを作成します
func (f *RepositoryFactory) NewSystemSettingsRepository() (repositories.SystemSettingsRepository, error) {
	return NewSystemSettingsRepository(f.db), nil
}

// NewActivityLogRepository はGORM用ActivityLogRepositoryを作成します
func (f *RepositoryFactory) NewActivityLogRepository() (repositories.ActivityLogRepository, error) {
	return NewActivityLogRepository(f.db), nil
}

// NewBackupRepository はGORM用BackupRepositoryを作成します
func (f *RepositoryFactory) NewBackupRepository() (repositories.BackupRepository, error) {
	return NewBackupRepository(f.db), nil
}

// NewRepositoryRepository はGORM用RepositoryRepositoryを作成します
func (f *RepositoryFactory) NewRepositoryRepository() (repositories.RepositoryRepository, error) {
	return NewRepositoryRepository(f.db), nil
}

// Close はデータベース接続をクローズします
func (f *RepositoryFactory) Close() error {
	db, err := f.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}
