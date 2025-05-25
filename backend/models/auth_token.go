package models

import (
	"time"
)

// TokenType はトークンの種類を表す
type TokenType string

const (
	// AccessToken はアクセストークンを表す
	AccessToken TokenType = "access"
	// RefreshToken はリフレッシュトークンを表す
	RefreshToken TokenType = "refresh"
)

// AuthToken は認証トークン情報を表す構造体
type AuthToken struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	TokenType TokenType `json:"token_type"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	RevokedAt time.Time `json:"revoked_at,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
	IPAddress string    `json:"ip_address,omitempty"`
}

// NewAuthToken は新しいAuthTokenインスタンスを作成する
func NewAuthToken(userID int64, tokenType TokenType, token string, expiresIn time.Duration, userAgent, ipAddress string) *AuthToken {
	now := time.Now()
	return &AuthToken{
		UserID:    userID,
		TokenType: tokenType,
		Token:     token,
		ExpiresAt: now.Add(expiresIn),
		CreatedAt: now,
		UserAgent: userAgent,
		IPAddress: ipAddress,
	}
}

// IsExpired はトークンが期限切れかどうかを判定する
func (at *AuthToken) IsExpired() bool {
	return time.Now().After(at.ExpiresAt)
}

// IsRevoked はトークンが無効化されているかどうかを判定する
func (at *AuthToken) IsRevoked() bool {
	return !at.RevokedAt.IsZero()
}

// Revoke はトークンを無効化する
func (at *AuthToken) Revoke() {
	at.RevokedAt = time.Now()
}

// IsValid はトークンが有効かどうかを判定する
func (at *AuthToken) IsValid() bool {
	return at.Token != "" && !at.IsExpired() && !at.IsRevoked()
}
