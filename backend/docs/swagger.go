package docs

// @title TicketHub API
// @version 1.0
// @description GitHubライクなTicketHub管理システムのRESTful API

// @contact.name API Support
// @contact.url https://github.com/shimauma0312/module-tickethub
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @tag.name issues
// @tag.description Issue関連のエンドポイント

// @tag.name discussions
// @tag.description Discussion関連のエンドポイント

// @tag.name comments
// @tag.description コメント関連のエンドポイント

// @tag.name labels
// @tag.description ラベル関連のエンドポイント

// @tag.name milestones
// @tag.description マイルストーン関連のエンドポイント

// @tag.name assignments
// @tag.description アサイン関連のエンドポイント

// @tag.name markdown
// @tag.description Markdownレンダリング関連のエンドポイント

// @tag.name drafts
// @tag.description ドラフト保存関連のエンドポイント
