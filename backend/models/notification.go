package models

import (
	"time"
)

// Notification は通知情報を表す構造体
type Notification struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Type       string    `json:"type"`        // mention/assign/comment/etc
	SourceType string    `json:"source_type"` // issue/discussion/comment
	SourceID   int64     `json:"source_id"`
	ActorID    int64     `json:"actor_id"`
	Message    string    `json:"message"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

// NewNotification は新しいNotificationインスタンスを作成する
func NewNotification(userID int64, notificationType, sourceType string, sourceID, actorID int64, message string) *Notification {
	return &Notification{
		UserID:     userID,
		Type:       notificationType,
		SourceType: sourceType,
		SourceID:   sourceID,
		ActorID:    actorID,
		Message:    message,
		IsRead:     false,
		CreatedAt:  time.Now(),
	}
}

// IsValid は通知情報の検証を行う
func (n *Notification) IsValid() bool {
	return n.UserID > 0 && n.ActorID > 0 && n.SourceID > 0 &&
		n.Type != "" && n.SourceType != "" && n.Message != ""
}

// MarkAsRead は通知を既読状態に設定する
func (n *Notification) MarkAsRead() {
	n.IsRead = true
}

// MarkAsUnread は通知を未読状態に設定する
func (n *Notification) MarkAsUnread() {
	n.IsRead = false
}
