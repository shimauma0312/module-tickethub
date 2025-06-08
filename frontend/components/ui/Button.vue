<template>
  <button 
    :class="[
      'gh-button', 
      `gh-button--${variant}`,
      `gh-button--${size}`,
      { 'gh-button--full': full },
      { 'gh-button--loading': loading }
    ]" 
    :disabled="disabled || loading"
    :type="type"
    @click="$emit('click', $event)"
  >
    <Icon v-if="iconLeft && !loading" :name="iconLeft" class="gh-button__icon gh-button__icon--left" />
    <span v-if="loading" class="gh-button__loader"></span>
    <span class="gh-button__content">
      <slot></slot>
    </span>
    <Icon v-if="iconRight && !loading" :name="iconRight" class="gh-button__icon gh-button__icon--right" />
  </button>
</template>

<script setup lang="ts">
defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (value: string) => ['primary', 'secondary', 'danger', 'outline', 'text'].includes(value)
  },
  size: {
    type: String,
    default: 'md',
    validator: (value: string) => ['sm', 'md', 'lg'].includes(value)
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  full: {
    type: Boolean,
    default: false
  },
  type: {
    type: String,
    default: 'button'
  },
  iconLeft: {
    type: String,
    default: ''
  },
  iconRight: {
    type: String,
    default: ''
  }
});

defineEmits(['click']);
</script>

<style scoped>
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
  position: relative;
  overflow: hidden;
  border: 1px solid transparent;
}

.gh-button:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(3, 102, 214, 0.3);
}

.gh-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* サイズバリエーション */
.gh-button--sm {
  padding: var(--spacing-xxs) var(--spacing-xs);
  font-size: var(--font-size-xs);
}

.gh-button--md {
  padding: var(--spacing-xs) var(--spacing-md);
}

.gh-button--lg {
  padding: var(--spacing-sm) var(--spacing-lg);
  font-size: var(--font-size-md);
}

/* バリアントスタイル */
.gh-button--primary {
  background-color: var(--color-accent-primary);
  color: white;
  border-color: var(--color-accent-primary);
}

.gh-button--primary:hover:not(:disabled) {
  background-color: #0550ae;
  border-color: #0550ae;
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

.gh-button--danger {
  background-color: var(--color-danger);
  color: white;
  border-color: var(--color-danger);
}

.gh-button--danger:hover:not(:disabled) {
  background-color: #cb2431;
  border-color: #cb2431;
}

.gh-button--outline {
  background-color: transparent;
  color: var(--color-text-primary);
  border-color: var(--color-border-primary);
}

.gh-button--outline:hover:not(:disabled) {
  background-color: var(--color-bg-secondary);
}

.gh-button--text {
  background-color: transparent;
  color: var(--color-accent-primary);
  border: none;
  padding-left: var(--spacing-xs);
  padding-right: var(--spacing-xs);
}

.gh-button--text:hover:not(:disabled) {
  color: #0550ae;
  background-color: rgba(3, 102, 214, 0.1);
}

/* Full幅 */
.gh-button--full {
  width: 100%;
}

/* アイコン */
.gh-button__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
}

.gh-button__icon--left {
  margin-right: var(--spacing-xxs);
}

.gh-button__icon--right {
  margin-left: var(--spacing-xxs);
}

/* ローディング状態 */
.gh-button--loading .gh-button__content {
  visibility: hidden;
}

.gh-button__loader {
  position: absolute;
  width: 16px;
  height: 16px;
  border: 2px solid currentColor;
  border-radius: 50%;
  border-right-color: transparent;
  animation: gh-button-spin 0.8s linear infinite;
}

@keyframes gh-button-spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
</style>
