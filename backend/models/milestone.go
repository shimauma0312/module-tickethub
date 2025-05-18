package models

import (
	"time"
)

// MilestoneStatus はマイルストーンのステータスを表す型
type MilestoneStatus string

const (
	// OpenMilestone は進行中のマイルストーンを表します
	OpenMilestone MilestoneStatus = "open"
	// ClosedMilestone は完了したマイルストーンを表します
	ClosedMilestone MilestoneStatus = "closed"
)

// Milestone はマイルストーンを表す構造体
type Milestone struct {
	ID          int64           `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description,omitempty"`
	DueDate     *time.Time      `json:"due_date,omitempty"`
	Status      MilestoneStatus `json:"status"`
	CreatorID   int64           `json:"creator_id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
}

// NewMilestone は新しいMilestoneインスタンスを作成する
func NewMilestone(title, description string, dueDate *time.Time, creatorID int64) *Milestone {
	now := time.Now()
	return &Milestone{
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      OpenMilestone,
		CreatorID:   creatorID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// IsValid はMilestoneの検証を行う
func (m *Milestone) IsValid() bool {
	return m.Title != "" && m.CreatorID > 0
}

// Close はマイルストーンを完了状態に設定する
func (m *Milestone) Close() {
	m.Status = ClosedMilestone
	now := time.Now()
	m.CompletedAt = &now
	m.UpdatedAt = now
}

// Reopen はマイルストーンを再オープン状態に設定する
func (m *Milestone) Reopen() {
	m.Status = OpenMilestone
	m.CompletedAt = nil
	m.UpdatedAt = time.Now()
}

// Update はマイルストーン情報を更新する
func (m *Milestone) Update(title, description string, dueDate *time.Time) {
	m.Title = title
	m.Description = description
	m.DueDate = dueDate
	m.UpdatedAt = time.Now()
}
