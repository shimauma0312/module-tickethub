package models

import (
	"time"
)

// Milestone はマイルストーン情報を表す構造体
type Milestone struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date,omitempty"`
	Status      string    `json:"status"` // open/closed
	CreatorID   int64     `json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// NewMilestone は新しいMilestoneインスタンスを作成する
func NewMilestone(title, description string, dueDate time.Time, creatorID int64) *Milestone {
	now := time.Now()
	return &Milestone{
		Title:       title,
		Description: description,
		DueDate:     dueDate,
		Status:      "open",
		CreatorID:   creatorID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// IsValid はマイルストーン情報の検証を行う
func (m *Milestone) IsValid() bool {
	return m.Title != "" && m.CreatorID > 0
}

// Close はマイルストーンをクローズ状態に設定する
func (m *Milestone) Close() {
	now := time.Now()
	m.Status = "closed"
	m.CompletedAt = now
	m.UpdatedAt = now
}

// Reopen はマイルストーンを再オープン状態に設定する
func (m *Milestone) Reopen() {
	var zeroTime time.Time
	m.Status = "open"
	m.CompletedAt = zeroTime
	m.UpdatedAt = time.Now()
}

// Update はマイルストーン情報を更新する
func (m *Milestone) Update(title, description string, dueDate time.Time) {
	m.Title = title
	m.Description = description
	m.DueDate = dueDate
	m.UpdatedAt = time.Now()
}
