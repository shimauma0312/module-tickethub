<template>
  <div class="gh-table-wrapper">
    <div v-if="$slots.header" class="gh-table__header">
      <slot name="header"></slot>
    </div>
    
    <div class="gh-table__container" :class="{ 'gh-table__container--loading': loading }">
      <table class="gh-table">
        <thead>
          <tr>
            <th 
              v-for="(column, idx) in columns" 
              :key="idx"
              :class="[
                { 'gh-table__th--sortable': column.sortable },
                { 'gh-table__th--sorted-asc': sortBy === column.key && sortDirection === 'asc' },
                { 'gh-table__th--sorted-desc': sortBy === column.key && sortDirection === 'desc' }
              ]"
              :style="column.width ? { width: column.width } : {}"
              @click="column.sortable ? handleSort(column.key) : null"
            >
              <div class="gh-table__th-content">
                {{ column.label }}
                <Icon 
                  v-if="column.sortable"
                  :name="getSortIcon(column.key)"
                  class="gh-table__sort-icon"
                />
              </div>
            </th>
            <th v-if="$slots.rowActions" class="gh-table__th--actions"></th>
          </tr>
        </thead>
        
        <tbody>
          <template v-if="items.length">
            <tr 
              v-for="(item, rowIdx) in items" 
              :key="getItemKey(item, rowIdx)"
              :class="{ 'gh-table__row--selected': isRowSelected(item) }"
              @click="handleRowClick(item, rowIdx)"
            >
              <td 
                v-for="(column, colIdx) in columns" 
                :key="colIdx"
                :class="column.cellClass"
              >
                <slot 
                  :name="`cell(${column.key})`" 
                  :item="item" 
                  :value="item[column.key]"
                  :index="rowIdx"
                >
                  {{ item[column.key] }}
                </slot>
              </td>
              <td v-if="$slots.rowActions" class="gh-table__cell--actions">
                <slot name="rowActions" :item="item" :index="rowIdx"></slot>
              </td>
            </tr>
          </template>
          <tr v-else-if="!loading" class="gh-table__empty">
            <td :colspan="$slots.rowActions ? columns.length + 1 : columns.length">
              <slot name="emptyState">
                <div class="gh-table__empty-state">
                  <Icon name="mdi:table-off" size="36" />
                  <p>{{ emptyText }}</p>
                </div>
              </slot>
            </td>
          </tr>
        </tbody>
      </table>
      
      <div v-if="loading" class="gh-table__loading">
        <slot name="loadingState">
          <div class="gh-table__loading-spinner"></div>
          <p>{{ loadingText }}</p>
        </slot>
      </div>
    </div>
    
    <div v-if="$slots.footer || pagination" class="gh-table__footer">
      <slot name="footer">
        <div v-if="pagination" class="gh-table__pagination">
          <button 
            class="gh-table__pagination-btn"
            :disabled="pagination.currentPage <= 1"
            @click="$emit('page-change', pagination.currentPage - 1)"
          >
            <Icon name="mdi:chevron-left" />
          </button>
          
          <span class="gh-table__pagination-info">
            {{ getPaginationText() }}
          </span>
          
          <button 
            class="gh-table__pagination-btn"
            :disabled="pagination.currentPage >= pagination.totalPages"
            @click="$emit('page-change', pagination.currentPage + 1)"
          >
            <Icon name="mdi:chevron-right" />
          </button>
        </div>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

interface Column {
  key: string;
  label: string;
  sortable?: boolean;
  width?: string;
  cellClass?: string | string[];
}

interface Pagination {
  currentPage: number;
  totalPages: number;
  totalItems: number;
  itemsPerPage: number;
}

const props = defineProps({
  columns: {
    type: Array as () => Column[],
    required: true
  },
  items: {
    type: Array as () => any[],
    default: () => []
  },
  itemKey: {
    type: String,
    default: 'id'
  },
  loading: {
    type: Boolean,
    default: false
  },
  emptyText: {
    type: String,
    default: 'データがありません'
  },
  loadingText: {
    type: String,
    default: '読み込み中...'
  },
  sortBy: {
    type: String,
    default: ''
  },
  sortDirection: {
    type: String,
    default: 'asc',
    validator: (value: string) => ['asc', 'desc'].includes(value)
  },
  selectedRows: {
    type: Array,
    default: () => []
  },
  pagination: {
    type: Object as () => Pagination,
    default: null
  }
});

const emit = defineEmits([
  'row-click', 
  'sort', 
  'page-change'
]);

// 行の一意のキーを取得
function getItemKey(item: any, index: number): string | number {
  return props.itemKey && item[props.itemKey] ? item[props.itemKey] : index;
}

// ソートアイコンの取得
function getSortIcon(columnKey: string): string {
  if (columnKey !== props.sortBy) {
    return 'mdi:unfold-more-horizontal';
  }
  return props.sortDirection === 'asc' ? 'mdi:arrow-up' : 'mdi:arrow-down';
}

