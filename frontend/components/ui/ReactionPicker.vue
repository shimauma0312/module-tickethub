<template>
  <div class="gh-reaction-picker" :class="{ 'gh-reaction-picker--open': isOpen }">
    <button
      v-if="!isOpen"
      class="gh-reaction-picker__trigger"
      @click="isOpen = true"
      aria-label="リアクションを追加"
    >
      <slot name="trigger">
        <Icon name="mdi:emoticon-outline" />
        <span class="gh-reaction-picker__trigger-text">リアクション</span>
      </slot>
    </button>
    
    <div v-if="isOpen" class="gh-reaction-picker__container">
      <div class="gh-reaction-picker__header">
        <span class="gh-reaction-picker__title">リアクション</span>
        <button
          class="gh-reaction-picker__close"
          @click="isOpen = false"
          aria-label="閉じる"
        >
          <Icon name="mdi:close" />
        </button>
      </div>
      
      <div class="gh-reaction-picker__content">
        <div class="gh-reaction-picker__section">
          <h3 class="gh-reaction-picker__section-title">よく使われる</h3>
          <div class="gh-reaction-picker__emoji-grid">
            <button
              v-for="emoji in popularEmojis"
              :key="emoji.id"
              class="gh-reaction-picker__emoji"
              @click="selectEmoji(emoji)"
              :aria-label="emoji.name"
            >
              {{ emoji.character }}
            </button>
          </div>
        </div>
        
        <div class="gh-reaction-picker__section">
          <h3 class="gh-reaction-picker__section-title">すべて</h3>
          <div class="gh-reaction-picker__tabs">
            <button
              v-for="category in categories"
              :key="category.id"
              class="gh-reaction-picker__tab"
              :class="{ 'gh-reaction-picker__tab--active': activeCategory === category.id }"
              @click="activeCategory = category.id"
            >
              <span class="gh-reaction-picker__tab-icon">{{ category.icon }}</span>
              <span class="sr-only">{{ category.name }}</span>
            </button>
          </div>
          
          <div class="gh-reaction-picker__emoji-grid gh-reaction-picker__emoji-grid--scrollable">
            <button
              v-for="emoji in filteredEmojis"
              :key="emoji.id"
              class="gh-reaction-picker__emoji"
              @click="selectEmoji(emoji)"
              :aria-label="emoji.name"
            >
              {{ emoji.character }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';

interface Emoji {
  id: string;
  character: string;
  name: string;
  category: string;
}

interface Category {
  id: string;
  name: string;
  icon: string;
}

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  }
});

const emit = defineEmits(['update:modelValue', 'select']);

const isOpen = ref(false);
const activeCategory = ref('smileys');

// カテゴリ定義
const categories: Category[] = [
  { id: 'smileys', name: '笑顔と感情', icon: '😀' },
  { id: 'people', name: '人と体', icon: '👋' },
  { id: 'nature', name: '動物と自然', icon: '🐱' },
  { id: 'food', name: '食べ物と飲み物', icon: '🍔' },
  { id: 'activities', name: 'アクティビティ', icon: '⚽' },
  { id: 'travel', name: '旅行と場所', icon: '🚗' },
  { id: 'objects', name: 'オブジェクト', icon: '💡' },
  { id: 'symbols', name: 'シンボル', icon: '❤️' },
  { id: 'flags', name: '旗', icon: '🏁' }
];

// 絵文字データ（実際の実装では完全なリストを使用）
const emojis: Emoji[] = [
  // スマイリーと感情
  { id: 'smile', character: '😀', name: 'にっこり', category: 'smileys' },
  { id: 'laugh', character: '😂', name: '笑い泣き', category: 'smileys' },
  { id: 'blush', character: '😊', name: '照れ笑い', category: 'smileys' },
  { id: 'heart_eyes', character: '😍', name: 'ハートの目', category: 'smileys' },
  { id: 'thinking', character: '🤔', name: '考え中', category: 'smileys' },
  { id: 'confused', character: '😕', name: '困惑', category: 'smileys' },
  
  // 人と体
  { id: 'thumbs_up', character: '👍', name: 'いいね', category: 'people' },
  { id: 'thumbs_down', character: '👎', name: 'だめ', category: 'people' },
  { id: 'clap', character: '👏', name: '拍手', category: 'people' },
  { id: 'muscle', character: '💪', name: '筋肉', category: 'people' },
  { id: 'wave', character: '👋', name: '手を振る', category: 'people' },
  
  // 動物と自然
  { id: 'dog', character: '🐶', name: '犬', category: 'nature' },
  { id: 'cat', character: '🐱', name: '猫', category: 'nature' },
  { id: 'unicorn', character: '🦄', name: 'ユニコーン', category: 'nature' },
  
  // 食べ物と飲み物
  { id: 'pizza', character: '🍕', name: 'ピザ', category: 'food' },
  { id: 'burger', character: '🍔', name: 'ハンバーガー', category: 'food' },
  { id: 'coffee', character: '☕', name: 'コーヒー', category: 'food' },
  
  // アクティビティ
  { id: 'party', character: '🎉', name: 'パーティー', category: 'activities' },
  { id: 'trophy', character: '🏆', name: 'トロフィー', category: 'activities' },
  
  // オブジェクト
  { id: 'bulb', character: '💡', name: '電球', category: 'objects' },
  { id: 'rocket', character: '🚀', name: 'ロケット', category: 'objects' },
  { id: 'eyes', character: '👀', name: '目', category: 'objects' },
  
  // シンボル
  { id: 'heart', character: '❤️', name: 'ハート', category: 'symbols' },
  { id: 'fire', character: '🔥', name: '炎', category: 'symbols' },
  { id: 'check', character: '✅', name: 'チェック', category: 'symbols' },
  { id: 'warning', character: '⚠️', name: '警告', category: 'symbols' }
];

