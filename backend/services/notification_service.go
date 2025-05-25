package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
)

// NotificationService は通知関連の機能を提供するサービス
type NotificationService struct {
	notificationRepo         repositories.NotificationRepository
	userRepo                 repositories.UserRepository
	userSettingsRepo         repositories.UserSettingsRepository
	pushSubscriptionRepo     repositories.PushSubscriptionRepository
	notificationTemplateRepo repositories.NotificationTemplateRepository
	// WebPush設定
	vapidPrivateKey string
	vapidPublicKey  string
	// Email設定
	smtpHost     string
	smtpPort     int
	smtpUsername string
	smtpPassword string
	smtpFrom     string
	baseURL      string
}

// NewNotificationService は新しいNotificationServiceのインスタンスを生成します
func NewNotificationService(
	factory repositories.RepositoryFactory,
	vapidPrivateKey, vapidPublicKey string,
	smtpHost string, smtpPort int, smtpUsername, smtpPassword, smtpFrom string,
	baseURL string) (*NotificationService, error) {

	notificationRepo, err := factory.NewNotificationRepository()
	if err != nil {
		return nil, err
	}

	userRepo, err := factory.NewUserRepository()
	if err != nil {
		return nil, err
	}

	userSettingsRepo, err := factory.NewUserSettingsRepository()
	if err != nil {
		return nil, err
	}

	pushSubscriptionRepo, err := factory.NewPushSubscriptionRepository()
	if err != nil {
		return nil, err
	}

	notificationTemplateRepo, err := factory.NewNotificationTemplateRepository()
	if err != nil {
		return nil, err
	}

	return &NotificationService{
		notificationRepo:         notificationRepo,
		userRepo:                 userRepo,
		userSettingsRepo:         userSettingsRepo,
		pushSubscriptionRepo:     pushSubscriptionRepo,
		notificationTemplateRepo: notificationTemplateRepo,
		vapidPrivateKey:          vapidPrivateKey,
		vapidPublicKey:           vapidPublicKey,
		smtpHost:                 smtpHost,
		smtpPort:                 smtpPort,
		smtpUsername:             smtpUsername,
		smtpPassword:             smtpPassword,
		smtpFrom:                 smtpFrom,
		baseURL:                  baseURL,
	}, nil
}

// NotificationData は通知テンプレート用のデータ構造
type NotificationData struct {
	User struct {
		Name  string
		Email string
	}
	Actor struct {
		Name  string
		Email string
	}
	Source struct {
		Type  string
		Title string
		URL   string
	}
	Message string
}

// CreateNotification は新しい通知を作成し、必要に応じてプッシュ通知やメール通知を送信します
func (s *NotificationService) CreateNotification(
	ctx context.Context,
	userID int64,
	notificationType string,
	sourceType string,
	sourceID int64,
	actorID int64,
	message string) error {

	// 通知オブジェクトを作成
	notification := models.NewNotification(userID, notificationType, sourceType, sourceID, actorID, message)

	if !notification.IsValid() {
		return fmt.Errorf("invalid notification data")
	}

	// 通知をDBに保存
	err := s.notificationRepo.Create(ctx, notification)
	if err != nil {
		return err
	}

	// ユーザー設定を取得
	userSettings, err := s.userSettingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// 通知の種類が通知設定に含まれているかチェック
	notificationTypes := strings.Split(userSettings.NotificationTypes, ",")
	shouldNotify := false
	if userSettings.NotificationTypes == "all" {
		shouldNotify = true
	} else {
		for _, t := range notificationTypes {
			if t == notificationType {
				shouldNotify = true
				break
			}
		}
	}

	if !shouldNotify {
		// ユーザーが通知タイプを無効にしている場合は通知しない
		return nil
	}

	// 通知データを準備
	data, err := s.prepareNotificationData(ctx, notification)
	if err != nil {
		return err
	}

	// テンプレートを取得
	template, err := s.notificationTemplateRepo.GetByType(ctx, notificationType)
	if err != nil {
		return err
	}

	if template == nil {
		return fmt.Errorf("notification template not found for type: %s", notificationType)
	}

	// 非同期で通知を送信（ここではgoroutineで簡略化していますが、実際にはメッセージキューなどを使用するとよいでしょう）
	go func() {
		// ブラウザプッシュ通知
		if userSettings.PushNotification {
			// タイトルと本文をテンプレートから生成
			title, body := s.renderTemplate(template.TitleTemplate, template.BodyTemplate, data)

			// プッシュ通知を送信
			s.sendPushNotification(userID, title, body, data.Source.URL)
		}

		// Eメール通知
		if userSettings.EmailNotification {
			// 件名と本文をテンプレートから生成
			subject, body := s.renderTemplate(template.EmailSubjectTemplate, template.EmailBodyTemplate, data)

			// Eメール通知を送信
			user, _ := s.userRepo.GetByID(ctx, userID)
			if user != nil {
				s.sendEmailNotification(user.Email, subject, body)
			}
		}
	}()

	return nil
}

