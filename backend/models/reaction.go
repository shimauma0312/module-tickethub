package models

import (
	"time"
)

// Reaction はリアクション情報を表す構造体
type Reaction struct {
	ID        int64     `json:"id"`
	CommentID int64     `json:"comment_id"`
	UserID    int64     `json:"user_id"`
	Emoji     string    `json:"emoji"`
	CreatedAt time.Time `json:"created_at"`
}

// NewReaction は新しいReactionインスタンスを作成する
func NewReaction(commentID, userID int64, emoji string) *Reaction {
	return &Reaction{
		CommentID: commentID,
		UserID:    userID,
		Emoji:     emoji,
		CreatedAt: time.Now(),
	}
}

// IsValid はリアクション情報の検証を行う
func (r *Reaction) IsValid() bool {
	return r.CommentID > 0 && r.UserID > 0 && r.Emoji != ""
}
