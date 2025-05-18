package repositories

import (
	"context"

	"github.com/shimauma0312/module-tickethub/models"
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

// UserRepository はUser関連のデータベース操作を抽象化するインターフェース
type UserRepository interface {
	// Create は新しいUserを作成します
	Create(ctx context.Context, user *models.User) error
	// GetByID はIDによってUserを取得します
	GetByID(ctx context.Context, id int64) (*models.User, error)
	// GetByUsername はUsernameによってUserを取得します
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	// GetByEmail はEmailによってUserを取得します
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	// List は条件に一致するUserの一覧を取得します
	List(ctx context.Context, page, limit int) ([]*models.User, int, error)
	// Update は既存のUserを更新します
	Update(ctx context.Context, user *models.User) error
	// Delete はUserを削除します
	Delete(ctx context.Context, id int64) error
	// ChangePassword はパスワードを変更します
	ChangePassword(ctx context.Context, id int64, hashedPassword string) error
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
	Update(ctx context.Context, discussion *models.Discussion) error
	// Delete はDiscussionを削除します
	Delete(ctx context.Context, id int64) error
	// Search はDiscussionの全文検索を行います
	Search(ctx context.Context, query string, page, limit int) ([]*models.Discussion, int, error)
}

// CommentRepository はComment関連のデータベース操作を抽象化するインターフェース
type CommentRepository interface {
	// Create は新しいCommentを作成します
	Create(ctx context.Context, comment *models.Comment) error
	// GetByID はIDによってCommentを取得します
	GetByID(ctx context.Context, id int64) (*models.Comment, error)
	// ListByTarget は対象エンティティに関連するCommentの一覧を取得します
	ListByTarget(ctx context.Context, targetType models.CommentType, targetID int64, page, limit int) ([]*models.Comment, int, error)
	// ListReplies はコメントへの返信一覧を取得します
	ListReplies(ctx context.Context, parentCommentID int64, page, limit int) ([]*models.Comment, int, error)
	// Update は既存のCommentを更新します
	Update(ctx context.Context, comment *models.Comment) error
	// Delete はCommentを削除します
	Delete(ctx context.Context, id int64) error
	// AddReaction はコメントにリアクションを追加します
	AddReaction(ctx context.Context, reaction *models.Reaction) error
	// RemoveReaction はコメントからリアクションを削除します
	RemoveReaction(ctx context.Context, commentID, userID int64, emoji string) error
}

// LabelRepository はLabel関連のデータベース操作を抽象化するインターフェース
type LabelRepository interface {
	// Create は新しいLabelを作成します
	Create(ctx context.Context, label *models.Label) error
	// GetByID はIDによってLabelを取得します
	GetByID(ctx context.Context, id int64) (*models.Label, error)
	// GetByName は名前によってLabelを取得します
	GetByName(ctx context.Context, name string) (*models.Label, error)
	// List はラベルの一覧を取得します
	List(ctx context.Context, labelType models.LabelType, page, limit int) ([]*models.Label, int, error)
	// Update は既存のLabelを更新します
	Update(ctx context.Context, label *models.Label) error
	// Delete はLabelを削除します
	Delete(ctx context.Context, id int64) error
}

// MilestoneRepository はMilestone関連のデータベース操作を抽象化するインターフェース
type MilestoneRepository interface {
	// Create は新しいMilestoneを作成します
	Create(ctx context.Context, milestone *models.Milestone) error
	// GetByID はIDによってMilestoneを取得します
	GetByID(ctx context.Context, id int64) (*models.Milestone, error)
	// List はマイルストーンの一覧を取得します
	List(ctx context.Context, status models.MilestoneStatus, page, limit int) ([]*models.Milestone, int, error)
	// Update は既存のMilestoneを更新します
	Update(ctx context.Context, milestone *models.Milestone) error
	// Delete はMilestoneを削除します
	Delete(ctx context.Context, id int64) error
}

// RepositoryFactory はDBタイプに応じたリポジトリのインスタンスを生成するインターフェース
type RepositoryFactory interface {
	// NewIssueRepository はIssueRepositoryの新しいインスタンスを生成します
	NewIssueRepository() (IssueRepository, error)
	// NewUserRepository はUserRepositoryの新しいインスタンスを生成します
	NewUserRepository() (UserRepository, error)
	// NewDiscussionRepository はDiscussionRepositoryの新しいインスタンスを生成します
	NewDiscussionRepository() (DiscussionRepository, error)
	// NewCommentRepository はCommentRepositoryの新しいインスタンスを生成します
	NewCommentRepository() (CommentRepository, error)
	// NewLabelRepository はLabelRepositoryの新しいインスタンスを生成します
	NewLabelRepository() (LabelRepository, error)
	// NewMilestoneRepository はMilestoneRepositoryの新しいインスタンスを生成します
	NewMilestoneRepository() (MilestoneRepository, error)
	// InitDB はデータベースの初期化を行います
	InitDB() error
	// Migrate はマイグレーションを実行します
	Migrate(direction string, steps int) error
	// Close はデータベース接続をクローズします
	Close() error
}
