package services

import (
	"github.com/shimauma0312/module-tickethub/backend/repositories"
	gormrepo "github.com/shimauma0312/module-tickethub/backend/repositories/gorm" // Alias to avoid conflict
	"gorm.io/gorm"
)

// RepositoryFactory はリポジトリを作成するためのファクトリ
type RepositoryFactory struct {
	gormDB *gorm.DB
}

// NewRepositoryFactory は新しいRepositoryFactoryを作成します
// このファクトリはGORM DBインスタンスのみを受け入れるように変更されました。
func NewRepositoryFactory(db *gorm.DB) *RepositoryFactory {
	return &RepositoryFactory{
		gormDB: db,
	}
}

// NewIssueRepository はIssueRepositoryを作成します
func (f *RepositoryFactory) NewIssueRepository() (repositories.IssueRepository, error) {
	return gormrepo.NewIssueRepository(f.gormDB), nil
}

// NewUserRepository はUserRepositoryを作成します
func (f *RepositoryFactory) NewUserRepository() (repositories.UserRepository, error) {
	repo, err := gormrepo.NewUserRepository(f.gormDB)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// NewAuthTokenRepository はAuthTokenRepositoryを作成します
func (f *RepositoryFactory) NewAuthTokenRepository() (repositories.AuthTokenRepository, error) {
	repo, err := gormrepo.NewAuthTokenRepository(f.gormDB)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// NewPasswordResetRepository はPasswordResetRepositoryを作成します
func (f *RepositoryFactory) NewPasswordResetRepository() (repositories.PasswordResetRepository, error) {
	repo, err := gormrepo.NewPasswordResetRepository(f.gormDB)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// NewNotificationRepository はNotificationRepositoryを作成します
func (f *RepositoryFactory) NewNotificationRepository() (repositories.NotificationRepository, error) {
	return gormrepo.NewNotificationRepository(f.gormDB), nil
}

// NewPushSubscriptionRepository はPushSubscriptionRepositoryを作成します
func (f *RepositoryFactory) NewPushSubscriptionRepository() (repositories.PushSubscriptionRepository, error) {
	return gormrepo.NewPushSubscriptionRepository(f.gormDB), nil
}

// NewNotificationTemplateRepository はNotificationTemplateRepositoryを作成します
func (f *RepositoryFactory) NewNotificationTemplateRepository() (repositories.NotificationTemplateRepository, error) {
	return gormrepo.NewNotificationTemplateRepository(f.gormDB), nil
}

// NewUserSettingsRepository はUserSettingsRepositoryを作成します
func (f *RepositoryFactory) NewUserSettingsRepository() (repositories.UserSettingsRepository, error) {
	repo, err := gormrepo.NewUserSettingsRepository(f.gormDB)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// NewLabelRepository はLabelRepositoryを作成します
func (f *RepositoryFactory) NewLabelRepository() (repositories.LabelRepository, error) {
	repo, err := gormrepo.NewLabelRepository(f.gormDB)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// NewMilestoneRepository はMilestoneRepositoryを作成します
func (f *RepositoryFactory) NewMilestoneRepository() (repositories.MilestoneRepository, error) {
	repo, err := gormrepo.NewMilestoneRepository(f.gormDB)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// NewDiscussionRepository はDiscussionRepositoryを作成します
func (f *RepositoryFactory) NewDiscussionRepository() (repositories.DiscussionRepository, error) {
	return gormrepo.NewDiscussionRepository(f.gormDB), nil
}

// NewCommentRepository はCommentRepositoryを作成します
func (f *RepositoryFactory) NewCommentRepository() (repositories.CommentRepository, error) {
	return gormrepo.NewCommentRepository(f.gormDB), nil
}

// NewReactionRepository はReactionRepositoryを作成します
func (f *RepositoryFactory) NewReactionRepository() (repositories.ReactionRepository, error) {
	return gormrepo.NewReactionRepository(f.gormDB), nil
}

// NewSearchService は検索サービスを作成します
func (f *RepositoryFactory) NewSearchService() (SearchService, error) {
	issueRepo, err := f.NewIssueRepository()
	if err != nil {
		return nil, err
	}

	commentRepo, err := f.NewCommentRepository()
	if err != nil {
		return nil, err
	}

	return NewSearchService(f.gormDB, issueRepo, commentRepo)
}

// NewSystemSettingsRepository はSystemSettingsRepositoryを作成します
func (f *RepositoryFactory) NewSystemSettingsRepository() (repositories.SystemSettingsRepository, error) {
	return gormrepo.NewSystemSettingsRepository(f.gormDB), nil
}

// NewActivityLogRepository はActivityLogRepositoryを作成します
func (f *RepositoryFactory) NewActivityLogRepository() (repositories.ActivityLogRepository, error) {
	return gormrepo.NewActivityLogRepository(f.gormDB), nil
}

// NewBackupRepository はBackupRepositoryを作成します
func (f *RepositoryFactory) NewBackupRepository() (repositories.BackupRepository, error) {
	return gormrepo.NewBackupRepository(f.gormDB), nil
}

// NewRepositoryRepository はRepositoryRepositoryを作成します
func (f *RepositoryFactory) NewRepositoryRepository() (repositories.RepositoryRepository, error) {
	return gormrepo.NewRepositoryRepository(f.gormDB), nil
}

// Close はデータベース接続をクローズします (GORMでは通常不要ですが、インターフェース互換性のために残すことも検討)
// GORMでは *gorm.DB のクローズは sql.DB 経由で行うため、このファクトリレベルでの明示的なCloseは不要かもしれません。
// もしアプリケーション終了時にDB接続を確実に閉じる必要がある場合は、main関数などで *gorm.DB から sql.DB を取得して Close() を呼び出してください。
func (f *RepositoryFactory) Close() error {
	// sqlDB, err := f.gormDB.DB()
	// if err != nil {
	// 	return err
	// }
	// return sqlDB.Close()
	return nil // GORMが接続管理を行うため、通常は不要
}
