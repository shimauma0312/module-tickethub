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
	"github.com/shimauma0312/module-tickethub/config"
	"github.com/shimauma0312/module-tickethub/services"
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

	// リポジトリファクトリーの作成
	repoFactory := services.NewRepositoryFactory(db, dbConfig.Type)

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
		v1 := apiGroup.Group("/v1")
		{
			// Issueリポジトリの作成
			issueRepo, err := repoFactory.NewIssueRepository()
			if err != nil {
				log.Fatalf("Failed to create issue repository: %v", err)
			}

			// 各種エンドポイントを設定
			v1.GET("/issues", func(c *gin.Context) {
				// クエリパラメータの取得
				page := 1
				limit := 10

				// データベースから取得
				issues, total, err := issueRepo.List(c, map[string]interface{}{}, page, limit)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"issues": issues,
					"total":  total,
					"page":   page,
					"limit":  limit,
				})
			})
		}
	}

	// サーバー起動
	fmt.Printf("Server starting on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
