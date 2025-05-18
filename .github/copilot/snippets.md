# 参照コードスニペット

このファイルには重要なコードパターンや例が含まれています。実際のコードを書く際は、プロジェクト内の既存コードを参考にしてください。

## よく使う設計パターン

### バックエンド (Go)

#### レイヤ構造
- **Model**: データ構造の定義 (/backend/models/)
- **Repository**: データアクセス層 (/backend/repositories/)
- **Service**: ビジネスロジック (/backend/services/)
- **API**: HTTPハンドラ (/backend/api/)

#### エラーハンドリング
```go
if err != nil {
    // ログ出力
    log.Printf("エラー詳細: %v", err)
    
    // 適切なエラー分類とラップ
    if errors.Is(err, sql.ErrNoRows) {
        return nil, services.ErrNotFound
    }
    return nil, fmt.Errorf("操作に失敗しました: %w", err)
}
```

### フロントエンド (Nuxt/Vue)

#### Composition API パターン
```typescript
// 再利用可能なComposable関数
export function useResource<T>(path: string) {
  const items = ref<T[]>([])
  const loading = ref(false)
  const error = ref<Error | null>(null)
  
  async function fetchItems() {
    loading.value = true
    try {
      const { data } = await useFetch<T[]>(`/api/${path}`)
      items.value = data.value || []
    } catch (err: any) {
      error.value = err
    } finally {
      loading.value = false
    }
  }
  
  return {
    items,
    loading,
    error,
    fetchItems
  }
}
```

## 既存コンポーネント

プロジェクト内には参考になるコード例が多数あります。例えば：

- Issue モデル: `/backend/models/issue.go`
- API ハンドラ例: バックエンドのメインコントローラ
- Vue コンポーネント例: `/frontend/pages/`
- テスト例: `/backend/models/issue_test.go`

実装の際には、既存のコードパターンに合わせて開発を行ってください。
```
