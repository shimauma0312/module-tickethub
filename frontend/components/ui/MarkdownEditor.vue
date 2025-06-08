<template>
  <div class="gh-markdown-editor" :class="{ 'gh-markdown-editor--fullscreen': isFullscreen }">
    <div class="gh-markdown-editor__header">
      <div class="gh-markdown-editor__tabs">
        <button 
          class="gh-markdown-editor__tab" 
          :class="{ 'gh-markdown-editor__tab--active': activeTab === 'write' }"
          @click="activeTab = 'write'"
        >
          <Icon name="mdi:pencil" />
          編集
        </button>
        <button 
          class="gh-markdown-editor__tab" 
          :class="{ 'gh-markdown-editor__tab--active': activeTab === 'preview' }"
          @click="activeTab = 'preview'"
        >
          <Icon name="mdi:eye" />
          プレビュー
        </button>
      </div>
      
      <div class="gh-markdown-editor__actions">
        <button class="gh-markdown-editor__action" @click="toggleFullscreen">
          <Icon :name="isFullscreen ? 'mdi:fullscreen-exit' : 'mdi:fullscreen'" />
          <span class="sr-only">{{ isFullscreen ? '全画面解除' : '全画面' }}</span>
        </button>
      </div>
    </div>
    
    <div class="gh-markdown-editor__content">
      <div v-show="activeTab === 'write'" class="gh-markdown-editor__write">
        <textarea
          ref="textareaRef"
          class="gh-markdown-editor__textarea"
          :value="modelValue"
          :placeholder="placeholder"
          @input="$emit('update:modelValue', ($event.target as HTMLTextAreaElement).value)"
        ></textarea>
        
        <div class="gh-markdown-editor__toolbar">
          <button class="gh-markdown-editor__tool" @click="insertMarkdown('**', '**', 'bold text')">
            <Icon name="mdi:format-bold" />
            <span class="sr-only">太字</span>
          </button>
          <button class="gh-markdown-editor__tool" @click="insertMarkdown('*', '*', 'italic text')">
            <Icon name="mdi:format-italic" />
            <span class="sr-only">斜体</span>
          </button>
          <button class="gh-markdown-editor__tool" @click="insertMarkdown('[', '](url)', 'link text')">
            <Icon name="mdi:link" />
            <span class="sr-only">リンク</span>
          </button>
          <button class="gh-markdown-editor__tool" @click="insertMarkdown('- ', '', 'List item')">
            <Icon name="mdi:format-list-bulleted" />
            <span class="sr-only">リスト</span>
          </button>
          <button class="gh-markdown-editor__tool" @click="insertMarkdown('1. ', '', 'List item')">
            <Icon name="mdi:format-list-numbered" />
            <span class="sr-only">番号付きリスト</span>
          </button>
          <button class="gh-markdown-editor__tool" @click="insertMarkdown('```\n', '\n```', 'code')">
            <Icon name="mdi:code-tags" />
            <span class="sr-only">コード</span>
          </button>
          <button class="gh-markdown-editor__tool" @click="insertMarkdown('> ', '', 'quote')">
            <Icon name="mdi:format-quote-close" />
            <span class="sr-only">引用</span>
          </button>
          <button class="gh-markdown-editor__tool" @click="insertMarkdown('### ', '', 'Heading')">
            <Icon name="mdi:format-header-3" />
            <span class="sr-only">見出し</span>
          </button>
        </div>
      </div>
      
      <div v-show="activeTab === 'preview'" class="gh-markdown-editor__preview">
        <div class="gh-markdown-content" v-html="renderedMarkdown"></div>
      </div>
    </div>
    
    <p v-if="hint" class="gh-markdown-editor__hint">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue';
import { marked } from 'marked';

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: 'Markdownを入力...'
  },
  hint: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['update:modelValue']);

const activeTab = ref('write');
const isFullscreen = ref(false);
const textareaRef = ref<HTMLTextAreaElement | null>(null);

const renderedMarkdown = computed(() => {
  return marked(props.modelValue || '', { 
    gfm: true,
    breaks: true,
    sanitize: true
  });
});

