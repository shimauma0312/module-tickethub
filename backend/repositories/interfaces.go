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
	Update(ctx context.Context, issue *models.Issue) error // Delete はIssueを削除します
	Delete(ctx context.Context, id int64) error
	// Search はIssueの全文検索を行います
	Search(ctx context.Context, query string, page, limit int) ([]*models.Issue, int, error)
	// GetAll はすべてのIssueを取得します（検索インデックス構築用）
	GetAll(ctx context.Context) ([]*models.Issue, error)
	// CountIssues は総Issue数を取得します
	CountIssues(ctx context.Context) (int64, error)
	// CountOpenIssues はオープンなIssue数を取得します
	CountOpenIssues(ctx context.Context) (int64, error)
}

// UserRepository はUser関連のデータベース操作を抽象化するインターフェース
type UserRepository interface {
	// Create は新しいUserを作成します
	Create(ctx context.Context, user *models.User) error
	// GetByID はIDによってUserを取得します
	GetByID(ctx context.Context, id int64) (*models.User, error)
	// GetByUsername はユーザー名によってUserを取得します
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	// GetByEmail はメールアドレスによってUserを取得します
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	// List は条件に一致するUserの一覧を取得します
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.User, int, error)
	// Update は既存のUserを更新します
	Update(ctx context.Context, user *models.User) error // Delete はUserを削除します
	Delete(ctx context.Context, id int64) error
	// Search はUserの全文検索を行います
	Search(ctx context.Context, query string, page, limit int) ([]*models.User, int, error)
	// CountUsers は総ユーザー数を取得します
	CountUsers(ctx context.Context) (int64, error)
	// CountActiveUsers はアクティブユーザー数を取得します（指定日数以内にログインしたユーザー）
	CountActiveUsers(ctx context.Context, days int) (int64, error)
}

// UserSettingsRepository はUserSettings関連のデータベース操作を抽象化するインターフェース
type UserSettingsRepository interface {
	// GetByUserID はユーザーIDによってUserSettingsを取得します
	GetByUserID(ctx context.Context, userID int64) (*models.UserSettings, error)
	// CreateOrUpdate はUserSettingsを作成または更新します
	CreateOrUpdate(ctx context.Context, settings *models.UserSettings) error
}

// AuthTokenRepository は認証トークン関連のデータベース操作を抽象化するインターフェース
type AuthTokenRepository interface {
	// Create は新しいAuthTokenを作成します
	Create(ctx context.Context, token *models.AuthToken) error
	// GetByID はIDによってAuthTokenを取得します
	GetByID(ctx context.Context, id int64) (*models.AuthToken, error)
	// GetByToken はトークン文字列によってAuthTokenを取得します
	GetByToken(ctx context.Context, token string) (*models.AuthToken, error)
	// GetByUserIDAndToken はユーザーIDとトークン文字列によってAuthTokenを取得します
	GetByUserIDAndToken(ctx context.Context, userID int64, token string) (*models.AuthToken, error)
	// GetValidTokensByUserID はユーザーIDによって有効なAuthTokenの一覧を取得します
	GetValidTokensByUserID(ctx context.Context, userID int64, tokenType string) ([]*models.AuthToken, error)
	// Update は既存のAuthTokenを更新します
	Update(ctx context.Context, token *models.AuthToken) error
	// Delete はAuthTokenを削除します
	Delete(ctx context.Context, id int64) error
	// RevokeAllForUser はユーザーの全トークンを無効化します
	RevokeAllForUser(ctx context.Context, userID int64) error
}

// PasswordResetRepository はパスワードリセット関連のデータベース操作を抽象化するインターフェース
type PasswordResetRepository interface {
	// Create は新しいPasswordResetを作成します
	Create(ctx context.Context, reset *models.PasswordReset) error
	// GetByID はIDによってPasswordResetを取得します
	GetByID(ctx context.Context, id int64) (*models.PasswordReset, error)
	// GetByToken はトークン文字列によってPasswordResetを取得します
	GetByToken(ctx context.Context, token string) (*models.PasswordReset, error)
	// GetByUserID はユーザーIDによってPasswordResetの一覧を取得します
	GetByUserID(ctx context.Context, userID int64) ([]*models.PasswordReset, error)
	// Update は既存のPasswordResetを更新します
	Update(ctx context.Context, reset *models.PasswordReset) error
	// Delete はPasswordResetを削除します
	Delete(ctx context.Context, id int64) error
	// RevokeAllForUser はユーザーの全トークンを無効化します
	RevokeAllForUser(ctx context.Context, userID int64) error
}

