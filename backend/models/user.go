package models

import (
	"time"
)

// User はユーザー情報を表す構造体
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // パスワードはJSONに含めない
	FullName  string    `json:"full_name"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LastLogin time.Time `json:"last_login,omitempty"`
	IsAdmin   bool      `json:"is_admin"`
	IsActive  bool      `json:"is_active"`
}

// NewUser は新しいUserインスタンスを作成する
func NewUser(username, email, password, fullName string) *User {
	now := time.Now()
	return &User{
		Username:  username,
		Email:     email,
		Password:  password, // 注: パスワードは暗号化済みであることを前提
		FullName:  fullName,
		CreatedAt: now,
		UpdatedAt: now,
		IsAdmin:   false,
		IsActive:  true,
	}
}

// IsValid はユーザー情報の検証を行う
func (u *User) IsValid() bool {
	return u.Username != "" && u.Email != "" && u.Password != ""
}

// SetPassword はパスワードを設定する（実際のハッシュ化は呼び出し側で行う）
func (u *User) SetPassword(hashedPassword string) {
	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
}

// UpdateProfile はプロフィール情報を更新する
func (u *User) UpdateProfile(fullName, avatarURL string) {
	u.FullName = fullName
	u.AvatarURL = avatarURL
	u.UpdatedAt = time.Now()
}

// Deactivate はユーザーを非アクティブ状態に設定する
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate はユーザーをアクティブ状態に設定する
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// SetAdmin は管理者権限を設定する
func (u *User) SetAdmin(isAdmin bool) {
	u.IsAdmin = isAdmin
	u.UpdatedAt = time.Now()
}

// RecordLogin はログイン日時を記録する
func (u *User) RecordLogin() {
	u.LastLogin = time.Now()
}