// prepareNotificationData は通知データを準備します
func (s *NotificationService) prepareNotificationData(ctx context.Context, notification *models.Notification) (*NotificationData, error) {
	data := &NotificationData{}

	// ユーザー情報を取得
	user, err := s.userRepo.GetByID(ctx, notification.UserID)
	if err != nil {
		return nil, err
	}
	data.User.Name = user.Name
	data.User.Email = user.Email

	// アクター情報を取得
	actor, err := s.userRepo.GetByID(ctx, notification.ActorID)
	if err != nil {
		return nil, err
	}
	data.Actor.Name = actor.Name
	data.Actor.Email = actor.Email

	// ソース情報を取得
	data.Source.Type = notification.SourceType

	// ソースのタイトルとURLを取得
	switch notification.SourceType {
	case "issue":
		// Issueリポジトリからデータを取得
		// 本来であればIssueRepositoryを使ってタイトルを取得するべきですが、
		// 簡略化のためここではダミーデータを設定しています
		data.Source.Title = "Issue Title" // 実際には取得したタイトル
		data.Source.URL = fmt.Sprintf("%s/issues/%d", s.baseURL, notification.SourceID)
	case "discussion":
		// Discussionリポジトリからデータを取得
		data.Source.Title = "Discussion Title" // 実際には取得したタイトル
		data.Source.URL = fmt.Sprintf("%s/discussions/%d", s.baseURL, notification.SourceID)
	case "comment":
		// Commentリポジトリからデータを取得
		data.Source.Title = "Comment Title" // 実際には取得したタイトル
		data.Source.URL = fmt.Sprintf("%s/comments/%d", s.baseURL, notification.SourceID)
	default:
		data.Source.Title = "Unknown Source"
		data.Source.URL = s.baseURL
	}

	data.Message = notification.Message

	return data, nil
}

// renderTemplate はテンプレート文字列をレンダリングします
func (s *NotificationService) renderTemplate(titleTmpl, bodyTmpl string, data *NotificationData) (string, string) {
	title := new(bytes.Buffer)
	body := new(bytes.Buffer)

	t1, _ := template.New("title").Parse(titleTmpl)
	t2, _ := template.New("body").Parse(bodyTmpl)

	t1.Execute(title, data)
	t2.Execute(body, data)

	return title.String(), body.String()
}

