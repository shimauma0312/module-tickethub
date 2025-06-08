<template>
  <div class="gh-comment" :class="{ 'gh-comment--highlighted': isHighlighted }">
    <div class="gh-comment__header">
      <div class="gh-comment__author">
        <div class="gh-comment__avatar">
          <img 
            v-if="author.avatarUrl" 
            :src="author.avatarUrl" 
            :alt="`${author.name}のアバター`"
          />
          <div v-else class="gh-comment__avatar-placeholder">
            {{ author.name.charAt(0).toUpperCase() }}
          </div>
        </div>
        <div class="gh-comment__author-info">
          <div class="gh-comment__author-name">{{ author.name }}</div>
          <div class="gh-comment__timestamp">
            {{ formatDate(createdAt) }}
            <span v-if="isEdited"> · 編集済み</span>
          </div>
        </div>
      </div>
      
      <div class="gh-comment__actions">
        <button 
          v-if="canEdit" 
          class="gh-comment__action-btn"
          @click="$emit('edit')"
        >
          <Icon name="mdi:pencil" />
          <span class="sr-only">編集</span>
        </button>
        <button 
          v-if="canDelete" 
          class="gh-comment__action-btn gh-comment__action-btn--danger"
          @click="$emit('delete')"
        >
          <Icon name="mdi:delete" />
          <span class="sr-only">削除</span>
        </button>
      </div>
    </div>
    
    <div class="gh-comment__body">
      <div v-if="isEditing">
        <MarkdownEditor
          v-model="editContent"
          :placeholder="'コメントを編集...'"
        />
        <div class="gh-comment__edit-actions">
          <button 
            class="gh-button gh-button--sm gh-button--secondary"
            @click="saveEdit"
          >
            保存
          </button>
          <button 
            class="gh-button gh-button--sm gh-button--outline"
            @click="cancelEdit"
          >
            キャンセル
          </button>
        </div>
      </div>
      <div v-else class="gh-markdown-content" v-html="renderedContent"></div>
    </div>
    
    <div v-if="showReactions" class="gh-comment__reactions">
      <div 
        v-for="(reaction, idx) in reactions" 
        :key="idx"
        class="gh-comment__reaction"
        :class="{ 'gh-comment__reaction--active': reaction.hasReacted }"
        @click="$emit('react', reaction.type)"
      >
        <span class="gh-comment__reaction-emoji">{{ reaction.emoji }}</span>
        <span class="gh-comment__reaction-count">{{ reaction.count }}</span>
      </div>
      
      <button 
        v-if="showAddReaction"
        class="gh-comment__add-reaction"
        @click="$emit('show-reaction-picker')"
      >
        <Icon name="mdi:emoticon-outline" />
        <span class="sr-only">リアクションを追加</span>
      </button>
    </div>
    
    <div v-if="showReplyForm" class="gh-comment__reply-form">
      <MarkdownEditor
        v-model="replyContent"
        :placeholder="'返信を入力...'"
      />
      <div class="gh-comment__reply-actions">
        <button 
          class="gh-button gh-button--sm gh-button--secondary"
          :disabled="!replyContent.trim()"
          @click="submitReply"
        >
          返信
        </button>
        <button 
          class="gh-button gh-button--sm gh-button--outline"
          @click="$emit('cancel-reply')"
        >
          キャンセル
        </button>
      </div>
    </div>
    
    <div v-if="replies && replies.length > 0" class="gh-comment__replies">
      <div 
        v-for="(reply, idx) in replies" 
        :key="reply.id || idx"
        class="gh-comment__reply"
      >
        <slot name="reply" :reply="reply"></slot>
      </div>
    </div>
    
    <div v-if="!showReplyForm && showReplyButton" class="gh-comment__reply-button-container">
      <button 
        class="gh-button gh-button--sm gh-button--outline"
        @click="$emit('reply')"
      >
        <Icon name="mdi:reply" />
        返信
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { marked } from 'marked';

// コメント投稿者の情報
interface Author {
  id: string | number;
  name: string;
  avatarUrl?: string;
}

// リアクション情報
interface Reaction {
  type: string;
  emoji: string;
  count: number;
  hasReacted: boolean;
}

