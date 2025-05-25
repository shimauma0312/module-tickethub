package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shimauma0312/module-tickethub/backend/config"
	"github.com/shimauma0312/module-tickethub/backend/migrations"
	"github.com/shimauma0312/module-tickethub/backend/models"
	"github.com/shimauma0312/module-tickethub/backend/repositories"
	"github.com/shimauma0312/module-tickethub/backend/repositories/sqlite"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	log.Println("Starting demo data generation...")

	// 設定の読み込み
	dbConfig, err := config.NewDatabaseConfig()
	if err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	// マイグレーションの実行
	if err := migrations.MigrateDB(dbConfig, "up"); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// リポジトリファクトリの作成
	var factory repositories.RepositoryFactory
	switch dbConfig.Type {
	case config.SQLite:
		db, err := config.OpenSQLiteDB(dbConfig.SQLitePath)
		if err != nil {
			log.Fatalf("Failed to open SQLite database: %v", err)
		}
		factory = sqlite.NewRepositoryFactory(db)
	case config.Postgres:
		log.Fatalf("PostgreSQL demo data generation not implemented yet")
	default:
		log.Fatalf("Unsupported database type: %v", dbConfig.Type)
	}
	defer factory.Close()

	// コンテキストの作成
	ctx := context.Background()

	// デモデータの生成
	if err := generateDemoData(ctx, factory); err != nil {
		log.Fatalf("Failed to generate demo data: %v", err)
	}

	log.Println("Demo data generation completed successfully!")
}

func generateDemoData(ctx context.Context, factory repositories.RepositoryFactory) error {
	log.Println("Generating demo data...")

	// ユーザーデータの生成
	users, err := generateUsers(ctx, factory)
	if err != nil {
		return fmt.Errorf("failed to generate users: %w", err)
	}

	// ラベルデータの生成
	labels, err := generateLabels(ctx, factory)
	if err != nil {
		return fmt.Errorf("failed to generate labels: %w", err)
	}

	// マイルストーンデータの生成
	milestones, err := generateMilestones(ctx, factory, users)
	if err != nil {
		return fmt.Errorf("failed to generate milestones: %w", err)
	}

	// Issue データの生成
	if err := generateIssues(ctx, factory, users, labels, milestones); err != nil {
		return fmt.Errorf("failed to generate issues: %w", err)
	}

	// Discussion データの生成
	if err := generateDiscussions(ctx, factory, users, labels); err != nil {
		return fmt.Errorf("failed to generate discussions: %w", err)
	}

	return nil
}

func generateUsers(ctx context.Context, factory repositories.RepositoryFactory) ([]*models.User, error) {
	log.Println("Generating user data...")

	userRepo, err := factory.NewUserRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create user repository: %w", err)
	}

	settingsRepo, err := factory.NewUserSettingsRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create user settings repository: %w", err)
	}

	// デモユーザーデータ
	demoUsers := []struct {
		username string
		email    string
		password string
		fullName string
		isAdmin  bool
	}{
		{
			username: "admin",
			email:    "admin@example.com",
			password: "admin123",
			fullName: "System Administrator",
			isAdmin:  true,
		},
		{
			username: "user1",
			email:    "user1@example.com",
			password: "user123",
			fullName: "Test User 1",
			isAdmin:  false,
		},
		{
			username: "user2",
			email:    "user2@example.com",
			password: "user123",
			fullName: "Test User 2",
			isAdmin:  false,
		},
		{
			username: "developer",
			email:    "dev@example.com",
			password: "dev123",
			fullName: "Developer",
			isAdmin:  false,
		},
		{
			username: "tester",
			email:    "test@example.com",
			password: "test123",
			fullName: "Quality Tester",
			isAdmin:  false,
		},
	}

	users := make([]*models.User, 0, len(demoUsers))

	for _, u := range demoUsers {
		// パスワードのハッシュ化
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}

		user := models.NewUser(u.username, u.email, string(hashedPassword), u.fullName)
		user.SetAdmin(u.isAdmin)

		// ユーザーの作成
		if err := userRepo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		// ユーザー設定の作成
		settings := models.NewUserSettings(user.ID)
		if err := settingsRepo.CreateOrUpdate(ctx, settings); err != nil {
			return nil, fmt.Errorf("failed to create user settings: %w", err)
		}

		users = append(users, user)
		log.Printf("Created user: %s (ID: %d)", user.Username, user.ID)
	}

	return users, nil
}