// WebPushSubscription はWebPush購読情報を表します
type WebPushSubscription struct {
	Endpoint       string `json:"endpoint"`
	ExpirationTime any    `json:"expirationTime"`
	Keys           struct {
		P256DH string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

// AddPushSubscription はブラウザのプッシュ通知サブスクリプションを追加します
func (s *NotificationService) AddPushSubscription(ctx context.Context, userID int64, subscriptionJSON string) error {
	var subscription WebPushSubscription
	if err := json.Unmarshal([]byte(subscriptionJSON), &subscription); err != nil {
		return err
	}

	// サブスクリプションをDBに保存
	pushSubscription := models.NewPushSubscription(
		userID,
		subscription.Endpoint,
		subscription.Keys.P256DH,
		subscription.Keys.Auth,
	)

	return s.pushSubscriptionRepo.Create(ctx, pushSubscription)
}

// RemovePushSubscription はブラウザのプッシュ通知サブスクリプションを削除します
func (s *NotificationService) RemovePushSubscription(ctx context.Context, endpoint string) error {
	return s.pushSubscriptionRepo.DeleteByEndpoint(ctx, endpoint)
}

// sendPushNotification はブラウザプッシュ通知を送信します
func (s *NotificationService) sendPushNotification(userID int64, title, body, url string) error {
	// ユーザーのプッシュサブスクリプションを取得
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	subscriptions, err := s.pushSubscriptionRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// 通知データを準備
	payload := map[string]interface{}{
		"title": title,
		"body":  body,
		"url":   url,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// 各サブスクリプションに通知を送信
	for _, sub := range subscriptions {
		// WebPush通知を送信
		s.sendWebPushNotification(sub.Endpoint, sub.P256dh, sub.Auth, jsonPayload)
	}

	return nil
}

// sendWebPushNotification はWeb Push通知を送信します
func (s *NotificationService) sendWebPushNotification(endpoint, p256dh, auth string, payload []byte) error {
	// WebPush通知を送信
	subscription := webpush.Subscription{
		Endpoint: endpoint,
		Keys: webpush.Keys{
			P256dh: p256dh,
			Auth:   auth,
		},
	}

	// 通知オプションを設定
	options := webpush.Options{
		Subscriber:      s.smtpFrom, // 通知の送信者（通常はメールアドレス）
		VAPIDPublicKey:  s.vapidPublicKey,
		VAPIDPrivateKey: s.vapidPrivateKey,
		TTL:             30,
	}

	// 通知を送信
	_, err := webpush.SendNotification(payload, &subscription, &options)
	return err
}

// sendEmailNotification はEメール通知を送信します
func (s *NotificationService) sendEmailNotification(to, subject, body string) error {
	// SMTPサーバーに接続
	auth := smtp.PlainAuth("", s.smtpUsername, s.smtpPassword, s.smtpHost)

	// メールヘッダーを設定
	headers := map[string]string{
		"From":         s.smtpFrom,
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	// ヘッダーをメッセージに追加
	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	// メールを送信
	err := smtp.SendMail(
		fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort),
		auth,
		s.smtpFrom,
		[]string{to},
		[]byte(message),
	)

	return err
}

// GetNotifications はユーザーの通知一覧を取得します
func (s *NotificationService) GetNotifications(ctx context.Context, userID int64, isRead *bool, page, limit int) ([]*models.Notification, int, error) {
	return s.notificationRepo.ListByUser(ctx, userID, isRead, page, limit)
}

// MarkAsRead は通知を既読状態に更新します
func (s *NotificationService) MarkAsRead(ctx context.Context, id int64) error {
	return s.notificationRepo.MarkAsRead(ctx, id)
}

// MarkAllAsRead はユーザーの全通知を既読状態に更新します
func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID int64) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}

// GetVAPIDPublicKey はVAPID公開鍵を取得します
func (s *NotificationService) GetVAPIDPublicKey() string {
	return s.vapidPublicKey
}

// UpdateUserNotificationSettings はユーザーの通知設定を更新します
func (s *NotificationService) UpdateUserNotificationSettings(
	ctx context.Context,
	userID int64,
	emailNotification *bool,
	pushNotification *bool,
	notificationTypes *string) error {

	// ユーザー設定を取得
	settings, err := s.userSettingsRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// 更新が指定されている場合のみ更新
	if emailNotification != nil {
		settings.UpdateEmailNotification(*emailNotification)
	}

	if pushNotification != nil {
		settings.UpdatePushNotification(*pushNotification)
	}

	if notificationTypes != nil {
		settings.UpdateNotificationTypes(*notificationTypes)
	}

	// 更新した設定を保存
	return s.userSettingsRepo.Update(ctx, settings)
}

// GetUserSettings はユーザーの通知設定を取得します
func (s *NotificationService) GetUserSettings(ctx context.Context, userID int64) (*models.UserSettings, error) {
	return s.userSettingsRepo.GetByUserID(ctx, userID)
}
