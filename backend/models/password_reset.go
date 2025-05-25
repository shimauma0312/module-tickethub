package models

import (
	"time"
)

// PasswordReset はパスワードリセット情報を表す構造体
type PasswordReset struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UsedAt    time.Time `json:"used_at,omitempty"`
}

// NewPasswordReset は新しいPasswordResetインスタンスを作成する
func NewPasswordReset(userID int64, token string, expiresInHours int) *PasswordReset {
	now := time.Now()
	return &PasswordReset{
		UserID:    userID,
		Token:     token,
		ExpiresAt: now.Add(time.Duration(expiresInHours) * time.Hour),
		CreatedAt: now,
	}
}

// IsExpired はトークンが期限切れかどうかを判定する
func (pr *PasswordReset) IsExpired() bool {
	return time.Now().After(pr.ExpiresAt)
}

// MarkAsUsed はトークンを使用済みとしてマークする
func (pr *PasswordReset) MarkAsUsed() {
	pr.UsedAt = time.Now()
}

// IsValid はトークンが有効かどうかを判定する
func (pr *PasswordReset) IsValid() bool {
	return pr.Token != "" && !pr.IsExpired() && pr.UsedAt.IsZero()
}
