package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/shimauma0312/module-tickethub/backend/models"
)

// NotificationTemplateRepository はPostgreSQL固有のNotificationTemplateリポジトリ実装です
type NotificationTemplateRepository struct {
	db *sql.DB
}

// NewNotificationTemplateRepository は新しいPostgreSQL NotificationTemplateRepositoryを生成します
func NewNotificationTemplateRepository(db *sql.DB) *NotificationTemplateRepository {
	return &NotificationTemplateRepository{
		db: db,
	}
}

// Create は新しいNotificationTemplateを作成します
func (r *NotificationTemplateRepository) Create(ctx context.Context, template *models.NotificationTemplate) error {
	query := `
		INSERT INTO notification_templates (type, title_template, body_template, email_subject_template, email_body_template, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		template.Type,
		template.TitleTemplate,
		template.BodyTemplate,
		template.EmailSubjectTemplate,
		template.EmailBodyTemplate,
		template.CreatedAt,
		template.UpdatedAt,
	).Scan(&template.ID)

	return err
}

// GetByType はタイプによってNotificationTemplateを取得します
func (r *NotificationTemplateRepository) GetByType(ctx context.Context, templateType string) (*models.NotificationTemplate, error) {
	query := `
		SELECT id, type, title_template, body_template, email_subject_template, email_body_template, created_at, updated_at
		FROM notification_templates
		WHERE type = $1
	`
	row := r.db.QueryRowContext(ctx, query, templateType)

	template := &models.NotificationTemplate{}
	err := row.Scan(
		&template.ID,
		&template.Type,
		&template.TitleTemplate,
		&template.BodyTemplate,
		&template.EmailSubjectTemplate,
		&template.EmailBodyTemplate,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

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
		err := rows.Scan(
			&template.ID,
			&template.Type,
			&template.TitleTemplate,
			&template.BodyTemplate,
			&template.EmailSubjectTemplate,
			&template.EmailBodyTemplate,
			&template.CreatedAt,
			&template.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

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
		SET title_template = $1, body_template = $2, email_subject_template = $3, email_body_template = $4, updated_at = $5
		WHERE id = $6
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
	query := `DELETE FROM notification_templates WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