// よく使われる絵文字
const popularEmojis = computed(() => {
  return [
    emojis.find(e => e.id === 'thumbs_up'),
    emojis.find(e => e.id === 'heart'),
    emojis.find(e => e.id === 'laugh'),
    emojis.find(e => e.id === 'fire'),
    emojis.find(e => e.id === 'clap'),
    emojis.find(e => e.id === 'thinking'),
    emojis.find(e => e.id === 'rocket'),
    emojis.find(e => e.id === 'eyes')
  ].filter(Boolean) as Emoji[];
});

// 現在のカテゴリに基づいてフィルタリングした絵文字
const filteredEmojis = computed(() => {
  return emojis.filter(emoji => emoji.category === activeCategory.value);
});

// 絵文字を選択
function selectEmoji(emoji: Emoji) {
  emit('update:modelValue', emoji.character);
  emit('select', emoji);
  isOpen.value = false;
}

// クリックアウトで閉じるための監視を設定
function onClickOutside(event: MouseEvent) {
  // コンポーネント外のクリックで閉じる
  if (isOpen.value) {
    const target = event.target as HTMLElement;
    const picker = document.querySelector('.gh-reaction-picker');
    if (picker && !picker.contains(target)) {
      isOpen.value = false;
    }
  }
}

// コンポーネントマウント時にイベントリスナーを追加
onMounted(() => {
  document.addEventListener('click', onClickOutside);
});

// コンポーネントアンマウント時にイベントリスナーを削除
onUnmounted(() => {
  document.removeEventListener('click', onClickOutside);
});
</script>

<style scoped>
.gh-reaction-picker {
  position: relative;
  display: inline-block;
}

.gh-reaction-picker__trigger {
  display: inline-flex;
  align-items: center;
  gap: var(--spacing-xxs);
  padding: var(--spacing-xs) var(--spacing-sm);
  background: none;
  border: 1px solid var(--color-border-primary);
  border-radius: var(--border-radius-md);
  color: var(--color-text-secondary);
  font-size: var(--font-size-sm);
  cursor: pointer;
  transition: var(--transition-base);
}

.gh-reaction-picker__trigger:hover {
  background-color: var(--color-bg-secondary);
  color: var(--color-text-primary);
}

.gh-reaction-picker__trigger-text {
  margin-left: var(--spacing-xxs);
}

.gh-reaction-picker__container {
  position: absolute;
  bottom: calc(100% + 8px);
  left: 0;
  width: 320px;
  background-color: var(--color-bg-canvas);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-lg);
  z-index: 100;
}

.gh-reaction-picker__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-xs) var(--spacing-sm);
  border-bottom: 1px solid var(--color-border-primary);
}

.gh-reaction-picker__title {
  font-weight: 600;
  font-size: var(--font-size-sm);
}

.gh-reaction-picker__close {
  background: none;
  border: none;
  color: var(--color-text-tertiary);
  cursor: pointer;
  padding: var(--spacing-xxs);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--border-radius-md);
}

.gh-reaction-picker__close:hover {
  background-color: var(--color-bg-secondary);
  color: var(--color-text-primary);
}

.gh-reaction-picker__content {
  padding: var(--spacing-sm);
}

.gh-reaction-picker__section {
  margin-bottom: var(--spacing-md);
}

.gh-reaction-picker__section:last-child {
  margin-bottom: 0;
}

.gh-reaction-picker__section-title {
  font-size: var(--font-size-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  margin-bottom: var(--spacing-xs);
}

.gh-reaction-picker__emoji-grid {
  display: grid;
  grid-template-columns: repeat(8, 1fr);
  gap: var(--spacing-xxs);
}

.gh-reaction-picker__emoji-grid--scrollable {
  max-height: 200px;
  overflow-y: auto;
  padding-right: var(--spacing-xs);
}

.gh-reaction-picker__emoji {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  font-size: 20px;
  background: none;
  border: none;
  border-radius: var(--border-radius-md);
  cursor: pointer;
  transition: var(--transition-base);
}

.gh-reaction-picker__emoji:hover {
  background-color: var(--color-bg-secondary);
}

.gh-reaction-picker__tabs {
  display: flex;
  border-bottom: 1px solid var(--color-border-primary);
  margin-bottom: var(--spacing-xs);
  overflow-x: auto;
  padding-bottom: var(--spacing-xxs);
}

.gh-reaction-picker__tab {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-xxs);
  background: none;
  border: none;
  border-radius: var(--border-radius-md);
  cursor: pointer;
  transition: var(--transition-base);
  margin-right: var(--spacing-xxs);
}

.gh-reaction-picker__tab:hover {
  background-color: var(--color-bg-secondary);
}

.gh-reaction-picker__tab--active {
  background-color: var(--color-bg-secondary);
  border-bottom: 2px solid var(--color-accent-primary);
}

.gh-reaction-picker__tab-icon {
  font-size: 18px;
}
</style>
