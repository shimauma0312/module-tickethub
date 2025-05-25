package models

import (
	"time"
)

// Discussion はディスカッション情報を表す構造体
type Discussion struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Status    string    `json:"status"`   // open/closed/answered
	Category  string    `json:"category"` // general/question/announcement
	Labels    []string  `json:"labels"`
	CreatorID int64     `json:"creator_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsDraft   bool      `json:"is_draft"`
	ClosedAt  time.Time `json:"closed_at,omitempty"`
}

// NewDiscussion は新しいDiscussionインスタンスを作成する
func NewDiscussion(title, body, category string, creatorID int64) *Discussion {
	now := time.Now()
	return &Discussion{
		Title:     title,
		Body:      body,
		Status:    "open",
		Category:  category,
		CreatorID: creatorID,
		CreatedAt: now,
		UpdatedAt: now,
		IsDraft:   false,
		Labels:    []string{},
	}
}

// IsValid はDiscussionの検証を行う
func (d *Discussion) IsValid() bool {
	return d.Title != "" && d.CreatorID > 0 && isValidCategory(d.Category)
}

// isValidCategory はカテゴリが有効かどうかを検証する
func isValidCategory(category string) bool {
	validCategories := []string{"general", "question", "announcement", "idea"}
	for _, c := range validCategories {
		if category == c {
			return true
		}
	}
	return false
}

// Close はディスカッションをクローズ状態に設定する
func (d *Discussion) Close() {
	now := time.Now()
	d.Status = "closed"
	d.ClosedAt = now
	d.UpdatedAt = now
}

// MarkAsAnswered はディスカッションを回答済み状態に設定する
func (d *Discussion) MarkAsAnswered() {
	d.Status = "answered"
	d.UpdatedAt = time.Now()
}

// Reopen はディスカッションを再オープン状態に設定する
func (d *Discussion) Reopen() {
	var zeroTime time.Time
	d.Status = "open"
	d.ClosedAt = zeroTime
	d.UpdatedAt = time.Now()
}

// AddLabel はディスカッションにラベルを追加する
func (d *Discussion) AddLabel(label string) {
	for _, l := range d.Labels {
		if l == label {
			return // 既に存在するラベルは追加しない
		}
	}
	d.Labels = append(d.Labels, label)
	d.UpdatedAt = time.Now()
}

// RemoveLabel はディスカッションからラベルを削除する
func (d *Discussion) RemoveLabel(label string) {
	for j, l := range d.Labels {
		if l == label {
			d.Labels = append(d.Labels[:j], d.Labels[j+1:]...)
			d.UpdatedAt = time.Now()
			return
		}
	}
}

// SetDraft はディスカッションの下書き状態を設定する
func (d *Discussion) SetDraft(isDraft bool) {
	d.IsDraft = isDraft
	d.UpdatedAt = time.Now()
}
