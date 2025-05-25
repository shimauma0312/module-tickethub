package gorm

import (
	"context"
	"fmt"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"gorm.io/gorm"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *commentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *commentRepository) GetByID(ctx context.Context, id int64) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.WithContext(ctx).First(&comment, id).Error
	return &comment, err
}

func (r *commentRepository) ListByTarget(ctx context.Context, targetID int64, targetType string, page, limit int) ([]*models.Comment, int, error) {
	var comments []*models.Comment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Comment{}).Where("target_id = ? AND target_type = ?", targetID, targetType)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&comments).Error
	return comments, int(total), err
}

func (r *commentRepository) ListReplies(ctx context.Context, parentCommentID int64, page, limit int) ([]*models.Comment, int, error) {
	var comments []*models.Comment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Comment{}).Where("parent_comment_id = ?", parentCommentID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&comments).Error
	return comments, int(total), err
}

func (r *commentRepository) Update(ctx context.Context, comment *models.Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *commentRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Comment{}, id).Error
}

func (r *commentRepository) Search(ctx context.Context, queryString string, page, limit int) ([]*models.Comment, int, error) {
	var comments []*models.Comment
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Comment{}).
		Where("body LIKE ?", "%"+queryString+"%")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&comments).Error
	return comments, int(total), err
}

func (r *commentRepository) GetAllOfType(ctx context.Context, commentType string) ([]*models.Comment, error) {
	var comments []*models.Comment

	err := r.db.WithContext(ctx).
		Where("type = ?", commentType).
		Find(&comments).Error

	if err != nil {
		return nil, err
	}

	return comments, nil
}

// CountComments は総コメント数を取得します
func (r *commentRepository) CountComments(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Comment{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count comments: %w", err)
	}
	return count, nil
}
