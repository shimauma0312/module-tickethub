<template>
  <div class="discussion-create">
    <div class="discussion-create-header">
      <h1 class="discussion-create-title">新規ディスカッション作成</h1>
    </div>
    
    <div class="discussion-form">
      <div class="form-main">
        <div class="form-group">
          <label for="discussion-title" class="form-label">タイトル</label>
          <Input 
            id="discussion-title" 
            v-model="discussion.title" 
            placeholder="ディスカッションのタイトルを入力" 
            class="discussion-title-input"
          />
        </div>
        
        <div class="form-group">
          <label for="discussion-category" class="form-label">カテゴリ</label>
          <Select 
            id="discussion-category"
            v-model="discussion.categoryId"
            :options="categoryOptions"
            placeholder="カテゴリを選択"
          />
        </div>
        
        <div class="form-group">
          <label for="discussion-body" class="form-label">内容</label>
          <MarkdownEditor 
            id="discussion-body"
            v-model="discussion.content" 
            :placeholder="bodyPlaceholder" 
            :minHeight="400"
            :showPreview="true"
          />
        </div>
      </div>
      
      <div class="form-sidebar">
        <Card class="sidebar-section">
          <h3 class="sidebar-heading">ディスカッション作成のヒント</h3>
          <ul class="tips-list">
            <li>タイトルは具体的かつ簡潔に記述しましょう</li>
            <li>内容は構造化して読みやすくしましょう</li>
            <li>質問の場合は、試したこと・エラー内容を明記しましょう</li>
            <li>関連する資料やリンクがあれば追加しましょう</li>
            <li>適切なカテゴリを選択しましょう</li>
          </ul>
        </Card>
        
        <Card class="sidebar-section">
          <h3 class="sidebar-heading">Markdownの使い方</h3>
          <div class="markdown-tips">
            <div class="markdown-tip">
              <code># 見出し</code>
              <span>見出しを作成</span>
            </div>
            <div class="markdown-tip">
              <code>**太字**</code>
              <span>テキストを強調</span>
            </div>
            <div class="markdown-tip">
              <code>[リンク](URL)</code>
              <span>ハイパーリンクを追加</span>
            </div>
            <div class="markdown-tip">
              <code>- リスト項目</code>
              <span>箇条書きリスト</span>
            </div>
            <div class="markdown-tip">
              <code>```コード```</code>
              <span>コードブロック</span>
            </div>
          </div>
        </Card>
      </div>
      
      <div class="form-actions">
        <Button variant="outline" @click="saveDraft" class="mr-2">
          下書き保存
        </Button>
        <Button variant="primary" @click="submitDiscussion" :disabled="!isFormValid">
          ディスカッション作成
        </Button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

// Discussion データ
const discussion = ref({
  title: '',
  content: '',
  categoryId: ''
})

// ドラフト保存用のID
const draftId = ref(null)

// カテゴリオプション（デモデータ）
const categoryOptions = [
  { value: '1', label: '全般' },
  { value: '2', label: '開発環境' },
  { value: '3', label: 'Q&A' },
  { value: '4', label: 'アイデア' },
  { value: '5', label: 'お知らせ' }
]

// プレースホルダー
const bodyPlaceholder = `## 概要
ディスカッションの目的や背景を簡潔に説明してください。

## 詳細
さらに詳しい情報や質問内容を記載してください。

## 検討したい点
議論したいポイントや質問を具体的に記載してください。

## 参考リンク
関連する資料やリンクを記載してください。`

// フォームバリデーション
const isFormValid = computed(() => {
  return discussion.value.title.trim() && 
         discussion.value.categoryId && 
         discussion.value.content.trim()
})

// ディスカッション作成
async function submitDiscussion() {
  if (!isFormValid.value) return
  
  try {
    // APIに送信（実装予定）
    console.log('ディスカッション作成:', discussion.value)
    
    // 成功したら下書きを削除
    if (draftId.value) {
      await deleteDraft(draftId.value)
    }
    
    // ディスカッション一覧ページに遷移
    router.push('/discussions')
  } catch (error) {
    console.error('ディスカッション作成中にエラーが発生しました:', error)
    // エラーハンドリング
  }
}

// 下書き保存
async function saveDraft() {
  if (!discussion.value.title.trim() && !discussion.value.content.trim()) {
    // タイトルと内容が両方空の場合は保存しない
    return
  }
  
  try {
    // APIに下書き保存（実装予定）
    console.log('下書き保存:', discussion.value)
    
    // 下書きIDを保存
    draftId.value = 'draft-' + Date.now()
    
    // 通知表示
    alert('下書きを保存しました')
  } catch (error) {
    console.error('下書き保存中にエラーが発生しました:', error)
    // エラーハンドリング
  }
}

// 下書き削除
async function deleteDraft(id) {
  try {
    // APIから下書き削除（実装予定）
    console.log('下書き削除:', id)
  } catch (error) {
    console.error('下書き削除中にエラーが発生しました:', error)
  }
}
</script>

<style scoped>
.discussion-create {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.discussion-create-header {
  margin-bottom: 2rem;
}

.discussion-create-title {
  font-size: 1.75rem;
  font-weight: 600;
}

.discussion-form {
  display: grid;
  grid-template-columns: 2fr 1fr;
  grid-template-areas: 
    "main sidebar"
    "actions actions";
  gap: 2rem;
}

.form-main {
  grid-area: main;
}

.form-sidebar {
  grid-area: sidebar;
}

.form-actions {
  grid-area: actions;
  display: flex;
  justify-content: flex-end;
  margin-top: 1.5rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--color-border-primary);
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-label {
  display: block;
  font-weight: 500;
  margin-bottom: 0.5rem;
}

.discussion-title-input {
  width: 100%;
  font-size: 1.125rem;
}

.sidebar-section {
  margin-bottom: 1.5rem;
}

.sidebar-heading {
  font-size: 1rem;
  font-weight: 600;
  margin-top: 0;
  margin-bottom: 1rem;
}

.tips-list {
  margin: 0;
  padding-left: 1.5rem;
}

.tips-list li {
  margin-bottom: 0.5rem;
  font-size: 0.875rem;
}

.markdown-tips {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.markdown-tip {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 0.875rem;
}

.markdown-tip code {
  padding: 0.25rem 0.5rem;
  background-color: var(--color-bg-secondary);
  border-radius: 0.25rem;
  font-family: monospace;
}

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .discussion-form {
    grid-template-columns: 1fr;
    grid-template-areas: 
      "sidebar"
      "main"
      "actions";
  }
}
</style>
