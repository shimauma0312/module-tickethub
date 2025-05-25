package models

import (
	"time"

	"gorm.io/gorm"
)

// SystemSettings はシステム全体の設定を表す構造体
type SystemSettings struct {
	ID                  int64     `json:"id"`
	SiteName            string    `json:"site_name"`
	SiteDescription     string    `json:"site_description"`
	SiteURL             string    `json:"site_url"`
	AllowSignup         bool      `json:"allow_signup"`
	DefaultLanguage     string    `json:"default_language"`
	DefaultTheme        string    `json:"default_theme"`
	EmailEnabled        bool      `json:"email_enabled"`
	EmailFromAddress    string    `json:"email_from_address"`
	EmailFromName       string    `json:"email_from_name"`
	SMTPHost            string    `json:"smtp_host,omitempty"`
	SMTPPort            int       `json:"smtp_port,omitempty"`
	SMTPUsername        string    `json:"smtp_username,omitempty"`
	SMTPPassword        string    `json:"smtp_password,omitempty"`
	SMTPUseTLS          bool      `json:"smtp_use_tls"`
	MaxFileUploadSize   int64     `json:"max_file_upload_size"`
	RequireEmailVerify  bool      `json:"require_email_verify"`
	AllowGuestAccess    bool      `json:"allow_guest_access"`
	MaintenanceMode     bool      `json:"maintenance_mode"`
	MaintenanceMessage  string    `json:"maintenance_message,omitempty"`
	BackupRetentionDays int       `json:"backup_retention_days"`
	LogRetentionDays    int       `json:"log_retention_days"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// NewDefaultSystemSettings はデフォルトのシステム設定を作成する
func NewDefaultSystemSettings() *SystemSettings {
	now := time.Now()
	return &SystemSettings{
		SiteName:            "TicketHub",
		SiteDescription:     "GitHub-like ticket management system",
		SiteURL:             "http://localhost:3000",
		AllowSignup:         true,
		DefaultLanguage:     "en",
		DefaultTheme:        "light",
		EmailEnabled:        false,
		EmailFromAddress:    "noreply@tickethub.local",
		EmailFromName:       "TicketHub",
		SMTPHost:            "",
		SMTPPort:            587,
		SMTPUsername:        "",
		SMTPPassword:        "",
		SMTPUseTLS:          true,
		MaxFileUploadSize:   10 * 1024 * 1024, // 10MB
		RequireEmailVerify:  false,
		AllowGuestAccess:    false,
		MaintenanceMode:     false,
		MaintenanceMessage:  "",
		BackupRetentionDays: 30,
		LogRetentionDays:    90,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}

// Update はシステム設定を更新する
func (s *SystemSettings) Update(settings map[string]interface{}) {
	s.UpdatedAt = time.Now()

	if val, exists := settings["site_name"]; exists {
		if str, ok := val.(string); ok {
			s.SiteName = str
		}
	}
	if val, exists := settings["site_description"]; exists {
		if str, ok := val.(string); ok {
			s.SiteDescription = str
		}
	}
	if val, exists := settings["site_url"]; exists {
		if str, ok := val.(string); ok {
			s.SiteURL = str
		}
	}
	if val, exists := settings["allow_signup"]; exists {
		if b, ok := val.(bool); ok {
			s.AllowSignup = b
		}
	}
	if val, exists := settings["default_language"]; exists {
		if str, ok := val.(string); ok {
			s.DefaultLanguage = str
		}
	}
	if val, exists := settings["default_theme"]; exists {
		if str, ok := val.(string); ok {
			s.DefaultTheme = str
		}
	}
	if val, exists := settings["email_enabled"]; exists {
		if b, ok := val.(bool); ok {
			s.EmailEnabled = b
		}
	}
	if val, exists := settings["email_from_address"]; exists {
		if str, ok := val.(string); ok {
			s.EmailFromAddress = str
		}
	}
	if val, exists := settings["email_from_name"]; exists {
		if str, ok := val.(string); ok {
			s.EmailFromName = str
		}
	}
	if val, exists := settings["smtp_host"]; exists {
		if str, ok := val.(string); ok {
			s.SMTPHost = str
		}
	}
	if val, exists := settings["smtp_port"]; exists {
		if i, ok := val.(int); ok {
			s.SMTPPort = i
		}
	}
	if val, exists := settings["smtp_username"]; exists {
		if str, ok := val.(string); ok {
			s.SMTPUsername = str
		}
	}
	if val, exists := settings["smtp_password"]; exists {
		if str, ok := val.(string); ok {
			s.SMTPPassword = str
		}
	}
	if val, exists := settings["smtp_use_tls"]; exists {
		if b, ok := val.(bool); ok {
			s.SMTPUseTLS = b
		}
	}
	if val, exists := settings["max_file_upload_size"]; exists {
		if i, ok := val.(int64); ok {
			s.MaxFileUploadSize = i
		}
	}
	if val, exists := settings["require_email_verify"]; exists {
		if b, ok := val.(bool); ok {
			s.RequireEmailVerify = b
		}
	}
	if val, exists := settings["allow_guest_access"]; exists {
		if b, ok := val.(bool); ok {
			s.AllowGuestAccess = b
		}
	}
	if val, exists := settings["maintenance_mode"]; exists {
		if b, ok := val.(bool); ok {
			s.MaintenanceMode = b
		}
	}
	if val, exists := settings["maintenance_message"]; exists {
		if str, ok := val.(string); ok {
			s.MaintenanceMessage = str
		}
	}
	if val, exists := settings["backup_retention_days"]; exists {
		if i, ok := val.(int); ok {
			s.BackupRetentionDays = i
		}
	}
	if val, exists := settings["log_retention_days"]; exists {
		if i, ok := val.(int); ok {
			s.LogRetentionDays = i
		}
	}
}

// AutoMigrateSystemSettings はSystemSettingsテーブルのマイグレーションを実行します
func AutoMigrateSystemSettings(db *gorm.DB) error {
	return db.AutoMigrate(&SystemSettings{})
}
