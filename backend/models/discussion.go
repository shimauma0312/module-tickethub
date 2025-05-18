package models

import (
	"time"
)

// DiscussionCategory はディスカッションのカテゴリを表す型
type DiscussionCategory string

const (
	// GeneralDiscussion は一般的な議論を表します
	GeneralDiscussion DiscussionCategory = "general"
	// QuestionDiscussion は質問を表します
	QuestionDiscussion DiscussionCategory = "question"
	// AnnouncementDiscussion はお知らせを表します
	AnnouncementDiscussion DiscussionCategory = "announcement"
)

// Discussion はディスカッションを表す構造体
type Discussion struct {
	ID          int64             `json:"id"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Category    DiscussionCategory `json:"category"`
	CreatorID   int64             `json:"creator_id"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	IsClosed    bool              `json:"is_closed"`
	IsPinned    bool              `json:"is_pinned"`
	Labels      []string          `json:"labels"`
}

// NewDiscussion は新しいDiscussionインスタンスを作成する
func NewDiscussion(title, body string, category DiscussionCategory, creatorID int64) *Discussion {
	now := time.Now()
	return &Discussion{
		Title:     title,
		Body:      body,
		Category:  category,
		CreatorID: creatorID,
		CreatedAt: now,
		UpdatedAt: now,
		IsClosed:  false,
		IsPinned:  false,
		Labels:    []string{},
	}
}

// IsValid はDiscussionの検証を行う
func (d *Discussion) IsValid() bool {
	return d.Title != "" && d.CreatorID > 0
}

// Close はディスカッションをクローズ状態に設定する
func (d *Discussion) Close() {
	d.IsClosed = true
	d.UpdatedAt = time.Now()
}

// Reopen はディスカッションを再オープン状態に設定する
func (d *Discussion) Reopen() {
	d.IsClosed = false
	d.UpdatedAt = time.Now()
}

// Pin はディスカッションをピン留めする
func (d *Discussion) Pin() {
	d.IsPinned = true
	d.UpdatedAt = time.Now()
}

// Unpin はディスカッションのピン留めを解除する
func (d *Discussion) Unpin() {
	d.IsPinned = false
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
