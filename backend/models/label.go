package models

import (
	"time"
)

// Label はラベル情報を表す構造体
type Label struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Type        string    `json:"type"` // issue/discussion/both
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewLabel は新しいLabelインスタンスを作成する
func NewLabel(name, description, color, labelType string) *Label {
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

// IsValid はラベル情報の検証を行う
func (l *Label) IsValid() bool {
	return l.Name != "" && l.Color != "" && isValidLabelType(l.Type)
}

// isValidLabelType はラベルタイプが有効かどうかを検証する
func isValidLabelType(labelType string) bool {
	return labelType == "issue" || labelType == "discussion" || labelType == "both"
}

// Update はラベル情報を更新する
func (l *Label) Update(name, description, color, labelType string) {
	l.Name = name
	l.Description = description
	l.Color = color
	if isValidLabelType(labelType) {
		l.Type = labelType
	}
	l.UpdatedAt = time.Now()
}
