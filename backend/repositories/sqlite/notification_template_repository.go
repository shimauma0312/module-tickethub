package sqlite

import (
	"context"
	"database/sql"
	"time"

	"github.com/shimauma0312/module-tickethub/backend/models"
)

// NotificationTemplateRepository はSQLite固有のNotificationTemplateリポジトリ実装です
type NotificationTemplateRepository struct {
	db *sql.DB
}

// NewNotificationTemplateRepository は新しいSQLite NotificationTemplateRepositoryを生成します
func NewNotificationTemplateRepository(db *sql.DB) *NotificationTemplateRepository {
	return &NotificationTemplateRepository{
		db: db,
	}
}

// Create は新しいNotificationTemplateを作成します
func (r *NotificationTemplateRepository) Create(ctx context.Context, template *models.NotificationTemplate) error {
	query := `
		INSERT INTO notification_templates (type, title_template, body_template, email_subject_template, email_body_template, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.ExecContext(
		ctx,
		query,
		template.Type,
		template.TitleTemplate,
		template.BodyTemplate,
		template.EmailSubjectTemplate,
		template.EmailBodyTemplate,
		template.CreatedAt,
		template.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	template.ID = id
	return nil
}

// GetByType はタイプによってNotificationTemplateを取得します
func (r *NotificationTemplateRepository) GetByType(ctx context.Context, templateType string) (*models.NotificationTemplate, error) {
	query := `
		SELECT id, type, title_template, body_template, email_subject_template, email_body_template, created_at, updated_at
		FROM notification_templates
		WHERE type = ?
	`
	row := r.db.QueryRowContext(ctx, query, templateType)

	template := &models.NotificationTemplate{}
	var createdAt, updatedAt string
	err := row.Scan(
		&template.ID,
		&template.Type,
		&template.TitleTemplate,
		&template.BodyTemplate,
		&template.EmailSubjectTemplate,
		&template.EmailBodyTemplate,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	template.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	template.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return template, nil
}

// GetAll は全てのNotificationTemplateを取得します
func (r *NotificationTemplateRepository) GetAll(ctx context.Context) ([]*models.NotificationTemplate, error) {
	query := `
		SELECT id, type, title_template, body_template, email_subject_template, email_body_template, created_at, updated_at
		FROM notification_templates
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var templates []*models.NotificationTemplate
	for rows.Next() {
		template := &models.NotificationTemplate{}
		var createdAt, updatedAt string
		err := rows.Scan(
			&template.ID,
			&template.Type,
			&template.TitleTemplate,
			&template.BodyTemplate,
			&template.EmailSubjectTemplate,
			&template.EmailBodyTemplate,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, err
		}

		template.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		template.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return templates, nil
}

// Update はNotificationTemplateを更新します
func (r *NotificationTemplateRepository) Update(ctx context.Context, template *models.NotificationTemplate) error {
	query := `
		UPDATE notification_templates
		SET title_template = ?, body_template = ?, email_subject_template = ?, email_body_template = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		template.TitleTemplate,
		template.BodyTemplate,
		template.EmailSubjectTemplate,
		template.EmailBodyTemplate,
		time.Now(),
		template.ID,
	)
	return err
}

// Delete はNotificationTemplateを削除します
func (r *NotificationTemplateRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM notification_templates WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
