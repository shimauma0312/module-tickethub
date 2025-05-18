package models

import (
	"time"
)

// CommentType はコメントの対象タイプを表す型
type CommentType string

const (
	// IssueComment はIssueに対するコメントを表します
	IssueComment CommentType = "issue"
	// DiscussionComment はDiscussionに対するコメントを表します
	DiscussionComment CommentType = "discussion"
	// ReplyComment は他のコメントへの返信を表します
	ReplyComment CommentType = "reply"
)

// Comment はコメントを表す構造体
type Comment struct {
	ID              int64       `json:"id"`
	Body            string      `json:"body"`
	CreatorID       int64       `json:"creator_id"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	Type            CommentType `json:"type"`
	TargetID        int64       `json:"target_id"`          // Issue/Discussion ID
	ParentCommentID int64       `json:"parent_comment_id,omitempty"` // 返信先コメントのID (存在する場合)
	Reactions       []*Reaction `json:"reactions,omitempty"`
	IsEdited        bool        `json:"is_edited"`
}

// Reaction はコメントへのリアクションを表す構造体
type Reaction struct {
	ID         int64     `json:"id"`
	CommentID  int64     `json:"comment_id"`
	UserID     int64     `json:"user_id"`
	Emoji      string    `json:"emoji"`      // 絵文字コード（例: ":thumbsup:", ":heart:"）
	CreatedAt  time.Time `json:"created_at"`
}

// NewComment は新しいCommentインスタンスを作成する
func NewComment(body string, creatorID, targetID int64, commentType CommentType) *Comment {
	now := time.Now()
	return &Comment{
		Body:      body,
		CreatorID: creatorID,
		TargetID:  targetID,
		Type:      commentType,
		CreatedAt: now,
		UpdatedAt: now,
		IsEdited:  false,
	}
}

// NewReplyComment は返信コメントを作成する
func NewReplyComment(body string, creatorID, targetID, parentCommentID int64) *Comment {
	now := time.Now()
	return &Comment{
		Body:            body,
		CreatorID:       creatorID,
		TargetID:        targetID,
		Type:            ReplyComment,
		ParentCommentID: parentCommentID,
		CreatedAt:       now,
		UpdatedAt:       now,
		IsEdited:        false,
	}
}

// IsValid はCommentの検証を行う
func (c *Comment) IsValid() bool {
	return c.Body != "" && c.CreatorID > 0 && c.TargetID > 0
}

// Edit はコメントを編集する
func (c *Comment) Edit(body string) {
	c.Body = body
	c.UpdatedAt = time.Now()
	c.IsEdited = true
}

// NewReaction は新しいリアクションを作成する
func NewReaction(commentID, userID int64, emoji string) *Reaction {
	return &Reaction{
		CommentID: commentID,
		UserID:    userID,
		Emoji:     emoji,
		CreatedAt: time.Now(),
	}
}

// AddReaction はコメントにリアクションを追加する
func (c *Comment) AddReaction(reaction *Reaction) {
	for _, r := range c.Reactions {
		if r.UserID == reaction.UserID && r.Emoji == reaction.Emoji {
			return // すでに同じリアクションが存在する場合は追加しない
		}
	}
	c.Reactions = append(c.Reactions, reaction)
}

// RemoveReaction はコメントからリアクションを削除する
func (c *Comment) RemoveReaction(userID int64, emoji string) {
	for i, r := range c.Reactions {
		if r.UserID == userID && r.Emoji == emoji {
			c.Reactions = append(c.Reactions[:i], c.Reactions[i+1:]...)
			return
		}
	}
}
