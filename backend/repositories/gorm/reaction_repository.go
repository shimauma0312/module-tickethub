package gorm

import (
	"context"

	"github.com/shimauma0312/module-tickethub/backend/models"
	"gorm.io/gorm"
)

type reactionRepository struct {
	db *gorm.DB
}

func NewReactionRepository(db *gorm.DB) *reactionRepository {
	return &reactionRepository{db: db}
}

func (r *reactionRepository) Create(ctx context.Context, reaction *models.Reaction) error {
	return r.db.WithContext(ctx).Create(reaction).Error
}

func (r *reactionRepository) GetByID(ctx context.Context, id int64) (*models.Reaction, error) {
	var reaction models.Reaction
	err := r.db.WithContext(ctx).First(&reaction, id).Error
	return &reaction, err
}

func (r *reactionRepository) GetByUserAndComment(ctx context.Context, userID, commentID int64, emoji string) (*models.Reaction, error) {
	var reaction models.Reaction
	err := r.db.WithContext(ctx).Where("user_id = ? AND comment_id = ? AND emoji = ?", userID, commentID, emoji).First(&reaction).Error
	return &reaction, err
}

func (r *reactionRepository) ListByComment(ctx context.Context, commentID int64) ([]*models.Reaction, error) {
	var reactions []*models.Reaction
	err := r.db.WithContext(ctx).Where("comment_id = ?", commentID).Find(&reactions).Error
	return reactions, err
}

func (r *reactionRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&models.Reaction{}, id).Error
}