func generateLabels(ctx context.Context, factory repositories.RepositoryFactory) ([]*models.Label, error) {
	log.Println("Generating label data...")

	labelRepo, err := factory.NewLabelRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create label repository: %w", err)
	}

	// デモラベルデータ
	demoLabels := []struct {
		name        string
		description string
		color       string
		labelType   string
	}{
		{
			name:        "バグ",
			description: "修正が必要な不具合",
			color:       "#FF0000",
			labelType:   "issue",
		},
		{
			name:        "機能要望",
			description: "新機能のリクエスト",
			color:       "#00FF00",
			labelType:   "issue",
		},
		{
			name:        "ドキュメント",
			description: "ドキュメントの変更や追加",
			color:       "#0000FF",
			labelType:   "both",
		},
		{
			name:        "質問",
			description: "質問や疑問点",
			color:       "#FFFF00",
			labelType:   "discussion",
		},
		{
			name:        "解決済み",
			description: "解決済みの課題",
			color:       "#00FFFF",
			labelType:   "both",
		},
		{
			name:        "優先度:高",
			description: "高優先度の課題",
			color:       "#FF00FF",
			labelType:   "issue",
		},
		{
			name:        "優先度:中",
			description: "中優先度の課題",
			color:       "#FFAAAA",
			labelType:   "issue",
		},
		{
			name:        "優先度:低",
			description: "低優先度の課題",
			color:       "#AAFFAA",
			labelType:   "issue",
		},
		{
			name:        "議論",
			description: "議論が必要なトピック",
			color:       "#AAAAFF",
			labelType:   "discussion",
		},
		{
			name:        "アナウンス",
			description: "お知らせや告知",
			color:       "#FFFFAA",
			labelType:   "discussion",
		},
	}

	labels := make([]*models.Label, 0, len(demoLabels))

	for _, l := range demoLabels {
		label := models.NewLabel(l.name, l.description, l.color, l.labelType)

		// ラベルの作成
		if err := labelRepo.Create(ctx, label); err != nil {
			return nil, fmt.Errorf("failed to create label: %w", err)
		}

		labels = append(labels, label)
		log.Printf("Created label: %s (ID: %d)", label.Name, label.ID)
	}

	return labels, nil
}

func generateMilestones(ctx context.Context, factory repositories.RepositoryFactory, users []*models.User) ([]*models.Milestone, error) {
	log.Println("Generating milestone data...")

	milestoneRepo, err := factory.NewMilestoneRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create milestone repository: %w", err)
	}

	// 作成者をランダムに選択
	getRandomUser := func() *models.User {
		return users[rand.Intn(len(users))]
	}

	now := time.Now()

	// デモマイルストーンデータ
	demoMilestones := []struct {
		title       string
		description string
		dueDate     time.Time
		status      string
	}{
		{
			title:       "バージョン 1.0",
			description: "初回リリース",
			dueDate:     now.Add(30 * 24 * time.Hour),
			status:      "open",
		},
		{
			title:       "バージョン 1.1",
			description: "バグ修正と小さな機能追加",
			dueDate:     now.Add(60 * 24 * time.Hour),
			status:      "open",
		},
		{
			title:       "バージョン 2.0",
			description: "メジャーアップデート",
			dueDate:     now.Add(90 * 24 * time.Hour),
			status:      "open",
		},
		{
			title:       "ドキュメント整備",
			description: "ユーザーマニュアルとAPI仕様書の整備",
			dueDate:     now.Add(45 * 24 * time.Hour),
			status:      "open",
		},
		{
			title:       "アルファテスト",
			description: "内部テスターによる検証",
			dueDate:     now.Add(-15 * 24 * time.Hour),
			status:      "closed",
		},
	}

	milestones := make([]*models.Milestone, 0, len(demoMilestones))

	for _, m := range demoMilestones {
		creator := getRandomUser()
		milestone := models.NewMilestone(m.title, m.description, m.dueDate, creator.ID)

		if m.status == "closed" {
			milestone.Close()
		}

		// マイルストーンの作成
		if err := milestoneRepo.Create(ctx, milestone); err != nil {
			return nil, fmt.Errorf("failed to create milestone: %w", err)
		}

		milestones = append(milestones, milestone)
		log.Printf("Created milestone: %s (ID: %d)", milestone.Title, milestone.ID)
	}

	return milestones, nil
}

