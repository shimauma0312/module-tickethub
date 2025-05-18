# TicketHub - 軽量 GitHub-ライク チケット管理システム

GitHubのIssue/Discussionに似た操作感で、小〜中規模チームでも自前でホストできる軽量・リアルタイムなチケット管理ツールです。

## 主な特徴

- GitHub ライクな UI/UX
- リアルタイム更新 (WebSocket)
- Markdownによる豊かなコンテンツ表現
- 軽量で導入しやすいインフラ構成
- クロスプラットフォーム (x86/ARM)

## 技術スタック

- **フロントエンド**: Nuxt 3 (Vue 3 + Vite), Pinia, TailwindCSS
- **バックエンド**: Go 1.22 / Gin
- **リアルタイム通信**: WebSocket (SockJS + STOMP)
- **データベース**: SQLite 3 (FTS5)
- **キャッシュ**: Redis (オプション)
- **コンテナ**: Docker

## クイックスタート

### 前提条件

- Docker と Docker Compose がインストールされていること

### 1. リポジトリをクローン

```bash
git clone https://github.com/shimauma0312/module-tickethub.git
cd module-tickethub
```

### 2. 環境変数の設定

```bash
cp .env.example .env
# 必要に応じて .env ファイルを編集
```

### 3. アプリケーションの起動

```bash
docker compose up
```

これだけです！ブラウザで http://localhost:3000 にアクセスして、TicketHubを利用開始できます。

## 開発環境

### CI/CD

このプロジェクトには自動テスト、ビルド、デプロイのためのCI/CDパイプラインが実装されています。
詳細は [CI/CD ドキュメント](./DOCS/CI-CD.md) を参照してください。

GitHub Actionsを使用して以下の自動化が行われます:
- プルリクエスト時の自動テストとコード品質チェック
- メインブランチへのマージ時の自動デプロイ
- セキュリティスキャンとコード品質の監視

### フロントエンド開発

```bash
cd frontend
yarn install
yarn dev
```

または npm を使用する場合：

```bash
cd frontend
npm install
npm run dev
```

### バックエンド開発

```bash
cd backend
go mod tidy
go run main.go
```

## 本番環境

本番環境では、最適化されたDockerイメージを使用します：

```bash
docker compose -f docker-compose.prod.yml up -d
```

## コントリビュート

バグ報告、機能リクエスト、プルリクエストは大歓迎です！
コントリビュートの前に、プロジェクトの行動規範とコントリビューションガイドを確認してください。

## ライセンス

このプロジェクトは MIT ライセンスの下で公開されています。
