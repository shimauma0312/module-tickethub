<template>
  <div 
    class="gh-card" 
    :class="[
      { 'gh-card--hoverable': hoverable },
      { 'gh-card--selected': selected }
    ]"
  >
    <div v-if="$slots.header || title" class="gh-card__header">
      <slot name="header">
        <h3 class="gh-card__title">{{ title }}</h3>
        <div v-if="$slots.headerActions" class="gh-card__actions">
          <slot name="headerActions"></slot>
        </div>
      </slot>
    </div>
    
    <div class="gh-card__body">
      <slot></slot>
    </div>
    
    <div v-if="$slots.footer" class="gh-card__footer">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  title: {
    type: String,
    default: ''
  },
  hoverable: {
    type: Boolean,
    default: false
  },
  selected: {
    type: Boolean,
    default: false
  }
});
</script>

<style scoped>
.gh-card {
  background-color: var(--color-bg-canvas);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--border-radius-md);
  overflow: hidden;
  transition: var(--transition-base);
}

.gh-card--hoverable {
  cursor: pointer;
}

.gh-card--hoverable:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
  border-color: var(--color-border-primary);
}

.gh-card--selected {
  border-color: var(--color-accent-primary);
  box-shadow: 0 0 0 1px var(--color-accent-primary);
}

.gh-card__header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--color-border-primary);
  background-color: var(--color-bg-secondary);
}

.gh-card__title {
  margin: 0;
  font-size: var(--font-size-md);
  font-weight: 600;
  color: var(--color-text-primary);
}

.gh-card__actions {
  display: flex;
  gap: var(--spacing-xs);
}

.gh-card__body {
  padding: var(--spacing-md);
}

.gh-card__footer {
  padding: var(--spacing-md);
  border-top: 1px solid var(--color-border-primary);
  background-color: var(--color-bg-secondary);
}

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .gh-card__header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--spacing-xs);
  }
  
  .gh-card__actions {
    align-self: flex-end;
  }
}
</style>