func generateIssues(ctx context.Context, factory repositories.RepositoryFactory, users []*models.User, labels []*models.Label, milestones []*models.Milestone) error {
	log.Println("Generating issue data...")

	issueRepo, err := factory.NewIssueRepository()
	if err != nil {
		return fmt.Errorf("failed to create issue repository: %w", err)
	}

	commentRepo, err := factory.NewCommentRepository()
	if err != nil {
		return fmt.Errorf("failed to create comment repository: %w", err)
	}

	// 作成者をランダムに選択
	getRandomUser := func() *models.User {
		return users[rand.Intn(len(users))]
	}

	// 担当者をランダムに選択（nullの場合もあり）
	getRandomAssignee := func() *models.User {
		if rand.Intn(10) < 3 { // 30%の確率でnull
			return nil
		}
		return users[rand.Intn(len(users))]
	}

	// ラベルをランダムに選択
	getRandomLabels := func() []*models.Label {
		var issueLabels []*models.Label
		numLabels := rand.Intn(3) + 1 // 1~3個のラベル

		// issue または both タイプのラベルのみ選択
		var validLabels []*models.Label
		for _, l := range labels {
			if l.Type == "issue" || l.Type == "both" {
				validLabels = append(validLabels, l)
			}
		}

		for i := 0; i < numLabels && i < len(validLabels); i++ {
			label := validLabels[rand.Intn(len(validLabels))]
			// 重複チェック
			duplicate := false
			for _, l := range issueLabels {
				if l.ID == label.ID {
					duplicate = true
					break
				}
			}
			if !duplicate {
				issueLabels = append(issueLabels, label)
			}
		}
		return issueLabels
	}

	// マイルストーンをランダムに選択（nullの場合もあり）
	getRandomMilestone := func() *models.Milestone {
		if rand.Intn(10) < 4 { // 40%の確率でnull
			return nil
		}
		return milestones[rand.Intn(len(milestones))]
	}

	// デモIssueデータ
	demoIssues := []struct {
		title       string
		body        string
		status      string
		numComments int
	}{
		{
			title:       "ログイン画面でエラーが発生する",
			body:        "特定の条件下でログイン時にエラーが発生します。\n\n再現手順:\n1. ユーザー名に特殊文字を入力\n2. ログインボタンをクリック\n\n期待される動作: エラーメッセージの表示\n実際の動作: 画面がフリーズ",
			status:      "open",
			numComments: 3,
		},
		{
			title:       "検索機能の実装",
			body:        "全文検索機能を実装する必要があります。\n\n要件:\n- FTS5による全文検索\n- タイトルと本文の検索\n- 複数キーワードによる検索\n- 検索結果のハイライト表示",
			status:      "open",
			numComments: 2,
		},
		{
			title:       "UI/UXの改善",
			body:        "モバイル端末での表示が最適化されていません。\n\n改善点:\n- レスポンシブデザインの適用\n- タッチ操作の最適化\n- フォントサイズの調整",
			status:      "open",
			numComments: 5,
		},
		{
			title:       "パフォーマンス最適化",
			body:        "大量のデータ表示時にパフォーマンスが低下します。\n\n問題箇所:\n- リスト表示時の遅延\n- メモリ使用量の増加\n\n対策案:\n- 仮想スクロールの実装\n- ページネーションの改善",
			status:      "closed",
			numComments: 7,
		},
		{
			title:       "セキュリティ脆弱性の修正",
			body:        "XSS脆弱性が発見されました。\n\n影響範囲:\n- コメント入力フォーム\n- プロフィール編集画面\n\n対策:\n- 入力値のサニタイズ\n- CSP（Content Security Policy）の実装",
			status:      "open",
			numComments: 4,
		},
		{
			title:       "ドキュメントの更新",
			body:        "API仕様書が古くなっています。\n\n更新内容:\n- 新しいエンドポイントの追加\n- リクエスト/レスポンスの例の更新\n- エラーコードの説明",
			status:      "closed",
			numComments: 1,
		},
		{
			title:       "通知機能の実装",
			body:        "ユーザーへの通知機能を実装します。\n\n要件:\n- ブラウザプッシュ通知\n- メール通知\n- 通知設定画面\n- 既読/未読管理",
			status:      "open",
			numComments: 3,
		},
		{
			title:       "テストカバレッジの向上",
			body:        "現在のテストカバレッジが低いです。\n\n目標:\n- ユニットテストのカバレッジ80%以上\n- E2Eテストの追加\n- CI/CDパイプラインへの統合",
			status:      "open",
			numComments: 2,
		},
		{
			title:       "多言語対応",
			body:        "アプリケーションを多言語対応にします。\n\n対応言語:\n- 英語（デフォルト）\n- 日本語\n- 中国語（簡体字）\n\n実装方法:\n- i18nライブラリの導入\n- 言語リソースファイルの作成",
			status:      "open",
			numComments: 0,
		},
		{
			title:       "データバックアップ機能",
			body:        "定期的なデータバックアップ機能を実装します。\n\n要件:\n- 自動バックアップ（日次/週次/月次）\n- バックアップの暗号化\n- リストア機能\n- バックアップ履歴管理",
			status:      "closed",
			numComments: 6,
		},
	}

	for _, i := range demoIssues {
		creator := getRandomUser()
		assignee := getRandomAssignee()
		milestone := getRandomMilestone()
		issueLabels := getRandomLabels()

		issue := models.NewIssue(i.title, i.body, creator.ID)

		if assignee != nil {
			issue.AssigneeID = assignee.ID
		}

		if milestone != nil {
			issue.MilestoneID = milestone.ID
		}

		for _, label := range issueLabels {
			issue.AddLabel(label.Name)
		}

		if i.status == "closed" {
			issue.Close()
		}

		// Issueの作成
		if err := issueRepo.Create(ctx, issue); err != nil {
			return fmt.Errorf("failed to create issue: %w", err)
		}

		log.Printf("Created issue: %s (ID: %d)", issue.Title, issue.ID)

		// コメントの生成
		for j := 0; j < i.numComments; j++ {
			commentCreator := getRandomUser()
			commentTime := issue.CreatedAt.Add(time.Duration(rand.Intn(100)) * time.Hour)

			commentBody := ""
			switch j % 3 {
			case 0:
				commentBody = fmt.Sprintf("これは問題ですね。詳細を確認します。 @%s さん何か情報はありますか？", users[rand.Intn(len(users))].Username)
			case 1:
				commentBody = "再現できました。修正に取り組みます。"
			case 2:
				commentBody = fmt.Sprintf("修正PRを作成しました。レビューをお願いします。\n\n```\nfunc fix() {\n  // 問題の修正コード\n}\n```\n\n@%s", users[rand.Intn(len(users))].Username)
			}

			comment := models.NewComment(commentBody, commentCreator.ID, issue.ID, "issue")
			comment.CreatedAt = commentTime
			comment.UpdatedAt = commentTime

			// コメントの作成
			if err := commentRepo.Create(ctx, comment); err != nil {
				return fmt.Errorf("failed to create comment: %w", err)
			}

			log.Printf("Created comment for issue %d (Comment ID: %d)", issue.ID, comment.ID)
		}
	}

	return nil
}