// LabelRepository はLabel関連のデータベース操作を抽象化するインターフェース
type LabelRepository interface {
	// Create は新しいLabelを作成します
	Create(ctx context.Context, label *models.Label) error
	// GetByID はIDによってLabelを取得します
	GetByID(ctx context.Context, id int64) (*models.Label, error)
	// GetByName は名前とタイプによってLabelを取得します
	GetByName(ctx context.Context, name, labelType string) (*models.Label, error)
	// List は条件に一致するLabelの一覧を取得します
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Label, int, error)
	// Update は既存のLabelを更新します
	Update(ctx context.Context, label *models.Label) error
	// Delete はLabelを削除します
	Delete(ctx context.Context, id int64) error
	// Search はLabelの全文検索を行います
	Search(ctx context.Context, query string, page, limit int) ([]*models.Label, int, error)
}

// MilestoneRepository はMilestone関連のデータベース操作を抽象化するインターフェース
type MilestoneRepository interface {
	// Create は新しいMilestoneを作成します
	Create(ctx context.Context, milestone *models.Milestone) error
	// GetByID はIDによってMilestoneを取得します
	GetByID(ctx context.Context, id int64) (*models.Milestone, error)
	// List は条件に一致するMilestoneの一覧を取得します
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Milestone, int, error)
	// Update は既存のMilestoneを更新します
	Update(ctx context.Context, milestone *models.Milestone) error
	// Delete はMilestoneを削除します
	Delete(ctx context.Context, id int64) error
	// Search はMilestoneの全文検索を行います
	Search(ctx context.Context, query string, page, limit int) ([]*models.Milestone, int, error)
}

// DiscussionRepository はDiscussion関連のデータベース操作を抽象化するインターフェース
type DiscussionRepository interface {
	// Create は新しいDiscussionを作成します
	Create(ctx context.Context, discussion *models.Discussion) error
	// GetByID はIDによってDiscussionを取得します
	GetByID(ctx context.Context, id int64) (*models.Discussion, error)
	// List は条件に一致するDiscussionの一覧を取得します
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.Discussion, int, error)
	// Update は既存のDiscussionを更新します
	Update(ctx context.Context, discussion *models.Discussion) error // Delete はDiscussionを削除します
	Delete(ctx context.Context, id int64) error
	// Search はDiscussionの全文検索を行います
	Search(ctx context.Context, query string, page, limit int) ([]*models.Discussion, int, error)
	// CountDiscussions は総Discussion数を取得します
	CountDiscussions(ctx context.Context) (int64, error)
	// CountOpenDiscussions はオープンなDiscussion数を取得します
	CountOpenDiscussions(ctx context.Context) (int64, error)
}

// CommentRepository はComment関連のデータベース操作を抽象化するインターフェース
type CommentRepository interface {
	// Create は新しいCommentを作成します
	Create(ctx context.Context, comment *models.Comment) error
	// GetByID はIDによってCommentを取得します
	GetByID(ctx context.Context, id int64) (*models.Comment, error)
	// ListByTarget はターゲットIDとタイプによってCommentの一覧を取得します
	ListByTarget(ctx context.Context, targetID int64, targetType string, page, limit int) ([]*models.Comment, int, error)
	// ListReplies は親コメントIDによって返信コメントの一覧を取得します
	ListReplies(ctx context.Context, parentCommentID int64, page, limit int) ([]*models.Comment, int, error)
	// Update は既存のCommentを更新します
	Update(ctx context.Context, comment *models.Comment) error
	// Delete はCommentを削除します
	Delete(ctx context.Context, id int64) error // Search はCommentの全文検索を行います
	Search(ctx context.Context, query string, page, limit int) ([]*models.Comment, int, error)
	// GetAllOfType は指定したタイプのすべてのコメントを取得します（検索インデックス構築用）
	GetAllOfType(ctx context.Context, commentType string) ([]*models.Comment, error)
	// CountComments は総コメント数を取得します
	CountComments(ctx context.Context) (int64, error)
}

