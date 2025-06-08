<template>
  <nav class="gh-nav" :class="{ 'gh-nav--expanded': isExpanded }">
    <div class="gh-nav__header">
      <button 
        v-if="collapsible" 
        class="gh-nav__toggle"
        @click="isExpanded = !isExpanded"
        :aria-expanded="isExpanded ? 'true' : 'false'"
        :aria-label="isExpanded ? 'ナビゲーションを折りたたむ' : 'ナビゲーションを展開する'"
      >
        <Icon :name="isExpanded ? 'mdi:menu-open' : 'mdi:menu'" />
      </button>
      
      <slot name="header">
        <div class="gh-nav__logo">
          <Icon name="mdi:ticket" class="gh-nav__logo-icon" />
          <span v-if="!collapsible || isExpanded" class="gh-nav__logo-text">TicketHub</span>
        </div>
      </slot>
    </div>
    
    <div class="gh-nav__content">
      <div v-for="(section, idx) in sections" :key="idx" class="gh-nav__section">
        <div v-if="section.title && (!collapsible || isExpanded)" class="gh-nav__section-title">
          {{ section.title }}
        </div>
        
        <div class="gh-nav__menu">
          <NuxtLink
            v-for="item in section.items"
            :key="item.to"
            :to="item.to"
            :aria-label="item.label"
            :title="collapsible && !isExpanded ? item.label : ''"
            class="gh-nav__item"
            :class="{ 'gh-nav__item--active': isActiveRoute(item.to) }"
          >
            <Icon v-if="item.icon" :name="item.icon" class="gh-nav__item-icon" />
            <span v-if="!collapsible || isExpanded" class="gh-nav__item-label">{{ item.label }}</span>
            <span 
              v-if="item.badge && (!collapsible || isExpanded)" 
              class="gh-nav__item-badge"
              :class="{ [`gh-nav__item-badge--${item.badgeVariant || 'default'}`]: true }"
            >
              {{ item.badge }}
            </span>
          </NuxtLink>
        </div>
      </div>
    </div>
    
    <div v-if="$slots.footer" class="gh-nav__footer">
      <slot name="footer"></slot>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';

interface NavItem {
  label: string;
  to: string;
  icon?: string;
  badge?: string | number;
  badgeVariant?: 'default' | 'primary' | 'secondary' | 'success' | 'danger' | 'warning';
}

interface NavSection {
  title?: string;
  items: NavItem[];
}

const props = defineProps({
  sections: {
    type: Array as () => NavSection[],
    default: () => []
  },
  collapsible: {
    type: Boolean,
    default: false
  },
  defaultExpanded: {
    type: Boolean,
    default: true
  }
});

const route = useRoute();
const isExpanded = ref(props.defaultExpanded);

// 現在のルートがアクティブかどうかを判定
function isActiveRoute(path: string): boolean {
  if (path === '/') {
    return route.path === '/';
  }
  return route.path.startsWith(path);
}
</script>

<style scoped>
.gh-nav {
  display: flex;
  flex-direction: column;
  background-color: var(--color-bg-canvas);
  border-right: 1px solid var(--color-border-primary);
  height: 100%;
  transition: var(--transition-base);
}

.gh-nav--expanded {
  width: 240px;
}

.gh-nav:not(.gh-nav--expanded) {
  width: 60px;
}

.gh-nav__header {
  display: flex;
  align-items: center;
  padding: var(--spacing-sm);
  border-bottom: 1px solid var(--color-border-primary);
  height: 60px;
}

.gh-nav__toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  background: none;
  border: none;
  border-radius: var(--border-radius-md);
  padding: var(--spacing-xxs);
  color: var(--color-text-primary);
  cursor: pointer;
  transition: var(--transition-base);
}

.gh-nav__toggle:hover {
  background-color: var(--color-bg-secondary);
}

.gh-nav__logo {
  display: flex;
  align-items: center;
  margin-left: var(--spacing-xs);
  font-weight: 600;
  font-size: var(--font-size-md);
  color: var(--color-text-primary);
}

.gh-nav__logo-icon {
  font-size: 24px;
  color: var(--color-accent-primary);
}

.gh-nav__logo-text {
  margin-left: var(--spacing-xs);
  white-space: nowrap;
}

.gh-nav__content {
  flex: 1;
  overflow-y: auto;
  padding: var(--spacing-sm) 0;
}

.gh-nav__section {
  margin-bottom: var(--spacing-md);
}

.gh-nav__section:last-child {
  margin-bottom: 0;
}

.gh-nav__section-title {
  padding: 0 var(--spacing-sm);
  margin-bottom: var(--spacing-xs);
  font-size: var(--font-size-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
}

.gh-nav__menu {
  display: flex;
  flex-direction: column;
}

.gh-nav__item {
  display: flex;
  align-items: center;
  padding: var(--spacing-xs) var(--spacing-sm);
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: var(--font-size-sm);
  border-left: 3px solid transparent;
  transition: var(--transition-base);
}

.gh-nav__item:hover {
  background-color: var(--color-bg-secondary);
  color: var(--color-text-primary);
}

.gh-nav__item--active {
  color: var(--color-accent-primary);
  background-color: rgba(3, 102, 214, 0.08);
  border-left-color: var(--color-accent-primary);
  font-weight: 500;
}

.gh-nav__item-icon {
  margin-right: var(--spacing-sm);
  font-size: 18px;
}

.gh-nav:not(.gh-nav--expanded) .gh-nav__item-icon {
  margin-right: 0;
}

.gh-nav__item-label {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.gh-nav__item-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 20px;
  height: 20px;
  padding: 0 var(--spacing-xxs);
  border-radius: var(--border-radius-full);
  font-size: var(--font-size-xs);
  font-weight: 600;
  margin-left: var(--spacing-xs);
}

.gh-nav__item-badge--default {
  background-color: var(--color-bg-secondary);
  color: var(--color-text-secondary);
}

.gh-nav__item-badge--primary {
  background-color: var(--color-accent-primary);
  color: white;
}

.gh-nav__item-badge--secondary {
  background-color: var(--color-accent-secondary);
  color: white;
}

.gh-nav__item-badge--success {
  background-color: var(--color-success);
  color: white;
}

.gh-nav__item-badge--danger {
  background-color: var(--color-danger);
  color: white;
}

.gh-nav__item-badge--warning {
  background-color: var(--color-warning);
  color: var(--color-text-primary);
}

.gh-nav__footer {
  padding: var(--spacing-sm);
  border-top: 1px solid var(--color-border-primary);
}

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .gh-nav {
    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    z-index: 1000;
    transform: translateX(-100%);
  }
  
  .gh-nav--expanded {
    transform: translateX(0);
  }
  
  .gh-nav__toggle {
    position: fixed;
    top: var(--spacing-sm);
    left: var(--spacing-sm);
    z-index: 1001;
    background-color: var(--color-bg-canvas);
    box-shadow: var(--shadow-md);
  }
}
</style>