function toggleFullscreen() {
  isFullscreen.value = !isFullscreen.value;
  // フルスクリーン時にフォーカスを戻す
  if (isFullscreen.value && activeTab.value === 'write') {
    setTimeout(() => {
      textareaRef.value?.focus();
    }, 10);
  }
}

// Markdownを挿入する
function insertMarkdown(prefix: string, suffix: string, placeholder: string) {
  const textarea = textareaRef.value;
  if (!textarea) return;
  
  const start = textarea.selectionStart;
  const end = textarea.selectionEnd;
  const text = textarea.value;
  
  let selectedText = text.substring(start, end);
  
  // 選択されていない場合はプレースホルダーを使用
  if (selectedText.length === 0) {
    selectedText = placeholder;
  }
  
  const newText = text.substring(0, start) + prefix + selectedText + suffix + text.substring(end);
  
  emit('update:modelValue', newText);
  
  // カーソル位置を調整（選択されたテキストの後ろに配置）
  setTimeout(() => {
    if (textarea) {
      const newCursorPos = start + prefix.length + selectedText.length + suffix.length;
      textarea.focus();
      textarea.setSelectionRange(newCursorPos, newCursorPos);
    }
  }, 10);
}

// ESCキーでフルスクリーン解除
onMounted(() => {
  const handleKeyDown = (e: KeyboardEvent) => {
    if (e.key === 'Escape' && isFullscreen.value) {
      isFullscreen.value = false;
    }
  };
  
  window.addEventListener('keydown', handleKeyDown);
  
  return () => {
    window.removeEventListener('keydown', handleKeyDown);
  };
});
</script>

<style scoped>
.gh-markdown-editor {
  border: 1px solid var(--color-border-primary);
  border-radius: var(--border-radius-md);
  background-color: var(--color-bg-canvas);
  transition: var(--transition-base);
  overflow: hidden;
}

.gh-markdown-editor:focus-within {
  border-color: var(--color-accent-primary);
  box-shadow: 0 0 0 3px rgba(3, 102, 214, 0.3);
}

.gh-markdown-editor--fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1000;
  border-radius: 0;
  border: none;
}

.gh-markdown-editor__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-xs);
  border-bottom: 1px solid var(--color-border-primary);
  background-color: var(--color-bg-secondary);
}

.gh-markdown-editor__tabs {
  display: flex;
}

.gh-markdown-editor__tab {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-xxs);
  padding: var(--spacing-xs) var(--spacing-sm);
  background: none;
  border: none;
  border-radius: var(--border-radius-md);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: var(--transition-base);
}

.gh-markdown-editor__tab:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.gh-markdown-editor__tab--active {
  color: var(--color-accent-primary);
  font-weight: 500;
}

.gh-markdown-editor__actions {
  display: flex;
  align-items: center;
}

.gh-markdown-editor__action {
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  border-radius: var(--border-radius-md);
  padding: var(--spacing-xxs);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: var(--transition-base);
}

.gh-markdown-editor__action:hover {
  background-color: rgba(0, 0, 0, 0.05);
  color: var(--color-text-primary);
}

.gh-markdown-editor__content {
  position: relative;
  height: 300px;
}

.gh-markdown-editor--fullscreen .gh-markdown-editor__content {
  height: calc(100vh - 100px);
}

.gh-markdown-editor__write,
.gh-markdown-editor__preview {
  height: 100%;
  overflow-y: auto;
}

.gh-markdown-editor__textarea {
  width: 100%;
  height: calc(100% - 40px);
  border: none;
  resize: none;
  padding: var(--spacing-md);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
  font-size: var(--font-size-sm);
  line-height: 1.5;
  color: var(--color-text-primary);
  background: transparent;
}

.gh-markdown-editor__textarea:focus {
  outline: none;
}

.gh-markdown-editor__toolbar {
  display: flex;
  padding: var(--spacing-xs);
  border-top: 1px solid var(--color-border-primary);
  background-color: var(--color-bg-secondary);
  overflow-x: auto;
}

.gh-markdown-editor__tool {
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  border-radius: var(--border-radius-md);
  padding: var(--spacing-xxs) var(--spacing-xs);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: var(--transition-base);
}