func generateDiscussions(ctx context.Context, factory repositories.RepositoryFactory, users []*models.User, labels []*models.Label) error {
	log.Println("Generating discussion data...")

	discussionRepo, err := factory.NewDiscussionRepository()
	if err != nil {
		return fmt.Errorf("failed to create discussion repository: %w", err)
	}

	commentRepo, err := factory.NewCommentRepository()
	if err != nil {
		return fmt.Errorf("failed to create comment repository: %w", err)
	}

	// 作成者をランダムに選択
	getRandomUser := func() *models.User {
		return users[rand.Intn(len(users))]
	}

	// ラベルをランダムに選択
	getRandomLabels := func() []*models.Label {
		var discussionLabels []*models.Label
		numLabels := rand.Intn(3) + 1 // 1~3個のラベル

		// discussion または both タイプのラベルのみ選択
		var validLabels []*models.Label
		for _, l := range labels {
			if l.Type == "discussion" || l.Type == "both" {
				validLabels = append(validLabels, l)
			}
		}

		for i := 0; i < numLabels && i < len(validLabels); i++ {
			label := validLabels[rand.Intn(len(validLabels))]
			// 重複チェック
			duplicate := false
			for _, l := range discussionLabels {
				if l.ID == label.ID {
					duplicate = true
					break
				}
			}
			if !duplicate {
				discussionLabels = append(discussionLabels, label)
			}
		}
		return discussionLabels
	}

	// デモDiscussionデータ
	demoDiscussions := []struct {
		title       string
		body        string
		category    string
		status      string
		numComments int
	}{
		{
			title:       "チーム開発のベストプラクティス",
			body:        "効率的なチーム開発のためのベストプラクティスについて議論しましょう。\n\n話題:\n- ブランチ戦略\n- コードレビュープロセス\n- ペアプログラミング\n- スクラム/カンバンの活用",
			category:    "general",
			status:      "open",
			numComments: 5,
		},
		{
			title:       "アーキテクチャ設計について",
			body:        "現在のアーキテクチャに関する質問です。\n\nマイクロサービスとモノリスのどちらが適していますか？それぞれのメリット・デメリットを教えてください。",
			category:    "question",
			status:      "answered",
			numComments: 3,
		},
		{
			title:       "新機能の提案: ダークモード",
			body:        "ダークモード機能を追加してはどうでしょうか？\n\n実装案:\n- CSSカスタム変数の活用\n- ユーザー設定での切り替え\n- システム設定との連動\n- アニメーション効果",
			category:    "idea",
			status:      "open",
			numComments: 7,
		},
		{
			title:       "リリースノート v1.0.0",
			body:        "v1.0.0がリリースされました！\n\n主な変更点:\n- 全文検索機能の追加\n- パフォーマンスの最適化\n- UIの刷新\n- バグ修正多数\n\nフィードバックをお待ちしています。",
			category:    "announcement",
			status:      "closed",
			numComments: 2,
		},
		{
			title:       "開発環境のセットアップガイド",
			body:        "開発環境のセットアップ手順をまとめました。\n\n手順:\n1. リポジトリのクローン\n2. 依存関係のインストール\n3. 環境変数の設定\n4. ビルドとテスト実行\n\n問題があれば教えてください。",
			category:    "general",
			status:      "open",
			numComments: 4,
		},
		{
			title:       "パフォーマンスチューニングのヒント",
			body:        "アプリケーションのパフォーマンスを向上させるヒントを共有しましょう。\n\n私の経験から:\n- インデックスの適切な設定\n- キャッシュの活用\n- 非同期処理の導入\n- バンドルサイズの最適化",
			category:    "general",
			status:      "open",
			numComments: 6,
		},
		{
			title:       "CI/CDパイプラインの改善案",
			body:        "現在のCI/CDパイプラインに問題があります。\n\n問題点:\n- ビルド時間が長い\n- テスト実行が不安定\n- デプロイの失敗が多い\n\n改善案を募集します。",
			category:    "question",
			status:      "open",
			numComments: 3,
		},
	}

	for _, d := range demoDiscussions {
		creator := getRandomUser()
		discussionLabels := getRandomLabels()

		discussion := models.NewDiscussion(d.title, d.body, d.category, creator.ID)

		for _, label := range discussionLabels {
			discussion.AddLabel(label.Name)
		}

		if d.status == "closed" {
			discussion.Close()
		} else if d.status == "answered" {
			discussion.MarkAsAnswered()
		}

		// Discussionの作成
		if err := discussionRepo.Create(ctx, discussion); err != nil {
			return fmt.Errorf("failed to create discussion: %w", err)
		}

		log.Printf("Created discussion: %s (ID: %d)", discussion.Title, discussion.ID)

		// コメントの生成
		var parentComment *models.Comment = nil

		for j := 0; j < d.numComments; j++ {
			commentCreator := getRandomUser()
			commentTime := discussion.CreatedAt.Add(time.Duration(rand.Intn(100)) * time.Hour)

			commentBody := ""
			isReply := rand.Intn(3) == 0 && parentComment != nil // 33%の確率で返信（親コメントがある場合）

			if isReply {
				commentBody = fmt.Sprintf("@%s さんのコメントに返信します。\n\n同意します。さらに追加すると...", users[rand.Intn(len(users))].Username)
				comment := models.NewReply(commentBody, commentCreator.ID, discussion.ID, parentComment.ID, "discussion")
				comment.CreatedAt = commentTime
				comment.UpdatedAt = commentTime

				// 返信コメントの作成
				if err := commentRepo.Create(ctx, comment); err != nil {
					return fmt.Errorf("failed to create reply comment: %w", err)
				}

				log.Printf("Created reply comment for discussion %d (Comment ID: %d, Parent: %d)", discussion.ID, comment.ID, parentComment.ID)
			} else {
				switch j % 4 {
				case 0:
					commentBody = fmt.Sprintf("これは良いトピックですね。@%s さんはどう思いますか？", users[rand.Intn(len(users))].Username)
				case 1:
					commentBody = "私の経験では、このアプローチが最も効果的でした。\n\n```\nfunction example() {\n  // サンプルコード\n}\n```"
				case 2:
					commentBody = "別の視点から考えると...\n\n1. まず最初のポイント\n2. 次に重要なこと\n3. 最後にまとめ"
				case 3:
					commentBody = fmt.Sprintf("参考になる記事を見つけました: https://example.com/article\n\n@%s さんもチェックしてみてください。", users[rand.Intn(len(users))].Username)
				}

				comment := models.NewComment(commentBody, commentCreator.ID, discussion.ID, "discussion")
				comment.CreatedAt = commentTime
				comment.UpdatedAt = commentTime

				// コメントの作成
				if err := commentRepo.Create(ctx, comment); err != nil {
					return fmt.Errorf("failed to create comment: %w", err)
				}

				log.Printf("Created comment for discussion %d (Comment ID: %d)", discussion.ID, comment.ID)

				// 次の返信用に親コメントとして保存（50%の確率）
				if rand.Intn(2) == 0 {
					parentComment = comment
				}
			}
		}
	}

	return nil
}
