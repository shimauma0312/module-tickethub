<template>
  <div class="gh-input-wrapper" :class="{ 'gh-input-wrapper--error': error }">
    <label v-if="label" :for="id" class="gh-input-label">
      {{ label }}
      <span v-if="required" class="gh-input-required">*</span>
    </label>
    
    <div class="gh-input-container">
      <Icon 
        v-if="iconLeft" 
        :name="iconLeft" 
        class="gh-input-icon gh-input-icon--left" 
      />
      
      <input
        :id="id"
        class="gh-input"
        :class="[
          `gh-input--${size}`,
          { 'gh-input--with-left-icon': iconLeft },
          { 'gh-input--with-right-icon': iconRight || clearable }
        ]"
        :type="type"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :required="required"
        :autocomplete="autocomplete"
        @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
        @blur="$emit('blur', $event)"
        @focus="$emit('focus', $event)"
      />
      
      <button
        v-if="clearable && modelValue"
        type="button"
        class="gh-input-clear"
        @click="$emit('update:modelValue', '')"
      >
        <Icon name="mdi:close-circle" />
        <span class="sr-only">クリア</span>
      </button>
      
      <Icon 
        v-else-if="iconRight" 
        :name="iconRight" 
        class="gh-input-icon gh-input-icon--right" 
      />
    </div>
    
    <p v-if="error" class="gh-input-error">{{ error }}</p>
    <p v-else-if="hint" class="gh-input-hint">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  id: {
    type: String,
    default: () => `input-${Math.random().toString(36).substr(2, 9)}`
  },
  label: {
    type: String,
    default: ''
  },
  type: {
    type: String,
    default: 'text'
  },
  placeholder: {
    type: String,
    default: ''
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
  required: {
    type: Boolean,
    default: false
  },
  error: {
    type: String,
    default: ''
  },
  hint: {
    type: String,
    default: ''
  },
  iconLeft: {
    type: String,
    default: ''
  },
  iconRight: {
    type: String,
    default: ''
  },
  clearable: {
    type: Boolean,
    default: false
  },
  autocomplete: {
    type: String,
    default: 'off'
  }
});

defineEmits(['update:modelValue', 'blur', 'focus']);
</script>

<style scoped>
.gh-input-wrapper {
  margin-bottom: var(--spacing-md);
}

.gh-input-label {
  display: block;
  margin-bottom: var(--spacing-xs);
  font-weight: 500;
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
}

.gh-input-required {
  color: var(--color-danger);
  margin-left: var(--spacing-xxs);
}

.gh-input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.gh-input {
  width: 100%;
  border: 1px solid var(--color-border-primary);
  border-radius: var(--border-radius-md);
  background-color: var(--color-bg-canvas);
  color: var(--color-text-primary);
  transition: var(--transition-base);
  font-size: var(--font-size-sm);
  line-height: 1.5;
}

.gh-input:focus {
  outline: none;
  border-color: var(--color-accent-primary);
  box-shadow: 0 0 0 3px rgba(3, 102, 214, 0.3);
}

.gh-input:disabled {
  background-color: var(--color-bg-secondary);
  cursor: not-allowed;
  opacity: 0.7;
}

.gh-input::placeholder {
  color: var(--color-text-tertiary);
}

/* サイズバリエーション */
.gh-input--sm {
  padding: var(--spacing-xxs) var(--spacing-xs);
  font-size: var(--font-size-xs);
}

.gh-input--md {
  padding: var(--spacing-xs) var(--spacing-md);
}

.gh-input--lg {
  padding: var(--spacing-sm) var(--spacing-md);
  font-size: var(--font-size-md);
}

/* アイコン付きインプット */
.gh-input-icon {
  position: absolute;
  color: var(--color-text-tertiary);
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.gh-input-icon--left {
  left: var(--spacing-xs);
}

.gh-input-icon--right {
  right: var(--spacing-xs);
}

.gh-input--with-left-icon {
  padding-left: calc(var(--spacing-md) + 16px);
}

.gh-input--with-right-icon {
  padding-right: calc(var(--spacing-md) + 16px);
}

/* クリアボタン */
.gh-input-clear {
  position: absolute;
  right: var(--spacing-xs);
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  color: var(--color-text-tertiary);
  cursor: pointer;
  padding: 0;
  height: 16px;
  width: 16px;
}

.gh-input-clear:hover {
  color: var(--color-text-primary);
}

/* エラーとヒント */
.gh-input-error {
  margin-top: var(--spacing-xxs);
  color: var(--color-danger);
  font-size: var(--font-size-xs);
}

.gh-input-hint {
  margin-top: var(--spacing-xxs);
  color: var(--color-text-tertiary);
  font-size: var(--font-size-xs);
}

.gh-input-wrapper--error .gh-input {
  border-color: var(--color-danger);
}

.gh-input-wrapper--error .gh-input:focus {
  box-shadow: 0 0 0 3px rgba(215, 58, 73, 0.3);
}
</style>
