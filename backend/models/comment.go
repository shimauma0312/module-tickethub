package models

import (
	"time"
)

// Comment はコメント情報を表す構造体
type Comment struct {
	ID              int64     `json:"id"`
	Body            string    `json:"body"`
	CreatorID       int64     `json:"creator_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Type            string    `json:"type"` // issue/discussion/reply
	TargetID        int64     `json:"target_id"`
	ParentCommentID int64     `json:"parent_comment_id,omitempty"`
	IsEdited        bool      `json:"is_edited"`
}

// NewComment は新しいCommentインスタンスを作成する
func NewComment(body string, creatorID, targetID int64, commentType string) *Comment {
	now := time.Now()
	return &Comment{
		Body:      body,
		CreatorID: creatorID,
		CreatedAt: now,
		UpdatedAt: now,
		Type:      commentType,
		TargetID:  targetID,
		IsEdited:  false,
	}
}

// NewReply は返信コメントを作成する
func NewReply(body string, creatorID, targetID, parentCommentID int64, commentType string) *Comment {
	comment := NewComment(body, creatorID, targetID, commentType)
	comment.ParentCommentID = parentCommentID
	comment.Type = "reply"
	return comment
}

// IsValid はコメント情報の検証を行う
func (c *Comment) IsValid() bool {
	return c.Body != "" && c.CreatorID > 0 && c.TargetID > 0 && isValidCommentType(c.Type)
}

// isValidCommentType はコメントタイプが有効かどうかを検証する
func isValidCommentType(commentType string) bool {
	return commentType == "issue" || commentType == "discussion" || commentType == "reply"
}

// Edit はコメント内容を編集する
func (c *Comment) Edit(body string) {
	c.Body = body
	c.IsEdited = true
	c.UpdatedAt = time.Now()
}

// IsReply は返信コメントかどうかを判定する
func (c *Comment) IsReply() bool {
	return c.Type == "reply" && c.ParentCommentID > 0
}
