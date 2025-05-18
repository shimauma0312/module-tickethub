package models

import (
	"time"
)

// IssueStatus はチケットの状態を表す型
type IssueStatus string

const (
	// OpenIssue は未解決の課題を表します
	OpenIssue IssueStatus = "open"
	// ClosedIssue は解決済みの課題を表します
	ClosedIssue IssueStatus = "closed"
)

// Issue は課題チケットを表す構造体
type Issue struct {
	ID          int64       `json:"id"`
	Title       string      `json:"title"`
	Body        string      `json:"body"`
	Status      IssueStatus `json:"status"`
	Labels      []string    `json:"labels"`
	AssigneeID  int64       `json:"assignee_id"`
	CreatorID   int64       `json:"creator_id"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	IsDraft     bool        `json:"is_draft"`
	MilestoneID int64       `json:"milestone_id,omitempty"`
	ClosedAt    *time.Time  `json:"closed_at,omitempty"`
}

// NewIssue は新しいIssueインスタンスを作成する
func NewIssue(title, body string, creatorID int64) *Issue {
	now := time.Now()
	return &Issue{
		Title:     title,
		Body:      body,
		Status:    OpenIssue,
		CreatorID: creatorID,
		CreatedAt: now,
		UpdatedAt: now,
		IsDraft:   false,
		Labels:    []string{},
	}
}

// IsValid はIssueの検証を行う
func (i *Issue) IsValid() bool {
	return i.Title != "" && i.CreatorID > 0
}

// Close はIssueをクローズ状態に設定する
func (i *Issue) Close() {
	i.Status = ClosedIssue
	now := time.Now()
	i.ClosedAt = &now
	i.UpdatedAt = now
}

// Reopen はIssueを再オープン状態に設定する
func (i *Issue) Reopen() {
	i.Status = OpenIssue
	i.ClosedAt = nil
	i.UpdatedAt = time.Now()
}

// AddLabel はIssueにラベルを追加する
func (i *Issue) AddLabel(label string) {
	for _, l := range i.Labels {
		if l == label {
			return // 既に存在するラベルは追加しない
		}
	}
	i.Labels = append(i.Labels, label)
	i.UpdatedAt = time.Now()
}

// RemoveLabel はIssueからラベルを削除する
func (i *Issue) RemoveLabel(label string) {
	for j, l := range i.Labels {
		if l == label {
			i.Labels = append(i.Labels[:j], i.Labels[j+1:]...)
			i.UpdatedAt = time.Now()
			return
		}
	}
}
