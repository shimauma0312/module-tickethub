package gorm

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shimauma0312/module-tickethub/backend/models" // モジュールパスを修正
	"gorm.io/gorm"
)

// AuthTokenRepository はGORMベースのAuthTokenリポジトリ実装
type AuthTokenRepository struct {
	db *gorm.DB
}

// NewAuthTokenRepository は新しいAuthTokenRepositoryを作成します
func NewAuthTokenRepository(db *gorm.DB) (*AuthTokenRepository, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &AuthTokenRepository{db: db}, nil
}

// Create は新しいAuthTokenを作成します
func (r *AuthTokenRepository) Create(ctx context.Context, token *models.AuthToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return fmt.Errorf("failed to create auth token: %w", err)
	}
	return nil
}

// GetByID はIDによってAuthTokenを取得します
func (r *AuthTokenRepository) GetByID(ctx context.Context, id int64) (*models.AuthToken, error) {
	var token models.AuthToken
	if err := r.db.WithContext(ctx).First(&token, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("auth token with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get auth token by id: %w", err)
	}
	return &token, nil
}

// GetByToken はトークン文字列によってAuthTokenを取得します
func (r *AuthTokenRepository) GetByToken(ctx context.Context, tokenStr string) (*models.AuthToken, error) {
	var token models.AuthToken
	if err := r.db.WithContext(ctx).Where("token = ?", tokenStr).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("auth token %s not found", tokenStr)
		}
		return nil, fmt.Errorf("failed to get auth token by token string: %w", err)
	}
	return &token, nil
}

// GetByUserIDAndToken はユーザーIDとトークン文字列によってAuthTokenを取得します
func (r *AuthTokenRepository) GetByUserIDAndToken(ctx context.Context, userID int64, tokenStr string) (*models.AuthToken, error) {
	var token models.AuthToken
	if err := r.db.WithContext(ctx).Where("user_id = ? AND token = ?", userID, tokenStr).First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("auth token not found for user_id %d and token %s", userID, tokenStr)
		}
		return nil, fmt.Errorf("failed to get auth token by user_id and token string: %w", err)
	}
	return &token, nil
}

// GetValidTokensByUserID はユーザーIDによって有効なAuthTokenの一覧を取得します
func (r *AuthTokenRepository) GetValidTokensByUserID(ctx context.Context, userID int64, tokenType string) ([]*models.AuthToken, error) {
	var tokens []*models.AuthToken
	query := r.db.WithContext(ctx).Where("user_id = ? AND token_type = ? AND revoked_at IS NULL AND expires_at > ?", userID, tokenType, time.Now())
	if err := query.Find(&tokens).Error; err != nil {
		return nil, fmt.Errorf("failed to get valid auth tokens by user_id: %w", err)
	}
	return tokens, nil
}

// Update は既存のAuthTokenを更新します
func (r *AuthTokenRepository) Update(ctx context.Context, token *models.AuthToken) error {
	if err := r.db.WithContext(ctx).Save(token).Error; err != nil {
		return fmt.Errorf("failed to update auth token: %w", err)
	}
	return nil
}

// Delete はAuthTokenを削除します (GORMでは通常SaveでRevokedAtを設定するか、物理削除ならUnscoped)
// ここではRevokeと同様の振る舞いとしてRevokedAtを更新する例を示します。
// 物理削除が必要な場合は、Deleteメソッドの動作を再検討してください。
func (r *AuthTokenRepository) Delete(ctx context.Context, id int64) error {
	// 物理削除の場合:
	// if err := r.db.WithContext(ctx).Delete(&models.AuthToken{}, id).Error; err != nil {
	// 	 return fmt.Errorf("failed to delete auth token: %w", err)
	// }
	// 論理削除 (Revoke) の場合:
	token, err := r.GetByID(ctx, id)
	if err != nil {
		return err // GetByIDがエラーを返す
	}
	if token.RevokedAt.IsZero() { // まだ無効化されていない場合のみ更新
		now := time.Now()
		token.RevokedAt = now
		return r.Update(ctx, token)
	}
	return nil // すでに無効化されている場合は何もしない
}

// RevokeAllForUser はユーザーの全トークンを無効化します
func (r *AuthTokenRepository) RevokeAllForUser(ctx context.Context, userID int64) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&models.AuthToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error; err != nil {
		return fmt.Errorf("failed to revoke all auth tokens for user_id %d: %w", userID, err)
	}
	return nil
}
