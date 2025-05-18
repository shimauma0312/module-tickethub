# CI/CD パイプライン

## ワークフロー概要

### 1. CI - Pull Request
プルリクエスト作成時に実行されるワークフロー:
- バックエンドのGoコードのリントとテスト
- フロントエンドのNuxtコードのリントとビルド
- テスト結果とカバレッジレポートの生成

### 2. CI/CD - Deploy
メインブランチへのマージ時に実行されるワークフロー:
- バックエンドとフロントエンドのテストとビルド
- E2Eテストの実行
- 本番環境へのデプロイ

### 3. コード品質 & セキュリティチェック
定期的に実行される品質チェック:
- セキュリティスキャン (Gosec, Trivy)
- コード品質チェック (gofmt, ESLint)
- SonarCloudによる静的解析

## テスト実行方法

### バックエンドテスト
```bash
cd backend
go test -v ./...
```

### フロントエンドテスト
```bash
cd frontend
yarn test
```

または npm を使用する場合：
```bash
cd frontend
npm test
```

### E2Eテスト
```bash
cd tests
yarn install
yarn playwright test
```

または npm を使用する場合：
```bash
cd tests
npm install
npx playwright test
```

## 開発時のベストプラクティス

1. **コミット前の確認**
   - バックエンド: `go test ./...` と `golangci-lint run`
   - フロントエンド: `yarn lint` と `yarn test` （または `npm run lint` と `npm test`）

2. **ブランチ戦略**
   - 機能開発: `feature/機能名`
   - バグ修正: `fix/問題の簡単な説明`
   - リリース: `release/vX.Y.Z`

3. **PRの作成**
   - 明確なタイトルと説明
   - 関連するIssueへのリンク
   - テストが追加されていることを確認

## カバレッジ目標

- バックエンド: 80%以上のコードカバレッジ
- フロントエンド: コンポーネントの80%以上をテスト
- E2E: 主要ユーザーフローをカバー
