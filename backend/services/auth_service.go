package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
	"golang.org/x/crypto/bcrypt"
)

const (
	// アクセストークンの有効期限（15分）
	accessTokenExpiration = 15 * time.Minute
	// リフレッシュトークンの有効期限（7日）
	refreshTokenExpiration = 7 * 24 * time.Hour
	// パスワードリセットトークンの有効期限（24時間）
	passwordResetExpiration = 24 * time.Hour
	// bcryptのコスト（セキュリティレベル）
	bcryptCost = 12
	// CSRFトークンの有効期限（1時間）
	csrfTokenExpiration = 1 * time.Hour
)

var (
	// ErrInvalidCredentials は無効な認証情報のエラー
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrTokenExpired はトークンの期限切れのエラー
	ErrTokenExpired = errors.New("token has expired")
	// ErrTokenInvalid は無効なトークンのエラー
	ErrTokenInvalid = errors.New("token is invalid")
	// ErrUserNotFound はユーザーが見つからないエラー
	ErrUserNotFound = errors.New("user not found")
	// ErrUserDisabled は無効化されたユーザーのエラー
	ErrUserDisabled = errors.New("user account is disabled")
	// ErrInvalidResetToken は無効なリセットトークンのエラー
	ErrInvalidResetToken = errors.New("invalid or expired password reset token")
)

