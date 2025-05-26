package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/shimauma0312/module-tickethub/backend/api"
	"github.com/shimauma0312/module-tickethub/backend/config"
	_ "github.com/shimauma0312/module-tickethub/backend/docs" // Swaggerドキュメント用
	"github.com/shimauma0312/module-tickethub/backend/services"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 環境変数の読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// サーバー環境の設定
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080" // デフォルトポート
	}

	// データベース設定の読み込み
	dbConfig, err := config.NewDatabaseConfig()
	if err != nil {
		log.Fatalf("Failed to create database config: %v", err)
	}

	// データベース接続の初期化
	db, err := config.InitDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// GORM DBの初期化
	gormDB, err := config.InitGormDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize GORM DB: %v", err)
	}

	// リポジトリファクトリーの作成
	repoFactory := services.NewRepositoryFactory(gormDB)

	// 認証サービスの初期化
	userRepo, err := repoFactory.NewUserRepository()
	if err != nil {
		log.Fatalf("Failed to create user repository: %v", err)
	}

	tokenRepo, err := repoFactory.NewAuthTokenRepository()
	if err != nil {
		log.Fatalf("Failed to create auth token repository: %v", err)
	}

	passwordResetRepo, err := repoFactory.NewPasswordResetRepository()
	if err != nil {
		log.Fatalf("Failed to create password reset repository: %v", err)
	}

	// JWT設定
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "tickethub-jwt-secret-key" // 本番環境では環境変数から取得する
	}

	// 認証サービスの作成
	authService := services.NewAuthService(
		userRepo,
		tokenRepo,
		passwordResetRepo,
		jwtSecret,
	)

	// Ginの設定
	r := gin.Default()

	// CORSの設定
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// ルートハンドラ
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "TicketHub API Server is running",
			"version": "0.1.0",
		})
	})

	// ヘルスチェック用エンドポイント
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// APIルートの設定
	apiGroup := r.Group("/api")
	{
		// CSRF保護トークン生成
		apiGroup.GET("/csrf-token", api.GenerateCSRFToken)

		// Swaggerドキュメント
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

		// 認証ルート
		authGroup := apiGroup.Group("/auth")
		{
			// 認証ハンドラーの作成
			authHandler := api.NewAuthHandler(authService)

			// 認証ルートの設定
			authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh-token", authHandler.RefreshToken)
			authGroup.POST("/password-reset", authHandler.InitiatePasswordReset)
			authGroup.POST("/password-reset/validate", authHandler.ValidatePasswordResetToken)
			authGroup.POST("/password-reset/complete", authHandler.CompletePasswordReset)

			// 認証が必要なルート
			authRequiredGroup := authGroup.Group("/")
			authRequiredGroup.Use(api.AuthMiddleware(authService))
			{
				authRequiredGroup.POST("/logout", authHandler.Logout)
				authRequiredGroup.POST("/logout-all", authHandler.LogoutAll)
				authRequiredGroup.POST("/change-password", authHandler.ChangePassword)
			}
		}

		v1 := apiGroup.Group("/v1")
		{
			// 各種リポジトリの作成
			issueRepo, err := repoFactory.NewIssueRepository()
			if err != nil {
				log.Fatalf("Failed to create issue repository: %v", err)
			}

			discussionRepo, err := repoFactory.NewDiscussionRepository()
			if err != nil {
				log.Fatalf("Failed to create discussion repository: %v", err)
			}

			commentRepo, err := repoFactory.NewCommentRepository()
			if err != nil {
				log.Fatalf("Failed to create comment repository: %v", err)
			}

			labelRepo, err := repoFactory.NewLabelRepository()
			if err != nil {
				log.Fatalf("Failed to create label repository: %v", err)
			}

			milestoneRepo, err := repoFactory.NewMilestoneRepository()
			if err != nil {
				log.Fatalf("Failed to create milestone repository: %v", err)
			}

			reactionRepo, err := repoFactory.NewReactionRepository()
			if err != nil {
				log.Fatalf("Failed to create reaction repository: %v", err)
			}

			userRepo, err := repoFactory.NewUserRepository()
			if err != nil {
				log.Fatalf("Failed to create user repository: %v", err)
			}

			// 管理者機能用リポジトリの作成
			systemSettingsRepo, err := repoFactory.NewSystemSettingsRepository()
			if err != nil {
				log.Fatalf("Failed to create system settings repository: %v", err)
			}

			activityLogRepo, err := repoFactory.NewActivityLogRepository()
			if err != nil {
				log.Fatalf("Failed to create activity log repository: %v", err)
			}

			backupRepo, err := repoFactory.NewBackupRepository()
			if err != nil {
				log.Fatalf("Failed to create backup repository: %v", err)
			}

			// 検索サービスの作成
			searchService, err := repoFactory.NewSearchService()
			if err != nil {
				log.Fatalf("Failed to create search service: %v", err)
			}

			// 各種ハンドラーの作成
			issueHandler := api.NewIssueHandler(issueRepo, labelRepo, milestoneRepo, userRepo)
			discussionHandler := api.NewDiscussionHandler(discussionRepo, labelRepo, userRepo)
			commentHandler := api.NewCommentHandler(commentRepo, issueRepo, discussionRepo, reactionRepo, userRepo)
			labelHandler := api.NewLabelHandler(labelRepo)
			milestoneHandler := api.NewMilestoneHandler(milestoneRepo)
			assignmentHandler := api.NewAssignmentHandler(issueRepo, userRepo)
			markdownHandler := api.NewMarkdownHandler()
			draftHandler := api.NewDraftHandler(issueRepo, discussionRepo)
			searchHandler := api.NewSearchHandler(searchService)

			// 管理者機能用サービスとハンドラーの作成
			activityLogService := services.NewActivityLogService(activityLogRepo)
			backupService := services.NewBackupService(backupRepo, "backups", string(dbConfig.Type), dbConfig.DSN())
			systemMetricsService := services.NewSystemMetricsService(userRepo, issueRepo, discussionRepo, commentRepo, backupRepo)
			adminHandler := api.NewAdminHandler(userRepo, systemSettingsRepo, activityLogService, backupService, systemMetricsService)

			// リポジトリ管理のハンドラー作成
			repoRepo, err := repoFactory.NewRepositoryRepository()
			if err != nil {
				log.Fatalf("Failed to create repository repository: %v", err)
			}
			repositoryHandler := api.NewRepositoryHandler(repoRepo, activityLogService)

			// 認証が必要なルートグループ
			authGroup := v1.Group("/")
			authGroup.Use(api.AuthMiddleware(authService))

			// 管理者権限が必要なルートグループ
			adminGroup := authGroup.Group("/")
			adminGroup.Use(api.AdminMiddleware())

			// Issue関連のエンドポイント
			v1.GET("/issues", issueHandler.ListIssues)
			v1.GET("/issues/:id", issueHandler.GetIssue)
			v1.GET("/issues/search", issueHandler.SearchIssues)
			authGroup.POST("/issues", issueHandler.CreateIssue)
			authGroup.PUT("/issues/:id", issueHandler.UpdateIssue)
			authGroup.DELETE("/issues/:id", issueHandler.DeleteIssue)
			authGroup.PATCH("/issues/:id/status", issueHandler.UpdateIssueStatus)
			authGroup.PATCH("/issues/:id/draft", issueHandler.UpdateIssueDraftStatus)

			// Discussion関連のエンドポイント
			v1.GET("/discussions", discussionHandler.ListDiscussions)
			v1.GET("/discussions/:id", discussionHandler.GetDiscussion)
			v1.GET("/discussions/search", discussionHandler.SearchDiscussions)
			authGroup.POST("/discussions", discussionHandler.CreateDiscussion)
			authGroup.PUT("/discussions/:id", discussionHandler.UpdateDiscussion)
			authGroup.DELETE("/discussions/:id", discussionHandler.DeleteDiscussion)
			authGroup.PATCH("/discussions/:id/status", discussionHandler.UpdateDiscussionStatus)
			authGroup.PATCH("/discussions/:id/draft", discussionHandler.UpdateDiscussionDraftStatus)

			// コメント関連のエンドポイント
			v1.GET("/comments/:id", commentHandler.GetComment)
			v1.GET("/:target_type/:target_id/comments", commentHandler.ListComments)
			v1.GET("/comments/:comment_id/replies", commentHandler.ListReplies)
			authGroup.POST("/:target_type/:target_id/comments", commentHandler.CreateComment)
			authGroup.POST("/:target_type/:target_id/comments/reply", commentHandler.CreateReplyComment)
			authGroup.PUT("/comments/:id", commentHandler.UpdateComment)
			authGroup.DELETE("/comments/:id", commentHandler.DeleteComment)

			// ラベル関連のエンドポイント
			v1.GET("/labels", labelHandler.ListLabels)
			v1.GET("/labels/:id", labelHandler.GetLabel)
			adminGroup.POST("/labels", labelHandler.CreateLabel)
			adminGroup.PUT("/labels/:id", labelHandler.UpdateLabel)
			adminGroup.DELETE("/labels/:id", labelHandler.DeleteLabel)

			// マイルストーン関連のエンドポイント
			v1.GET("/milestones", milestoneHandler.ListMilestones)
			v1.GET("/milestones/:id", milestoneHandler.GetMilestone)
			authGroup.POST("/milestones", milestoneHandler.CreateMilestone)
			authGroup.PUT("/milestones/:id", milestoneHandler.UpdateMilestone)
			authGroup.DELETE("/milestones/:id", milestoneHandler.DeleteMilestone)
			authGroup.PATCH("/milestones/:id/status", milestoneHandler.UpdateMilestoneStatus)

			// アサイン関連のエンドポイント
			authGroup.PUT("/issues/:id/assign", assignmentHandler.AssignIssue)
			authGroup.PUT("/issues/:id/unassign", assignmentHandler.UnassignIssue)

			// Markdown関連のエンドポイント
			v1.POST("/markdown", markdownHandler.RenderMarkdown)
			v1.POST("/markdown/raw", markdownHandler.RenderRawMarkdown)

			// ドラフト関連のエンドポイント
			authGroup.GET("/drafts", draftHandler.ListDrafts)
			authGroup.POST("/drafts/issues", draftHandler.SaveIssueDraft)
			authGroup.PUT("/drafts/issues/:id", draftHandler.SaveIssueDraft)
			authGroup.POST("/drafts/discussions", draftHandler.SaveDiscussionDraft)
			authGroup.PUT("/drafts/discussions/:id", draftHandler.SaveDiscussionDraft)

			// 検索関連のエンドポイント
			searchHandler.RegisterRoutes(r)

			// 管理者専用のエンドポイント
			adminGroup.GET("/users", adminHandler.GetUsers)
			adminGroup.PUT("/users/:id", adminHandler.UpdateUser)

			adminGroup.GET("/settings", adminHandler.GetSystemSettings)
			adminGroup.PUT("/settings", adminHandler.UpdateSystemSettings)

			adminGroup.GET("/activity-logs", adminHandler.GetActivityLogs)
			adminGroup.GET("/metrics", adminHandler.GetSystemMetrics)

			adminGroup.POST("/backups", adminHandler.CreateBackup)
			adminGroup.GET("/backups", adminHandler.GetBackups)
			adminGroup.POST("/backups/:id/restore", adminHandler.RestoreBackup)
			adminGroup.DELETE("/backups/:id", adminHandler.DeleteBackup)

			// リポジトリ管理エンドポイントの追加
			adminGroup.GET("/repositories", repositoryHandler.GetRepositories)
			adminGroup.GET("/repositories/:id", repositoryHandler.GetRepository)
			adminGroup.POST("/repositories", repositoryHandler.CreateRepository)
			adminGroup.PUT("/repositories/:id", repositoryHandler.UpdateRepository)
			adminGroup.DELETE("/repositories/:id", repositoryHandler.DeleteRepository)
		}
	}

	// サーバー起動
	fmt.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
