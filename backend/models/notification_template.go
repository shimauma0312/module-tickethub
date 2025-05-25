package models

import (
	"time"
)

// NotificationTemplate は通知テンプレート情報を表す構造体
type NotificationTemplate struct {
	ID                   int64     `json:"id"`
	Type                 string    `json:"type"`
	TitleTemplate        string    `json:"title_template"`
	BodyTemplate         string    `json:"body_template"`
	EmailSubjectTemplate string    `json:"email_subject_template"`
	EmailBodyTemplate    string    `json:"email_body_template"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// NewNotificationTemplate は新しいNotificationTemplateインスタンスを作成する
func NewNotificationTemplate(
	notificationType,
	titleTemplate,
	bodyTemplate,
	emailSubjectTemplate,
	emailBodyTemplate string) *NotificationTemplate {

	now := time.Now()
	return &NotificationTemplate{
		Type:                 notificationType,
		TitleTemplate:        titleTemplate,
		BodyTemplate:         bodyTemplate,
		EmailSubjectTemplate: emailSubjectTemplate,
		EmailBodyTemplate:    emailBodyTemplate,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

// IsValid はNotificationTemplateの検証を行う
func (nt *NotificationTemplate) IsValid() bool {
	return nt.Type != "" &&
		nt.TitleTemplate != "" &&
		nt.BodyTemplate != "" &&
		nt.EmailSubjectTemplate != "" &&
		nt.EmailBodyTemplate != ""
}