// JWTClaims はJWTトークンのクレーム
type JWTClaims struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	IsAdmin   bool   `json:"is_admin"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// AuthService は認証関連の機能を提供するサービス
type AuthService struct {
	userRepo          repositories.UserRepository
	tokenRepo         repositories.AuthTokenRepository
	passwordResetRepo repositories.PasswordResetRepository
	jwtSecret         []byte
}

// NewAuthService は新しいAuthServiceを作成します
func NewAuthService(
	userRepo repositories.UserRepository,
	tokenRepo repositories.AuthTokenRepository,
	passwordResetRepo repositories.PasswordResetRepository,
	jwtSecret string,
) *AuthService {
	return &AuthService{
		userRepo:          userRepo,
		tokenRepo:         tokenRepo,
		passwordResetRepo: passwordResetRepo,
		jwtSecret:         []byte(jwtSecret),
	}
}

// Register はユーザー登録を行います
func (s *AuthService) Register(ctx context.Context, username, email, password, fullName string) (*models.User, error) {
	// ユーザー名の重複チェック
	existingUser, err := s.userRepo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("username already exists")
	}

	// メールアドレスの重複チェック
	existingUser, err = s.userRepo.GetByEmail(ctx, email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	// パスワードのハッシュ化
	hashedPassword, err := s.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// ユーザーの作成
	user := models.NewUser(username, email, hashedPassword, fullName)
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login はユーザーログインを行います
func (s *AuthService) Login(ctx context.Context, usernameOrEmail, password, userAgent, ipAddress string) (*models.User, string, string, error) {
	// ユーザー名またはメールアドレスでユーザーを検索
	var user *models.User
	var err error

	// メールアドレスの場合
	user, err = s.userRepo.GetByEmail(ctx, usernameOrEmail)
	if err != nil || user == nil {
		// ユーザー名の場合
		user, err = s.userRepo.GetByUsername(ctx, usernameOrEmail)
		if err != nil || user == nil {
			return nil, "", "", ErrInvalidCredentials
		}
	}

	// アカウントが有効かチェック
	if !user.IsActive {
		return nil, "", "", ErrUserDisabled
	}

	// パスワードの検証
	if !s.CheckPasswordHash(password, user.Password) {
		return nil, "", "", ErrInvalidCredentials
	}

	// 最終ログイン日時を更新
	user.RecordLogin()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, "", "", fmt.Errorf("failed to update last login: %w", err)
	}

	// アクセストークンの生成
	accessToken, err := s.GenerateJWT(user.ID, user.Username, user.IsAdmin, string(models.AccessToken), accessTokenExpiration)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate access token: %w", err)
	}

	// リフレッシュトークンの生成
	refreshToken, err := s.GenerateRandomToken(32)
	if err != nil {
		return nil, "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// トークンをデータベースに保存
	refreshTokenModel := models.NewAuthToken(
		user.ID,
		models.RefreshToken,
		refreshToken,
		refreshTokenExpiration,
		userAgent,
		ipAddress,
	)
	if err := s.tokenRepo.Create(ctx, refreshTokenModel); err != nil {
		return nil, "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return user, accessToken, refreshToken, nil
}

// Logout はユーザーログアウトを行います
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	// リフレッシュトークンを検索
	token, err := s.tokenRepo.GetByToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("failed to find refresh token: %w", err)
	}

	// トークンを無効化
	token.Revoke()
	if err := s.tokenRepo.Update(ctx, token); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

// RefreshToken はトークンの更新を行います
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken, userAgent, ipAddress string) (string, string, error) {
	// リフレッシュトークンを検索
	token, err := s.tokenRepo.GetByToken(ctx, refreshToken)
	if err != nil || token == nil {
		return "", "", ErrTokenInvalid
	}

	// トークンの有効性をチェック
	if !token.IsValid() {
		return "", "", ErrTokenInvalid
	}

	// ユーザーを取得
	user, err := s.userRepo.GetByID(ctx, token.UserID)
	if err != nil || user == nil {
		return "", "", ErrUserNotFound
	}

	// アカウントが有効かチェック
	if !user.IsActive {
		return "", "", ErrUserDisabled
	}

	// 現在のリフレッシュトークンを無効化
	token.Revoke()
	if err := s.tokenRepo.Update(ctx, token); err != nil {
		return "", "", fmt.Errorf("failed to revoke old token: %w", err)
	}

	// 新しいアクセストークンを生成
	accessToken, err := s.GenerateJWT(user.ID, user.Username, user.IsAdmin, string(models.AccessToken), accessTokenExpiration)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new access token: %w", err)
	}

	// 新しいリフレッシュトークンを生成
	newRefreshToken, err := s.GenerateRandomToken(32)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new refresh token: %w", err)
	}

	// 新しいトークンをデータベースに保存
	refreshTokenModel := models.NewAuthToken(
		user.ID,
		models.RefreshToken,
		newRefreshToken,
		refreshTokenExpiration,
		userAgent,
		ipAddress,
	)
	if err := s.tokenRepo.Create(ctx, refreshTokenModel); err != nil {
		return "", "", fmt.Errorf("failed to save new refresh token: %w", err)
	}

	return accessToken, newRefreshToken, nil
}

// InitiatePasswordReset はパスワードリセットプロセスを開始します
func (s *AuthService) InitiatePasswordReset(ctx context.Context, email string) (*models.PasswordReset, error) {
	// メールアドレスでユーザーを検索
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	// アカウントが有効かチェック
	if !user.IsActive {
		return nil, ErrUserDisabled
	}

	// 既存のリセットトークンを無効化
	if err := s.passwordResetRepo.RevokeAllForUser(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to revoke existing tokens: %w", err)
	}

	// リセットトークンを生成
	resetToken, err := s.GenerateRandomToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate reset token: %w", err)
	}

	// パスワードリセットを作成
	hours := int(passwordResetExpiration.Hours())
	passwordReset := models.NewPasswordReset(user.ID, resetToken, hours)
	if err := s.passwordResetRepo.Create(ctx, passwordReset); err != nil {
		return nil, fmt.Errorf("failed to create password reset: %w", err)
	}

	return passwordReset, nil
}

// ResetPassword はパスワードをリセットします
func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// トークンを検索
	passwordReset, err := s.passwordResetRepo.GetByToken(ctx, token)
	if err != nil || passwordReset == nil {
		return ErrInvalidResetToken
	}

	// トークンの有効性をチェック
	if !passwordReset.IsValid() {
		return ErrInvalidResetToken
	}

	// ユーザーを取得
	user, err := s.userRepo.GetByID(ctx, passwordReset.UserID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	// アカウントが有効かチェック
	if !user.IsActive {
		return ErrUserDisabled
	}

	// パスワードをハッシュ化
	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// パスワードを更新
	user.SetPassword(hashedPassword)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// トークンを使用済みにマーク
	passwordReset.MarkAsUsed()
	if err := s.passwordResetRepo.Update(ctx, passwordReset); err != nil {
		return fmt.Errorf("failed to mark token as used: %w", err)
	}

	// ユーザーの全トークンを無効化（セキュリティのため）
	if err := s.tokenRepo.RevokeAllForUser(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to revoke user tokens: %w", err)
	}

	return nil
}

// ChangePassword はパスワードを変更します
func (s *AuthService) ChangePassword(ctx context.Context, userID int64, currentPassword, newPassword string) error {
	// ユーザーを取得
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	// アカウントが有効かチェック
	if !user.IsActive {
		return ErrUserDisabled
	}

	// 現在のパスワードを検証
	if !s.CheckPasswordHash(currentPassword, user.Password) {
		return ErrInvalidCredentials
	}

	// パスワードをハッシュ化
	hashedPassword, err := s.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// パスワードを更新
	user.SetPassword(hashedPassword)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	// ユーザーの全トークンを無効化（セキュリティのため）
	if err := s.tokenRepo.RevokeAllForUser(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to revoke user tokens: %w", err)
	}

	return nil
}

// VerifyJWT はJWTトークンを検証します
func (s *AuthService) VerifyJWT(tokenString string) (*JWTClaims, error) {
	// トークンを解析
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 署名アルゴリズムの検証
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// クレームを取得
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// トークンの有効期限をチェック
		if time.Now().After(claims.ExpiresAt.Time) {
			return nil, ErrTokenExpired
		}
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// GenerateJWT はJWTトークンを生成します
func (s *AuthService) GenerateJWT(userID int64, username string, isAdmin bool, tokenType string, expiration time.Duration) (string, error) {
	// クレームを作成
	now := time.Now()
	claims := JWTClaims{
		UserID:    userID,
		Username:  username,
		IsAdmin:   isAdmin,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "tickethub",
		},
	}

	// トークンを生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// GenerateCSRFToken はCSRFトークンを生成します
func (s *AuthService) GenerateCSRFToken(userID int64) (string, error) {
	// ランダムなトークンを生成
	token, err := s.GenerateRandomToken(32)
	if err != nil {
		return "", fmt.Errorf("failed to generate CSRF token: %w", err)
	}

	// CSRFトークンをデータベースに保存
	csrfToken := models.NewAuthToken(
		userID,
		"csrf",
		token,
		csrfTokenExpiration,
		"",
		"",
	)
	if err := s.tokenRepo.Create(context.Background(), csrfToken); err != nil {
		return "", fmt.Errorf("failed to save CSRF token: %w", err)
	}

	return token, nil
}

// VerifyCSRFToken はCSRFトークンを検証します
func (s *AuthService) VerifyCSRFToken(userID int64, token string) (bool, error) {
	// CSRFトークンを検索
	csrfToken, err := s.tokenRepo.GetByUserIDAndToken(context.Background(), userID, token)
	if err != nil || csrfToken == nil {
		return false, ErrTokenInvalid
	}

	// トークンの有効性をチェック
	if !csrfToken.IsValid() || csrfToken.TokenType != "csrf" {
		return false, ErrTokenInvalid
	}

	return true, nil
}

// HashPassword はパスワードをハッシュ化します
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

// CheckPasswordHash はパスワードとハッシュを比較します
func (s *AuthService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRandomToken はランダムなトークンを生成します
func (s *AuthService) GenerateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
