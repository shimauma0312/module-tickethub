package models

import (
	"time"

	"gorm.io/gorm"
)

// IssueGorm はGORM用のIssue構造体
type IssueGorm struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Body        string         `gorm:"type:text;not null" json:"body"`
	Status      string         `gorm:"not null;default:open" json:"status"`
	AssigneeID  *int64         `gorm:"default:null" json:"assignee_id,omitempty"`
	CreatorID   int64          `gorm:"not null" json:"creator_id"`
	CreatedAt   time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null" json:"updated_at"`
	IsDraft     bool           `gorm:"not null;default:false" json:"is_draft"`
	MilestoneID *int64         `gorm:"default:null" json:"milestone_id,omitempty"`
	Labels      []IssueLabel   `gorm:"foreignKey:IssueID" json:"labels"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// IssueLabel はIssueのラベルを表すGORM構造体
type IssueLabel struct {
	IssueID int64  `gorm:"primaryKey;not null"`
	Label   string `gorm:"primaryKey;not null"`
}

// TableName はテーブル名を指定
func (IssueGorm) TableName() string {
	return "issues"
}

// TableName はテーブル名を指定
func (IssueLabel) TableName() string {
	return "issue_labels"
}

// ToModel はGORMモデルを通常のモデルに変換
func (i *IssueGorm) ToModel() *Issue {
	issue := &Issue{
		ID:        i.ID,
		Title:     i.Title,
		Body:      i.Body,
		Status:    i.Status,
		CreatorID: i.CreatorID,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
		IsDraft:   i.IsDraft,
		Labels:    make([]string, 0, len(i.Labels)),
	}

	if i.AssigneeID != nil {
		issue.AssigneeID = *i.AssigneeID
	}

	if i.MilestoneID != nil {
		issue.MilestoneID = *i.MilestoneID
	}

	for _, label := range i.Labels {
		issue.Labels = append(issue.Labels, label.Label)
	}

	return issue
}

// FromModel は通常のモデルからGORMモデルを生成
func IssueFromModel(issue *Issue) *IssueGorm {
	gormIssue := &IssueGorm{
		ID:        issue.ID,
		Title:     issue.Title,
		Body:      issue.Body,
		Status:    issue.Status,
		CreatorID: issue.CreatorID,
		CreatedAt: issue.CreatedAt,
		UpdatedAt: issue.UpdatedAt,
		IsDraft:   issue.IsDraft,
		Labels:    make([]IssueLabel, 0, len(issue.Labels)),
	}

	if issue.AssigneeID > 0 {
		gormIssue.AssigneeID = &issue.AssigneeID
	}

	if issue.MilestoneID > 0 {
		gormIssue.MilestoneID = &issue.MilestoneID
	}

	for _, label := range issue.Labels {
		gormIssue.Labels = append(gormIssue.Labels, IssueLabel{
			IssueID: issue.ID,
			Label:   label,
		})
	}

	return gormIssue
}

// AutoMigrate はGORMのマイグレーションを実行
func AutoMigrateIssue(db *gorm.DB) error {
	return db.AutoMigrate(&IssueGorm{}, &IssueLabel{})
}
