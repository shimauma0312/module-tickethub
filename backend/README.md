# TicketHub データモデルとマイグレーション

このリポジトリには、TicketHubアプリケーションのデータモデル設計とデータベースマイグレーション機能が含まれています。

## データモデル

以下のエンティティのデータモデルが実装されています：

- User: ユーザー情報
- UserSettings: ユーザー設定
- Issue: 課題チケット
- Discussion: ディスカッション
- Comment: コメント（課題やディスカッションに対する）
- Label: ラベル
- Milestone: マイルストーン
- Reaction: リアクション（コメントに対する）
- Notification: 通知

## データベース対応

- SQLite 3: デフォルトのデータベース
- PostgreSQL: スケーラビリティのための外部データベース
- 抽象化レイヤー: リポジトリパターンによるデータベース操作の抽象化

## 全文検索対応

- SQLite: FTS5（Full Text Search）による全文検索
- PostgreSQL: tsvector型とGINインデックスによる全文検索

## マイグレーション

マイグレーションファイルは `migrations/{sqlite|postgres}` ディレクトリに配置されています。マイグレーションの実行には以下のコマンドを使用します：

```bash
# マイグレーションを実行（アップ）
go run main.go migrate up

# マイグレーションを元に戻す（ダウン）
go run main.go migrate down
```

## デモデータの生成

デモデータを生成するには以下のコマンドを実行します：

```bash
go run ./cmd/generate_demo_data/main.go
```

このコマンドは以下のデモデータを生成します：

- ユーザー（admin, user1, user2, developer, tester）
- ラベル（バグ, 機能要望, ドキュメント, 質問, 解決済み, 優先度, など）
- マイルストーン（バージョン 1.0, バージョン 1.1, バージョン 2.0, ドキュメント整備, アルファテスト）
- Issue（10件）
- Discussion（7件）
- コメント（各Issueとディスカッションに対する）

## 環境変数の設定

以下の環境変数で動作を制御できます：

```
# データベースタイプ（sqlite, postgres, sqlserver）
DB_TYPE=sqlite

# SQLite設定
SQLITE_DB_PATH=./data/tickethub.db

# PostgreSQL設定
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=tickethub
DB_SSLMODE=disable
```

## リポジトリパターン

アプリケーションは以下のリポジトリインターフェースを提供しています：

- IssueRepository
- UserRepository
- UserSettingsRepository
- LabelRepository
- MilestoneRepository
- DiscussionRepository
- CommentRepository
- ReactionRepository
- NotificationRepository

各インターフェースはSQLiteとPostgreSQLの実装があります。
