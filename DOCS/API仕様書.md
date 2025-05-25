# TicketHub API仕様書

## 概要

TicketHubは、GitHubライクなチケット管理システムのRESTful APIです。本APIを使用することで、Issue・Discussion管理、コメント、ラベル、マイルストーン、アサイン機能などを実装できます。

## ベースURL

```
http://localhost:8080/api/v1
```

## 認証

ほとんどのAPIエンドポイントでは認証が必要です。認証には、Bearer認証方式を使用します。

```
Authorization: Bearer {token}
```

`{token}`は、ログインAPIで取得したJWTトークンに置き換えてください。

## 共通レスポンス

### 成功レスポンス

リクエストが成功すると、適切なHTTPステータスコードと共に以下のようなレスポンスが返されます：

```json
{
  "data": {...}  // リクエストに応じたデータ
}
```

または、リスト取得の場合：

```json
{
  "items": [...],  // アイテムの配列
  "total": 100,    // 全件数
  "page": 1,       // 現在のページ
  "limit": 10      // 1ページあたりの件数
}
```

### エラーレスポンス

リクエストが失敗すると、適切なHTTPステータスコードと共に以下のようなレスポンスが返されます：

```json
{
  "error": "エラーメッセージ"
}
```

## APIエンドポイント

### 認証

#### ユーザー登録

```
POST /auth/register
```

**リクエスト**

```json
{
  "username": "user123",
  "email": "user@example.com",
  "password": "password123",
  "full_name": "Test User"
}
```

**レスポンス**

```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "username": "user123",
    "email": "user@example.com"
  }
}
```

#### ログイン

```
POST /auth/login
```

**リクエスト**

```json
{
  "username_or_email": "user123",
  "password": "password123"
}
```