// リプライ情報（必要に応じて型を拡張）
interface Reply {
  id: string | number;
  [key: string]: any;
}

const props = defineProps({
  id: {
    type: [String, Number],
    required: true
  },
  author: {
    type: Object as () => Author,
    required: true
  },
  content: {
    type: String,
    default: ''
  },
  createdAt: {
    type: [Date, String],
    required: true
  },
  updatedAt: {
    type: [Date, String],
    default: null
  },
  isEdited: {
    type: Boolean,
    default: false
  },
  isHighlighted: {
    type: Boolean,
    default: false
  },
  canEdit: {
    type: Boolean,
    default: false
  },
  canDelete: {
    type: Boolean,
    default: false
  },
  showReactions: {
    type: Boolean,
    default: true
  },
  reactions: {
    type: Array as () => Reaction[],
    default: () => []
  },
  showAddReaction: {
    type: Boolean,
    default: true
  },
  showReplyButton: {
    type: Boolean,
    default: true
  },
  showReplyForm: {
    type: Boolean,
    default: false
  },
  replies: {
    type: Array as () => Reply[],
    default: () => []
  }
});

const emit = defineEmits([
  'edit', 
  'delete', 
  'update', 
  'reply', 
  'cancel-reply', 
  'submit-reply',
  'react',
  'show-reaction-picker'
]);

// 編集モードのステート
const isEditing = ref(false);
const editContent = ref(props.content);

// 返信フォームのステート
const replyContent = ref('');

// Markdownをレンダリング
const renderedContent = computed(() => {
  return marked(props.content || '', { 
    gfm: true,
    breaks: true,
    sanitize: true
  });
});

// 日付フォーマット
function formatDate(date: Date | string): string {
  if (!date) return '';
  
  const d = new Date(date);
  
  // 今日の日付
  const today = new Date();
  const yesterday = new Date(today);
  yesterday.setDate(yesterday.getDate() - 1);
  
  // 日付が今日または昨日の場合は特別な表示
  if (d.toDateString() === today.toDateString()) {
    return `今日 ${d.getHours()}:${String(d.getMinutes()).padStart(2, '0')}`;
  } else if (d.toDateString() === yesterday.toDateString()) {
    return `昨日 ${d.getHours()}:${String(d.getMinutes()).padStart(2, '0')}`;
  }
  
  // それ以外は標準的な形式
  return `${d.getFullYear()}年${d.getMonth() + 1}月${d.getDate()}日`;
}

// 編集を開始
function startEdit() {
  isEditing.value = true;
  editContent.value = props.content;
}

// 編集を保存
function saveEdit() {
  emit('update', {
    id: props.id,
    content: editContent.value
  });
  isEditing.value = false;
}

// 編集をキャンセル
function cancelEdit() {
  isEditing.value = false;
  editContent.value = props.content;
}

// 返信を送信
function submitReply() {
  if (!replyContent.value.trim()) return;
  
  emit('submit-reply', {
    parentId: props.id,
    content: replyContent.value
  });
  
  replyContent.value = '';
}
</script>

<style scoped>
.gh-comment {
  border: 1px solid var(--color-border-primary);
  border-radius: var(--border-radius-md);
  background-color: var(--color-bg-canvas);
  margin-bottom: var(--spacing-md);
  transition: var(--transition-base);
}

.gh-comment--highlighted {
  box-shadow: 0 0 0 1px var(--color-accent-primary);
  border-color: var(--color-accent-primary);
}

.gh-comment__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm);
  background-color: var(--color-bg-secondary);
  border-bottom: 1px solid var(--color-border-primary);
  border-top-left-radius: var(--border-radius-md);
  border-top-right-radius: var(--border-radius-md);
}

.gh-comment__author {
  display: flex;
  align-items: center;
}

.gh-comment__avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  overflow: hidden;
  margin-right: var(--spacing-xs);
}

.gh-comment__avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.gh-comment__avatar-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--color-accent-primary);
  color: white;
  font-weight: 600;
  font-size: var(--font-size-md);
}

.gh-comment__author-info {
  display: flex;
  flex-direction: column;
}

.gh-comment__author-name {
  font-weight: 600;
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
}

.gh-comment__timestamp {
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
}

