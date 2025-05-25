package models

import (
	"time"
)

// PushSubscription はブラウザのプッシュ通知サブスクリプション情報を表す構造体
type PushSubscription struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Endpoint  string    `json:"endpoint"`
	P256dh    string    `json:"p256dh"`
	Auth      string    `json:"auth"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewPushSubscription は新しいPushSubscriptionインスタンスを作成する
func NewPushSubscription(userID int64, endpoint, p256dh, auth string) *PushSubscription {
	now := time.Now()
	return &PushSubscription{
		UserID:    userID,
		Endpoint:  endpoint,
		P256dh:    p256dh,
		Auth:      auth,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsValid はPushSubscriptionの検証を行う
func (ps *PushSubscription) IsValid() bool {
	return ps.UserID > 0 && ps.Endpoint != "" && ps.P256dh != "" && ps.Auth != ""
}