.gh-markdown-editor__tool:hover {
  background-color: rgba(0, 0, 0, 0.05);
  color: var(--color-text-primary);
}

.gh-markdown-editor__preview {
  padding: var(--spacing-md);
}

.gh-markdown-editor__hint {
  padding: var(--spacing-xs) var(--spacing-md);
  font-size: var(--font-size-xs);
  color: var(--color-text-tertiary);
  border-top: 1px solid var(--color-border-primary);
}

/* Markdownのスタイル（グローバルなので、scoped外に書く必要があります） */
:deep(.gh-markdown-content) {
  font-size: var(--font-size-sm);
  line-height: 1.6;
  word-wrap: break-word;
  color: var(--color-text-primary);
}

:deep(.gh-markdown-content h1),
:deep(.gh-markdown-content h2),
:deep(.gh-markdown-content h3),
:deep(.gh-markdown-content h4),
:deep(.gh-markdown-content h5),
:deep(.gh-markdown-content h6) {
  margin-top: 24px;
  margin-bottom: 16px;
  font-weight: 600;
  line-height: 1.25;
}

:deep(.gh-markdown-content h1) {
  font-size: 2em;
  border-bottom: 1px solid var(--color-border-primary);
  padding-bottom: 0.3em;
}

:deep(.gh-markdown-content h2) {
  font-size: 1.5em;
  border-bottom: 1px solid var(--color-border-primary);
  padding-bottom: 0.3em;
}

:deep(.gh-markdown-content h3) {
  font-size: 1.25em;
}

:deep(.gh-markdown-content p) {
  margin-top: 0;
  margin-bottom: 16px;
}

:deep(.gh-markdown-content ul),
:deep(.gh-markdown-content ol) {
  margin-top: 0;
  margin-bottom: 16px;
  padding-left: 2em;
}

:deep(.gh-markdown-content li) {
  margin-bottom: 0.25em;
}

:deep(.gh-markdown-content a) {
  color: var(--color-accent-primary);
  text-decoration: none;
}

:deep(.gh-markdown-content a:hover) {
  text-decoration: underline;
}

:deep(.gh-markdown-content code) {
  padding: 0.2em 0.4em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(27, 31, 35, 0.05);
  border-radius: 3px;
  font-family: SFMono-Regular, Consolas, Liberation Mono, Menlo, monospace;
}

:deep(.gh-markdown-content pre) {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: var(--color-bg-secondary);
  border-radius: var(--border-radius-md);
  margin-top: 0;
  margin-bottom: 16px;
}

:deep(.gh-markdown-content pre code) {
  padding: 0;
  margin: 0;
  font-size: 100%;
  background-color: transparent;
  border: 0;
  display: inline;
  overflow: visible;
  line-height: inherit;
  word-wrap: normal;
}

:deep(.gh-markdown-content blockquote) {
  padding: 0 1em;
  color: var(--color-text-tertiary);
  border-left: 0.25em solid var(--color-border-primary);
  margin: 0 0 16px 0;
}

:deep(.gh-markdown-content img) {
  max-width: 100%;
  box-sizing: border-box;
  border-style: none;
}

:deep(.gh-markdown-content hr) {
  height: 0.25em;
  padding: 0;
  margin: 24px 0;
  background-color: var(--color-border-primary);
  border: 0;
}

:deep(.gh-markdown-content table) {
  border-spacing: 0;
  border-collapse: collapse;
  display: block;
  width: 100%;
  overflow: auto;
  margin-top: 0;
  margin-bottom: 16px;
}

:deep(.gh-markdown-content table tr) {
  background-color: var(--color-bg-canvas);
  border-top: 1px solid var(--color-border-primary);
}

:deep(.gh-markdown-content table tr:nth-child(2n)) {
  background-color: var(--color-bg-secondary);
}

:deep(.gh-markdown-content table th),
:deep(.gh-markdown-content table td) {
  padding: 6px 13px;
  border: 1px solid var(--color-border-primary);
}

:deep(.gh-markdown-content table th) {
  font-weight: 600;
}
</style>