// ReactionRepository はReaction関連のデータベース操作を抽象化するインターフェース
type ReactionRepository interface {
	// Create は新しいReactionを作成します
	Create(ctx context.Context, reaction *models.Reaction) error
	// GetByID はIDによってReactionを取得します
	GetByID(ctx context.Context, id int64) (*models.Reaction, error)
	// GetByUserAndComment はユーザーIDとコメントIDによってReactionを取得します
	GetByUserAndComment(ctx context.Context, userID, commentID int64, emoji string) (*models.Reaction, error)
	// ListByComment はコメントIDによってReactionの一覧を取得します
	ListByComment(ctx context.Context, commentID int64) ([]*models.Reaction, error)
	// Delete はReactionを削除します
	Delete(ctx context.Context, id int64) error
}

// NotificationRepository はNotification関連のデータベース操作を抽象化するインターフェース
type NotificationRepository interface {
	// Create は新しいNotificationを作成します
	Create(ctx context.Context, notification *models.Notification) error
	// GetByID はIDによってNotificationを取得します
	GetByID(ctx context.Context, id int64) (*models.Notification, error)
	// ListByUser はユーザーIDによってNotificationの一覧を取得します
	ListByUser(ctx context.Context, userID int64, isRead *bool, page, limit int) ([]*models.Notification, int, error)
	// MarkAsRead はNotificationを既読状態に更新します
	MarkAsRead(ctx context.Context, id int64) error
	// MarkAllAsRead はユーザーの全Notificationを既読状態に更新します
	MarkAllAsRead(ctx context.Context, userID int64) error
	// Delete はNotificationを削除します
	Delete(ctx context.Context, id int64) error
}

// PushSubscriptionRepository はPushSubscription関連のデータベース操作を抽象化するインターフェース
type PushSubscriptionRepository interface {
	// Create は新しいPushSubscriptionを作成します
	Create(ctx context.Context, subscription *models.PushSubscription) error
	// GetByUserID はユーザーIDによってPushSubscriptionを取得します
	GetByUserID(ctx context.Context, userID int64) ([]*models.PushSubscription, error)
	// DeleteByEndpoint はEndpointによってPushSubscriptionを削除します
	DeleteByEndpoint(ctx context.Context, endpoint string) error
	// DeleteByUserID はユーザーIDによってPushSubscriptionを削除します
	DeleteByUserID(ctx context.Context, userID int64) error
}

// NotificationTemplateRepository はNotificationTemplate関連のデータベース操作を抽象化するインターフェース
type NotificationTemplateRepository interface {
	// Create は新しいNotificationTemplateを作成します
	Create(ctx context.Context, template *models.NotificationTemplate) error
	// GetByType はタイプによってNotificationTemplateを取得します
	GetByType(ctx context.Context, templateType string) (*models.NotificationTemplate, error)
	// GetAll は全てのNotificationTemplateを取得します
	GetAll(ctx context.Context) ([]*models.NotificationTemplate, error)
	// Update はNotificationTemplateを更新します
	Update(ctx context.Context, template *models.NotificationTemplate) error
	// Delete はNotificationTemplateを削除します
	Delete(ctx context.Context, id int64) error
}

// SystemSettingsRepository はSystemSettings関連のデータベース操作を抽象化するインターフェース
type SystemSettingsRepository interface {
	// Get はシステム設定を取得します
	Get(ctx context.Context) (*models.SystemSettings, error)
	// CreateOrUpdate はシステム設定を作成または更新します
	CreateOrUpdate(ctx context.Context, settings *models.SystemSettings) error
}

