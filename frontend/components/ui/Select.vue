<template>
  <div class="gh-select-wrapper" :class="{ 'gh-select-wrapper--error': error }">
    <label v-if="label" :for="id" class="gh-select-label">
      {{ label }}
      <span v-if="required" class="gh-select-required">*</span>
    </label>
    
    <div class="gh-select-container">
      <select
        :id="id"
        class="gh-select"
        :class="[`gh-select--${size}`]"
        :value="modelValue"
        :disabled="disabled"
        :required="required"
        :multiple="multiple"
        @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
        @blur="$emit('blur', $event)"
        @focus="$emit('focus', $event)"
      >
        <option v-if="placeholder" value="" disabled selected>{{ placeholder }}</option>
        <slot></slot>
      </select>
      
      <Icon 
        name="mdi:chevron-down" 
        class="gh-select-icon" 
      />
    </div>
    
    <p v-if="error" class="gh-select-error">{{ error }}</p>
    <p v-else-if="hint" class="gh-select-hint">{{ hint }}</p>
  </div>
</template>

<script setup lang="ts">
const props = defineProps({
  modelValue: {
    type: [String, Number, Array],
    default: ''
  },
  id: {
    type: String,
    default: () => `select-${Math.random().toString(36).substr(2, 9)}`
  },
  label: {
    type: String,
    default: ''
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
  multiple: {
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
  }
});

defineEmits(['update:modelValue', 'blur', 'focus']);
</script>

<style scoped>
.gh-select-wrapper {
  margin-bottom: var(--spacing-md);
}

.gh-select-label {
  display: block;
  margin-bottom: var(--spacing-xs);
  font-weight: 500;
  font-size: var(--font-size-sm);
  color: var(--color-text-primary);
}

.gh-select-required {
  color: var(--color-danger);
  margin-left: var(--spacing-xxs);
}

.gh-select-container {
  position: relative;
  display: flex;
  align-items: center;
}

.gh-select {
  width: 100%;
  appearance: none;
  border: 1px solid var(--color-border-primary);
  border-radius: var(--border-radius-md);
  background-color: var(--color-bg-canvas);
  color: var(--color-text-primary);
  transition: var(--transition-base);
  font-size: var(--font-size-sm);
  line-height: 1.5;
  padding-right: calc(var(--spacing-md) + 16px);
}

.gh-select:focus {
  outline: none;
  border-color: var(--color-accent-primary);
  box-shadow: 0 0 0 3px rgba(3, 102, 214, 0.3);
}

.gh-select:disabled {
  background-color: var(--color-bg-secondary);
  cursor: not-allowed;
  opacity: 0.7;
}

/* サイズバリエーション */
.gh-select--sm {
  padding: var(--spacing-xxs) var(--spacing-xs);
  font-size: var(--font-size-xs);
}

.gh-select--md {
  padding: var(--spacing-xs) var(--spacing-md);
}

.gh-select--lg {
  padding: var(--spacing-sm) var(--spacing-md);
  font-size: var(--font-size-md);
}

/* アイコン */
.gh-select-icon {
  position: absolute;
  right: var(--spacing-xs);
  color: var(--color-text-tertiary);
  pointer-events: none;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* エラーとヒント */
.gh-select-error {
  margin-top: var(--spacing-xxs);
  color: var(--color-danger);
  font-size: var(--font-size-xs);
}

.gh-select-hint {
  margin-top: var(--spacing-xxs);
  color: var(--color-text-tertiary);
  font-size: var(--font-size-xs);
}

.gh-select-wrapper--error .gh-select {
  border-color: var(--color-danger);
}

.gh-select-wrapper--error .gh-select:focus {
  box-shadow: 0 0 0 3px rgba(215, 58, 73, 0.3);
}
</style>