.gh-comment__actions {
  display: flex;
  gap: var(--spacing-xxs);
}

.gh-comment__action-btn {
  background: none;
  border: none;
  color: var(--color-text-tertiary);
  padding: var(--spacing-xxs);
  border-radius: var(--border-radius-md);
  cursor: pointer;
  transition: var(--transition-base);
  display: flex;
  align-items: center;
  justify-content: center;
}

.gh-comment__action-btn:hover {
  background-color: rgba(0, 0, 0, 0.05);
  color: var(--color-text-primary);
}

.gh-comment__action-btn--danger:hover {
  color: var(--color-danger);
}

.gh-comment__body {
  padding: var(--spacing-md);
}

.gh-comment__edit-actions,
.gh-comment__reply-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--spacing-xs);
  margin-top: var(--spacing-sm);
}

.gh-comment__reactions {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-xs);
  padding: 0 var(--spacing-md) var(--spacing-md);
}

.gh-comment__reaction {
  display: inline-flex;
  align-items: center;
  padding: var(--spacing-xxs) var(--spacing-xs);
  border-radius: var(--border-radius-full);
  border: 1px solid var(--color-border-primary);
  background-color: var(--color-bg-secondary);
  cursor: pointer;
  transition: var(--transition-base);
}

.gh-comment__reaction:hover {
  background-color: rgba(3, 102, 214, 0.1);
}

.gh-comment__reaction--active {
  background-color: rgba(3, 102, 214, 0.1);
  border-color: var(--color-accent-primary);
}

.gh-comment__reaction-emoji {
  margin-right: var(--spacing-xxs);
}

.gh-comment__reaction-count {
  font-size: var(--font-size-xs);
  font-weight: 500;
}

.gh-comment__add-reaction {
  background: none;
  border: none;
  color: var(--color-text-tertiary);
  padding: var(--spacing-xxs);
  border-radius: var(--border-radius-md);
  cursor: pointer;
  transition: var(--transition-base);
  display: flex;
  align-items: center;
  justify-content: center;
}

.gh-comment__add-reaction:hover {
  background-color: rgba(0, 0, 0, 0.05);
  color: var(--color-text-primary);
}

.gh-comment__reply-form {
  padding: var(--spacing-md);
  border-top: 1px solid var(--color-border-primary);
}

.gh-comment__replies {
  border-top: 1px solid var(--color-border-primary);
  padding: var(--spacing-md);
  background-color: var(--color-bg-secondary);
}

.gh-comment__reply {
  margin-bottom: var(--spacing-sm);
}

.gh-comment__reply:last-child {
  margin-bottom: 0;
}

.gh-comment__reply-button-container {
  padding: 0 var(--spacing-md) var(--spacing-md);
}

/* GitHub風ボタン（コンポーネント化されていれば不要） */
.gh-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: var(--font-size-sm);
  font-weight: 500;
  border-radius: var(--border-radius-md);
  transition: var(--transition-base);
  cursor: pointer;
  white-space: nowrap;
  border: 1px solid transparent;
}

.gh-button--sm {
  padding: var(--spacing-xxs) var(--spacing-xs);
  font-size: var(--font-size-xs);
}

.gh-button--secondary {
  background-color: var(--color-accent-secondary);
  color: white;
  border-color: var(--color-accent-secondary);
}

.gh-button--secondary:hover:not(:disabled) {
  background-color: #2c974b;
  border-color: #2c974b;
}

.gh-button--outline {
  background-color: transparent;
  color: var(--color-text-primary);
  border-color: var(--color-border-primary);
}

.gh-button--outline:hover:not(:disabled) {
  background-color: var(--color-bg-secondary);
}

.gh-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Markdownスタイル（MarkdownEditorコンポーネントにも同様のスタイルがあります） */
:deep(.gh-markdown-content) {
  font-size: var(--font-size-sm);
  line-height: 1.6;
  word-wrap: break-word;
  color: var(--color-text-primary);
}

/* 他のMarkdownスタイルはMarkdownEditorコンポーネントから共通化されるべきですが、ここでは省略 */

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .gh-comment__header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .gh-comment__actions {
    margin-top: var(--spacing-xs);
    align-self: flex-end;
  }
}
</style>