// ActivityLogRepository はActivityLog関連のデータベース操作を抽象化するインターフェース
type ActivityLogRepository interface {
	// Create は新しいActivityLogを作成します
	Create(ctx context.Context, log *models.ActivityLog) error
	// List は条件に一致するActivityLogの一覧を取得します
	List(ctx context.Context, filter map[string]interface{}, page, limit int) ([]*models.ActivityLog, int, error)
	// GetMetrics はシステムメトリクスを取得します
	GetMetrics(ctx context.Context) (*models.SystemMetrics, error)
	// DeleteOldLogs は古いログを削除します
	DeleteOldLogs(ctx context.Context, retentionDays int) error
}

// BackupRepository はBackup関連のデータベース操作を抽象化するインターフェース
type BackupRepository interface {
	// Create は新しいBackupInfoを作成します
	Create(ctx context.Context, backup *models.BackupInfo) error
	// List はBackupInfoの一覧を取得します
	List(ctx context.Context, page, limit int) ([]*models.BackupInfo, int, error)
	// GetByID はIDによってBackupInfoを取得します
	GetByID(ctx context.Context, id int64) (*models.BackupInfo, error)
	// Update はBackupInfoを更新します
	Update(ctx context.Context, backup *models.BackupInfo) error // Delete はBackupInfoを削除します
	Delete(ctx context.Context, id int64) error
	// DeleteOldBackups は古いバックアップを削除します
	DeleteOldBackups(ctx context.Context, retentionDays int) error
	// GetLatestBackup は最新のバックアップを取得します
	GetLatestBackup(ctx context.Context) (*models.BackupInfo, error)
}

// RepositoryFactory はDBタイプに応じたリポジトリのインスタンスを生成するインターフェース
type RepositoryFactory interface {
	// NewIssueRepository はIssueRepositoryの新しいインスタンスを生成します
	NewIssueRepository() (IssueRepository, error)
	// NewUserRepository はUserRepositoryの新しいインスタンスを生成します
	NewUserRepository() (UserRepository, error)
	// NewUserSettingsRepository はUserSettingsRepositoryの新しいインスタンスを生成します
	NewUserSettingsRepository() (UserSettingsRepository, error)
	// NewLabelRepository はLabelRepositoryの新しいインスタンスを生成します
	NewLabelRepository() (LabelRepository, error)
	// NewMilestoneRepository はMilestoneRepositoryの新しいインスタンスを生成します
	NewMilestoneRepository() (MilestoneRepository, error)
	// NewDiscussionRepository はDiscussionRepositoryの新しいインスタンスを生成します
	NewDiscussionRepository() (DiscussionRepository, error)
	// NewCommentRepository はCommentRepositoryの新しいインスタンスを生成します
	NewCommentRepository() (CommentRepository, error)
	// NewReactionRepository はReactionRepositoryの新しいインスタンスを生成します
	NewReactionRepository() (ReactionRepository, error)
	// NewNotificationRepository はNotificationRepositoryの新しいインスタンスを生成します
	NewNotificationRepository() (NotificationRepository, error)
	// NewPushSubscriptionRepository はPushSubscriptionRepositoryの新しいインスタンスを生成します
	NewPushSubscriptionRepository() (PushSubscriptionRepository, error)
	// NewNotificationTemplateRepository はNotificationTemplateRepositoryの新しいインスタンスを生成します
	NewNotificationTemplateRepository() (NotificationTemplateRepository, error)
	// NewSystemSettingsRepository はSystemSettingsRepositoryの新しいインスタンスを生成します
	NewSystemSettingsRepository() (SystemSettingsRepository, error)
	// NewActivityLogRepository はActivityLogRepositoryの新しいインスタンスを生成します
	NewActivityLogRepository() (ActivityLogRepository, error) // NewBackupRepository はBackupRepositoryの新しいインスタンスを生成します
	NewBackupRepository() (BackupRepository, error)
	// NewRepositoryRepository はRepositoryRepositoryの新しいインスタンスを生成します
	NewRepositoryRepository() (RepositoryRepository, error)
	// Close はデータベース接続をクローズします
	Close() error
}
