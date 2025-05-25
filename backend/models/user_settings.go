package models

import (
	"time"
)

// UserSettings はユーザー設定情報を表す構造体
type UserSettings struct {
	UserID            int64     `json:"user_id"`
	EmailNotification bool      `json:"email_notification"`
	PushNotification  bool      `json:"push_notification"`
	NotificationTypes string    `json:"notification_types"` // all, mention, assign, comment, update のカンマ区切り
	Language          string    `json:"language"`
	Theme             string    `json:"theme"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// NewUserSettings は新しいUserSettingsインスタンスを作成する
func NewUserSettings(userID int64) *UserSettings {
	now := time.Now()
	return &UserSettings{
		UserID:            userID,
		EmailNotification: true,
		PushNotification:  true,
		NotificationTypes: "all",
		Language:          "en",
		Theme:             "light",
		UpdatedAt:         now,
	}
}

// UpdateEmailNotification はメール通知設定を更新する
func (us *UserSettings) UpdateEmailNotification(enabled bool) {
	us.EmailNotification = enabled
	us.UpdatedAt = time.Now()
}

// UpdateLanguage は言語設定を更新する
func (us *UserSettings) UpdateLanguage(language string) {
	us.Language = language
	us.UpdatedAt = time.Now()
}

// UpdateTheme はテーマ設定を更新する
func (us *UserSettings) UpdateTheme(theme string) {
	us.Theme = theme
	us.UpdatedAt = time.Now()
}

// UpdatePushNotification はプッシュ通知設定を更新する
func (us *UserSettings) UpdatePushNotification(enabled bool) {
	us.PushNotification = enabled
	us.UpdatedAt = time.Now()
}

// UpdateNotificationTypes は通知タイプ設定を更新する
func (us *UserSettings) UpdateNotificationTypes(types string) {
	us.NotificationTypes = types
	us.UpdatedAt = time.Now()
}
