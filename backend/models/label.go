package models

import (
	"time"
)

// LabelType はラベルが適用できるエンティティの種類を表す型
type LabelType string

const (
	// IssueLabelType はIssueに適用できるラベルを表します
	IssueLabelType LabelType = "issue"
	// DiscussionLabelType はDiscussionに適用できるラベルを表します
	DiscussionLabelType LabelType = "discussion"
	// BothLabelType は両方に適用できるラベルを表します
	BothLabelType LabelType = "both"
)

// Label はラベルを表す構造体
type Label struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Color       string    `json:"color"` // HEX形式の色コード（例: #FF0000）
	Type        LabelType `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewLabel は新しいLabelインスタンスを作成する
func NewLabel(name, description, color string, labelType LabelType) *Label {
	now := time.Now()
	return &Label{
		Name:        name,
		Description: description,
		Color:       color,
		Type:        labelType,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// IsValid はLabelの検証を行う
func (l *Label) IsValid() bool {
	return l.Name != "" && l.Color != ""
}

// Update はラベル情報を更新する
func (l *Label) Update(name, description, color string) {
	l.Name = name
	l.Description = description
	l.Color = color
	l.UpdatedAt = time.Now()
}