**レスポンス**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 1800,
  "token_type": "Bearer"
}
```

### Issue

#### Issue一覧の取得

```
GET /issues
```

**クエリパラメータ**

- `page`: ページ番号（デフォルト: 1）
- `limit`: 1ページあたりの件数（デフォルト: 10）
- `status`: ステータス（open/closed、デフォルト: open）
- `label`: ラベル名
- `assignee`: 担当者ID
- `milestone`: マイルストーンID

**レスポンス**

```json
{
  "issues": [
    {
      "id": 1,
      "title": "Issue Title",
      "body": "Issue description",
      "status": "open",
      "labels": ["bug", "critical"],
      "assignee_id": 2,
      "creator_id": 1,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-02T00:00:00Z",
      "is_draft": false,
      "milestone_id": 1
    }
  ],
  "total": 100,
  "page": 1,
  "limit": 10
}
```

#### 特定のIssueの取得

```
GET /issues/:id
```

**レスポンス**

```json
{
  "id": 1,
  "title": "Issue Title",
  "body": "Issue description",
  "status": "open",
  "labels": ["bug", "critical"],
  "assignee_id": 2,
  "creator_id": 1,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-02T00:00:00Z",
  "is_draft": false,
  "milestone_id": 1
}
```

#### 新規Issue作成

```
POST /issues
```

**リクエスト**

```json
{
  "title": "New Issue",
  "body": "Issue description",
  "labels": ["bug", "critical"],
  "assignee_id": 2,
  "milestone_id": 1,
  "is_draft": false
}
```

**レスポンス**

```json
{
  "id": 1,
  "title": "New Issue",
  "body": "Issue description",
  "status": "open",
  "labels": ["bug", "critical"],
  "assignee_id": 2,
  "creator_id": 1,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z",
  "is_draft": false,
  "milestone_id": 1
}
```

#### Issue更新

```
PUT /issues/:id
```

**リクエスト**

```json
{
  "title": "Updated Issue",
  "body": "Updated description",
  "labels": ["bug", "critical", "high-priority"],
  "assignee_id": 3,
  "milestone_id": 2,
  "is_draft": false
}
```

**レスポンス**

```json
{
  "id": 1,
  "title": "Updated Issue",
  "body": "Updated description",
  "status": "open",
  "labels": ["bug", "critical", "high-priority"],
  "assignee_id": 3,
  "creator_id": 1,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-02T00:00:00Z",
  "is_draft": false,
  "milestone_id": 2
}
```

#### Issue削除

```
DELETE /issues/:id
```

**レスポンス**

```json
{
  "message": "Issue deleted successfully"
}
```

#### Issueステータス更新

```
PATCH /issues/:id/status
```

**リクエスト**

```json
{
  "status": "closed"
}
```

**レスポンス**

```json
{
  "id": 1,
  "title": "Issue Title",
  "status": "closed",
  "updated_at": "2023-01-03T00:00:00Z",
  ...
}
```

#### Issueドラフト状態更新

```
PATCH /issues/:id/draft
```

**リクエスト**

```json
{
  "is_draft": true
}
```

**レスポンス**

```json
{
  "id": 1,
  "title": "Issue Title",
  "is_draft": true,
  "updated_at": "2023-01-03T00:00:00Z",
  ...
}
```

#### Issue検索

```
GET /issues/search
```

**クエリパラメータ**

- `q`: 検索クエリ（必須）
- `page`: ページ番号（デフォルト: 1）
- `limit`: 1ページあたりの件数（デフォルト: 10）

**レスポンス**

```json
{
  "issues": [...],
  "total": 50,
  "page": 1,
  "limit": 10,
  "query": "バグ"
}
```

### Discussion

#### Discussion一覧の取得

```
GET /discussions
```

**クエリパラメータ**

- `page`: ページ番号（デフォルト: 1）
- `limit`: 1ページあたりの件数（デフォルト: 10）
- `status`: ステータス（open/closed/answered、デフォルト: open）
- `category`: カテゴリ（general/question/announcement/idea）
- `label`: ラベル名

**レスポンス**

```json
{
  "discussions": [
    {
      "id": 1,
      "title": "Discussion Title",
      "body": "Discussion content",
      "status": "open",
      "category": "question",
      "labels": ["feature", "discussion"],
      "creator_id": 1,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-02T00:00:00Z",
      "is_draft": false
    }
  ],
  "total": 100,
  "page": 1,
  "limit": 10
}
```

#### 特定のDiscussionの取得

```
GET /discussions/:id
```

**レスポンス**

```json
{
  "id": 1,
  "title": "Discussion Title",
  "body": "Discussion content",
  "status": "open",
  "category": "question",
  "labels": ["feature", "discussion"],
  "creator_id": 1,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-02T00:00:00Z",
  "is_draft": false
}
```

#### 新規Discussion作成

```
POST /discussions
```

**リクエスト**

```json
{
  "title": "New Discussion",
  "body": "Discussion content",
  "category": "question",
  "labels": ["feature", "discussion"],
  "is_draft": false
}
```

**レスポンス**

```json
{
  "id": 1,
  "title": "New Discussion",
  "body": "Discussion content",
  "status": "open",
  "category": "question",
  "labels": ["feature", "discussion"],
  "creator_id": 1,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z",
  "is_draft": false
}
```

#### Discussion更新

```
PUT /discussions/:id
```

**リクエスト**

```json
{
  "title": "Updated Discussion",
  "body": "Updated content",
  "category": "general",
  "labels": ["feature", "discussion", "important"],
  "is_draft": false
}
```

**レスポンス**

```json
{
  "id": 1,
  "title": "Updated Discussion",
  "body": "Updated content",
  "status": "open",
  "category": "general",
  "labels": ["feature", "discussion", "important"],
  "creator_id": 1,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-02T00:00:00Z",
  "is_draft": false
}
```

#### Discussion削除

```
DELETE /discussions/:id
```

**レスポンス**

```json
{
  "message": "Discussion deleted successfully"
}
```

#### Discussionステータス更新

```
PATCH /discussions/:id/status
```

**リクエスト**

```json
{
  "status": "answered"
}
```

**レスポンス**

```json
{
  "id": 1,
  "title": "Discussion Title",
  "status": "answered",
  "updated_at": "2023-01-03T00:00:00Z",
  ...
}
```

#### Discussionドラフト状態更新

```
PATCH /discussions/:id/draft
```

**リクエスト**

```json
{
  "is_draft": true
}
```

**レスポンス**

```json
{
  "id": 1,
  "title": "Discussion Title",
  "is_draft": true,
  "updated_at": "2023-01-03T00:00:00Z",
  ...
}
```

#### Discussion検索

```
GET /discussions/search
```

**クエリパラメータ**

- `q`: 検索クエリ（必須）
- `page`: ページ番号（デフォルト: 1）
- `limit`: 1ページあたりの件数（デフォルト: 10）

**レスポンス**

```json
{
  "discussions": [...],
  "total": 25,
  "page": 1,
  "limit": 10,
  "query": "機能"
}
```

### コメント

#### コメント一覧の取得

```
GET /:target_type/:target_id/comments
```

**パラメータ**

- `target_type`: コメント対象タイプ（issue/discussion）
- `target_id`: コメント対象ID

**クエリパラメータ**

- `page`: ページ番号（デフォルト: 1）
- `limit`: 1ページあたりの件数（デフォルト: 20）

**レスポンス**

```json
{
  "comments": [
    {
      "id": 1,
      "body": "Comment text",
      "creator_id": 2,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z",
      "type": "issue",
      "target_id": 1,
      "is_edited": false
    }
  ],
  "total": 50,
  "page": 1,
  "limit": 20
}
```

#### 特定のコメントの取得

```
GET /comments/:id
```

**レスポンス**

```json
{
  "id": 1,
  "body": "Comment text",
  "creator_id": 2,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z",
  "type": "issue",
  "target_id": 1,
  "is_edited": false
}
```

#### 新規コメント作成

```
POST /:target_type/:target_id/comments
```

**パラメータ**

- `target_type`: コメント対象タイプ（issue/discussion）
- `target_id`: コメント対象ID

**リクエスト**

```json
{
  "body": "New comment"
}
```

**レスポンス**

```json
{
  "id": 1,
  "body": "New comment",
  "creator_id": 2,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z",
  "type": "issue",
  "target_id": 1,
  "is_edited": false
}
```

#### 返信コメント作成

```
POST /:target_type/:target_id/comments/reply
```

**パラメータ**

- `target_type`: コメント対象タイプ（issue/discussion）
- `target_id`: コメント対象ID

**リクエスト**

```json
{
  "body": "Reply comment",
  "parent_comment_id": 1
}
```

**レスポンス**

```json
{
  "id": 2,
  "body": "Reply comment",
  "creator_id": 2,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z",
  "type": "reply",
  "target_id": 1,
  "parent_comment_id": 1,
  "is_edited": false
}
```

#### コメント更新

```
PUT /comments/:id
```

**リクエスト**

```json
{
  "body": "Updated comment"
}
```

**レスポンス**

```json
{
  "id": 1,
  "body": "Updated comment",
  "creator_id": 2,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-02T00:00:00Z",
  "type": "issue",
  "target_id": 1,
  "is_edited": true
}
```

#### コメント削除

```
DELETE /comments/:id
```

**レスポンス**

```json
{
  "message": "Comment deleted successfully"
}
```

#### 返信コメント一覧の取得

```
GET /comments/:comment_id/replies
```

**クエリパラメータ**

- `page`: ページ番号（デフォルト: 1）
- `limit`: 1ページあたりの件数（デフォルト: 20）

**レスポンス**

```json
{
  "comments": [
    {
      "id": 2,
      "body": "Reply comment",
      "creator_id": 2,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z",
      "type": "reply",
      "target_id": 1,
      "parent_comment_id": 1,
      "is_edited": false
    }
  ],
  "total": 5,
  "page": 1,
  "limit": 20
}
```

### ラベル

#### ラベル一覧の取得

```
GET /labels
```

**クエリパラメータ**

- `page`: ページ番号（デフォルト: 1）
- `limit`: 1ページあたりの件数（デフォルト: 50）
- `type`: ラベルタイプ（issue/discussion/both）

**レスポンス**

```json
{
  "labels": [
    {
      "id": 1,
      "name": "bug",
      "description": "Bug reports",
      "color": "#ff0000",
      "type": "issue",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "total": 30,
  "page": 1,
  "limit": 50
}
```

#### 特定のラベルの取得

```
GET /labels/:id
```

**レスポンス**

```json
{
  "id": 1,
  "name": "bug",
  "description": "Bug reports",
  "color": "#ff0000",
  "type": "issue",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

#### 新規ラベル作成（管理者のみ）

```
POST /labels
```

**リクエスト**

```json
{
  "name": "feature",
  "description": "Feature requests",
  "color": "#00ff00",
  "type": "both"
}
```

**レスポンス**

```json
{
  "id": 2,
  "name": "feature",
  "description": "Feature requests",
  "color": "#00ff00",
  "type": "both",
  "created_at": "2023-01-02T00:00:00Z",
  "updated_at": "2023-01-02T00:00:00Z"
}
```

#### ラベル更新（管理者のみ）

```
PUT /labels/:id
```

**リクエスト**

```json
{
  "name": "enhancement",
  "description": "Feature enhancements",
  "color": "#00ffff",
  "type": "both"
}
```

**レスポンス**

```json
{
  "id": 2,
  "name": "enhancement",
  "description": "Feature enhancements",
  "color": "#00ffff",
  "type": "both",
  "created_at": "2023-01-02T00:00:00Z",
  "updated_at": "2023-01-03T00:00:00Z"
}
```

#### ラベル削除（管理者のみ）

```
DELETE /labels/:id
```

**レスポンス**

```json
{
  "message": "Label deleted successfully"
}
```

### マイルストーン

#### マイルストーン一覧の取得

```
GET /milestones
```

**クエリパラメータ**

- `page`: ページ番号（デフォルト: 1）
- `limit`: 1ページあたりの件数（デフォルト: 10）
- `status`: ステータス（open/closed、デフォルト: open）

**レスポンス**

```json
{
  "milestones": [
    {
      "id": 1,
      "title": "v1.0",
      "description": "First stable release",
      "due_date": "2023-12-31T00:00:00Z",
      "status": "open",
      "creator_id": 1,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "total": 5,
  "page": 1,
  "limit": 10
}
```

#### 特定のマイルストーンの取得

```
GET /milestones/:id
```

**レスポンス**

```json
{
  "id": 1,
  "title": "v1.0",
  "description": "First stable release",
  "due_date": "2023-12-31T00:00:00Z",
  "status": "open",
  "creator_id": 1,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

#### 新規マイルストーン作成

```
POST /milestones
```

**リクエスト**

```json
{
  "title": "v1.1",
  "description": "Bugfix release",
  "due_date": "2024-03-31"
}
```

**レスポンス**

```json
{
  "id": 2,
  "title": "v1.1",
  "description": "Bugfix release",
  "due_date": "2024-03-31T00:00:00Z",
  "status": "open",
  "creator_id": 1,
  "created_at": "2023-01-02T00:00:00Z",
  "updated_at": "2023-01-02T00:00:00Z"
}
```

#### マイルストーン更新

```
PUT /milestones/:id
```

**リクエスト**

```json
{
  "title": "v1.1",
  "description": "Bugfix and performance release",
  "due_date": "2024-04-15"
}
```

**レスポンス**

```json
{
  "id": 2,
  "title": "v1.1",
  "description": "Bugfix and performance release",
  "due_date": "2024-04-15T00:00:00Z",
  "status": "open",
  "creator_id": 1,
  "created_at": "2023-01-02T00:00:00Z",
  "updated_at": "2023-01-03T00:00:00Z"
}
```

#### マイルストーン削除

```
DELETE /milestones/:id
```

**レスポンス**

```json
{
  "message": "Milestone deleted successfully"
}
```

#### マイルストーンステータス更新

```
PATCH /milestones/:id/status
```

**リクエスト**

```json
{
  "status": "closed"
}
```

**レスポンス**

```json
{
  "id": 2,
  "title": "v1.1",
  "status": "closed",
  "completed_at": "2023-01-04T00:00:00Z",
  "updated_at": "2023-01-04T00:00:00Z",
  ...
}
```

### アサイン

#### Issue担当者設定

```
PUT /issues/:id/assign
```

**リクエスト**

```json
{
  "assignee_id": 3
}
```

**レスポンス**

```json
{
  "message": "Issue assigned successfully",
  "issue_id": 1,
  "assignee_id": 3
}
```

#### Issue担当者削除

```
PUT /issues/:id/unassign
```

**レスポンス**

```json
{
  "message": "Issue unassigned successfully",
  "issue_id": 1
}
```

### Markdown

#### Markdownテキストのレンダリング（JSON）

```
POST /markdown
```

**リクエスト**

```json
{
  "text": "# Title\n\nThis is **bold** text."
}
```

**レスポンス**

```json
{
  "html": "<h1>Title</h1>\n<p>This is <strong>bold</strong> text.</p>"
}
```

#### Markdownテキストのレンダリング（HTML）

```
POST /markdown/raw
```

**リクエスト**

```json
{
  "text": "# Title\n\nThis is **bold** text."
}
```

**レスポンス**

```html
<h1>Title</h1>
<p>This is <strong>bold</strong> text.</p>
```

### ドラフト

#### ドラフト一覧の取得

```
GET /drafts
```

**クエリパラメータ**

- `type`: ドラフトタイプ（issues/discussions/all、デフォルト: all）
- `page`: ページ番号（デフォルト: 1）
- `limit`: 1ページあたりの件数（デフォルト: 10）

**レスポンス**

```json
{
  "issues": [...],
  "issues_total": 3,
  "discussions": [...],
  "discussions_total": 2,
  "page": 1,
  "limit": 10
}
```

#### Issueドラフト保存（新規）

```
POST /drafts/issues
```

**リクエスト**

```json
{
  "title": "Draft Issue",
  "body": "Work in progress...",
  "labels": ["feature"],
  "assignee_id": 2,
  "milestone_id": 1
}
```

**レスポンス**

```json
{
  "id": 5,
  "title": "Draft Issue",
  "body": "Work in progress...",
  "status": "open",
  "labels": ["feature"],
  "assignee_id": 2,
  "creator_id": 1,
  "created_at": "2023-01-05T00:00:00Z",
  "updated_at": "2023-01-05T00:00:00Z",
  "is_draft": true,
  "milestone_id": 1
}
```

#### Issueドラフト保存（更新）

```
PUT /drafts/issues/:id
```

**リクエスト**

```json
{
  "title": "Updated Draft Issue",
  "body": "More progress..."
}
```

**レスポンス**

```json
{
  "id": 5,
  "title": "Updated Draft Issue",
  "body": "More progress...",
  "status": "open",
  "labels": ["feature"],
  "assignee_id": 2,
  "creator_id": 1,
  "created_at": "2023-01-05T00:00:00Z",
  "updated_at": "2023-01-06T00:00:00Z",
  "is_draft": true,
  "milestone_id": 1
}
```

#### Discussionドラフト保存（新規）

```
POST /drafts/discussions
```

**リクエスト**

```json
{
  "title": "Draft Discussion",
  "body": "Work in progress...",
  "category": "general",
  "labels": ["question"]
}
```

**レスポンス**

```json
{
  "id": 3,
  "title": "Draft Discussion",
  "body": "Work in progress...",
  "status": "open",
  "category": "general",
  "labels": ["question"],
  "creator_id": 1,
  "created_at": "2023-01-05T00:00:00Z",
  "updated_at": "2023-01-05T00:00:00Z",
  "is_draft": true
}
```

#### Discussionドラフト保存（更新）

```
PUT /drafts/discussions/:id
```

**リクエスト**

```json
{
  "title": "Updated Draft Discussion",
  "body": "More progress..."
}
```

**レスポンス**

```json
{
  "id": 3,
  "title": "Updated Draft Discussion",
  "body": "More progress...",
  "status": "open",
  "category": "general",
  "labels": ["question"],
  "creator_id": 1,
  "created_at": "2023-01-05T00:00:00Z",
  "updated_at": "2023-01-06T00:00:00Z",
  "is_draft": true
}
```

### 検索

#### コンテンツ検索

```
GET /api/search
```

**クエリパラメータ**

- `query`: 検索キーワード（必須）
- `limit`: 結果の上限数（デフォルト: 20）
- `offset`: 結果のオフセット（デフォルト: 0）
- `labels`: ラベルフィルタ（カンマ区切り、例: `bug,feature`）
- `status`: ステータスフィルタ（`open` / `closed` / `all`、デフォルト: `all`）
- `assignee_id`: 担当者IDフィルタ
- `creator_id`: 作成者IDフィルタ

**レスポンス**

```json
{
  "results": [
    {
      "type": "issue", // "issue" または "comment"
      "id": 1,
      "title": "検索対象のIssueタイトル", // Issueの場合
      "body": "検索対象の本文...",
      "snippet": "...検索キーワードを含む本文の断片...",
      "labels": ["bug", "ui"], // Issueの場合
      "status": "open", // Issueの場合
      "assignee_id": 2,
      "creator_id": 1,
      "created_at": "2023-05-20T10:00:00Z",
      "updated_at": "2023-05-21T11:30:00Z",
      "target_id": 0, // Commentの場合、関連するIssue/DiscussionのID
      "rank": 12.345, // 検索ランキングスコア
      "highlighted": "本文中の<mark>検索キーワード</mark>がハイライトされます"
    }
    // ... 他の検索結果
  ],
  "total_count": 15,
  "current_page": 1,
  "total_pages": 2,
  "query": "検索キーワード"
}
```

**説明**

- Issueのタイトル、本文、およびコメントの本文を対象に全文検索を行います。
- `type`フィールドで結果がIssueかCommentかを示します。
- `snippet`フィールドには、検索キーワード周辺のテキストが表示されます。
- `highlighted`フィールドには、検索キーワードが`<mark>`タグで囲まれた本文（またはタイトル）が表示されます。
- FTS5のランキングアルゴリズムに基づいて関連性の高い順にソートされます。

#### 検索インデックスの再構築（管理者のみ）

```
POST /api/search/rebuild-index
```

**リクエストボディ**

なし

**レスポンス**

```json
{
  "message": "インデックスの再構築が完了しました"
}
```

**説明**

- データベース内のすべてのIssueとCommentを対象に、検索インデックスを再構築します。
- 通常は自動でインデックスが更新されますが、何らかの理由でインデックスが破損した場合や、システム初期構築時に使用します。
- この操作は時間がかかる可能性があるため、注意して実行してください。

### アサイン

#### Issue担当者設定

```
PUT /issues/:id/assign
```

**リクエスト**

```json
{
  "assignee_id": 3
}
```

**レスポンス**

```json
{
  "message": "Issue assigned successfully",
  "issue_id": 1,
  "assignee_id": 3
}
```

#### Issue担当者削除

```
PUT /issues/:id/unassign
```

**レスポンス**

```json
{
  "message": "Issue unassigned successfully",
  "issue_id": 1
}
```
