<template>
  <div 
    class="gh-label" 
    :class="[`gh-label--${size}`, { 'gh-label--clickable': clickable }]"
    :style="computedStyle"
    @click="handleClick"
  >
    <slot>{{ text }}</slot>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps({
  text: {
    type: String,
    default: ''
  },
  color: {
    type: String,
    default: '#0366d6'
  },
  size: {
    type: String,
    default: 'md',
    validator: (value: string) => ['sm', 'md', 'lg'].includes(value)
  },
  clickable: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['click']);

// ラベルの背景色と文字色を計算
const computedStyle = computed(() => {
  // カラーを16進数形式に正規化（#FFFFFFまたは#FFF形式）
  let color = props.color;
  if (!color.startsWith('#')) {
    color = `#${color}`;
  }
  
  // 色の明るさを計算（0-255）
  let r, g, b;
  if (color.length === 4) {
    // #RGB形式
    r = parseInt(color[1] + color[1], 16);
    g = parseInt(color[2] + color[2], 16);
    b = parseInt(color[3] + color[3], 16);
  } else {
    // #RRGGBB形式
    r = parseInt(color.substring(1, 3), 16);
    g = parseInt(color.substring(3, 5), 16);
    b = parseInt(color.substring(5, 7), 16);
  }
  
  // YIQ公式で明るさを計算
  const brightness = (r * 299 + g * 587 + b * 114) / 1000;
  
  // 明るさに基づいて文字色を決定（明るい背景には暗い文字、暗い背景には明るい文字）
  const textColor = brightness > 128 ? '#000000' : '#ffffff';
  
  // 少し薄い背景色（70%の不透明度）
  const alphaColor = `${color}B3`; // B3 is ~70% opacity in hex
  
  return {
    backgroundColor: alphaColor,
    color: textColor,
    borderColor: color
  };
});

// クリックイベントハンドラ
function handleClick(event: MouseEvent) {
  if (props.clickable) {
    emit('click', event);
  }
}
</script>

<style scoped>
.gh-label {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 500;
  border-radius: var(--border-radius-full);
  border: 1px solid transparent;
  white-space: nowrap;
  line-height: 1;
}

.gh-label--clickable {
  cursor: pointer;
  transition: var(--transition-base);
}

.gh-label--clickable:hover {
  filter: brightness(95%);
}

/* サイズバリエーション */
.gh-label--sm {
  padding: 2px 6px;
  font-size: 11px;
}

.gh-label--md {
  padding: 3px 8px;
  font-size: 12px;
}

.gh-label--lg {
  padding: 4px 10px;
  font-size: 13px;
}
</style>