// ソート処理
function handleSort(columnKey: string) {
  let direction = 'asc';
  
  if (props.sortBy === columnKey) {
    direction = props.sortDirection === 'asc' ? 'desc' : 'asc';
  }
  
  emit('sort', { key: columnKey, direction });
}

// 行クリック処理
function handleRowClick(item: any, index: number) {
  emit('row-click', { item, index });
}

// 行の選択状態
function isRowSelected(item: any): boolean {
  if (!props.selectedRows.length) return false;
  
  const itemId = item[props.itemKey];
  return props.selectedRows.some(row => 
    typeof row === 'object' ? row[props.itemKey] === itemId : row === itemId
  );
}

// ページネーションテキスト
function getPaginationText(): string {
  if (!props.pagination) return '';
  
  const { currentPage, totalPages, totalItems, itemsPerPage } = props.pagination;
  const start = (currentPage - 1) * itemsPerPage + 1;
  const end = Math.min(currentPage * itemsPerPage, totalItems);
  
  return `${start}-${end} / ${totalItems}件`;
}
</script>

<style scoped>
.gh-table-wrapper {
  width: 100%;
  overflow: hidden;
  border: 1px solid var(--color-border-primary);
  border-radius: var(--border-radius-md);
  background-color: var(--color-bg-canvas);
}

.gh-table__header {
  padding: var(--spacing-md);
  border-bottom: 1px solid var(--color-border-primary);
  background-color: var(--color-bg-secondary);
}

.gh-table__container {
  position: relative;
  overflow-x: auto;
}

.gh-table__container--loading {
  min-height: 200px;
}

.gh-table {
  width: 100%;
  border-collapse: collapse;
  font-size: var(--font-size-sm);
}

.gh-table thead {
  background-color: var(--color-bg-secondary);
}

.gh-table th {
  position: sticky;
  top: 0;
  padding: var(--spacing-xs) var(--spacing-md);
  text-align: left;
  font-weight: 600;
  color: var(--color-text-primary);
  border-bottom: 1px solid var(--color-border-primary);
  background-color: var(--color-bg-secondary);
  z-index: 1;
}

.gh-table__th-content {
  display: flex;
  align-items: center;
}

.gh-table__th--sortable {
  cursor: pointer;
  user-select: none;
}

.gh-table__th--sortable:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.gh-table__sort-icon {
  margin-left: var(--spacing-xxs);
  opacity: 0.5;
}

.gh-table__th--sorted-asc .gh-table__sort-icon,
.gh-table__th--sorted-desc .gh-table__sort-icon {
  opacity: 1;
}

.gh-table td {
  padding: var(--spacing-sm) var(--spacing-md);
  border-bottom: 1px solid var(--color-border-primary);
  color: var(--color-text-primary);
}

.gh-table tr:last-child td {
  border-bottom: none;
}

.gh-table tbody tr {
  transition: var(--transition-base);
}

.gh-table tbody tr:hover {
  background-color: rgba(0, 0, 0, 0.02);
}

.gh-table__row--selected {
  background-color: rgba(3, 102, 214, 0.08) !important;
}

.gh-table__th--actions,
.gh-table__cell--actions {
  width: 1%;
  white-space: nowrap;
  text-align: right;
}

.gh-table__empty {
  text-align: center;
}

.gh-table__empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-xl) 0;
  color: var(--color-text-tertiary);
}

.gh-table__empty-state p {
  margin-top: var(--spacing-sm);
}

.gh-table__loading {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: rgba(255, 255, 255, 0.8);
  z-index: 2;
}

.gh-table__loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid var(--color-border-primary);
  border-radius: 50%;
  border-top-color: var(--color-accent-primary);
  animation: gh-table-spin 1s linear infinite;
  margin-bottom: var(--spacing-sm);
}

@keyframes gh-table-spin {
  to {
    transform: rotate(360deg);
  }
}

.gh-table__footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--spacing-sm) var(--spacing-md);
  border-top: 1px solid var(--color-border-primary);
  background-color: var(--color-bg-secondary);
}

.gh-table__pagination {
  display: flex;
  align-items: center;
  margin-left: auto;
}

.gh-table__pagination-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid var(--color-border-primary);
  border-radius: var(--border-radius-md);
  background-color: var(--color-bg-canvas);
  color: var(--color-text-primary);
  cursor: pointer;
  transition: var(--transition-base);
}

.gh-table__pagination-btn:hover:not(:disabled) {
  background-color: var(--color-bg-secondary);
}

.gh-table__pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.gh-table__pagination-info {
  margin: 0 var(--spacing-sm);
  font-size: var(--font-size-sm);
  color: var(--color-text-secondary);
}

/* レスポンシブデザイン */
@media (max-width: 768px) {
  .gh-table__footer {
    flex-direction: column;
    gap: var(--spacing-sm);
  }
  
  .gh-table__pagination {
    margin: 0 auto;
  }
}
</style>
