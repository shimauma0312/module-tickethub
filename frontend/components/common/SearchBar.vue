<template>
  <div class="search-bar-container">
    <Input
      v-model="internalSearchQuery"
      :placeholder="placeholder"
      prepend-icon="mdi:magnify"
      :clearable="true"
      @keyup.enter="onSearch"
      @update:modelValue="onInput"
    />
    <Button @click="onSearch" :loading="loading">{{ buttonText }}</Button>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
// Input と Button コンポーネントがグローバル登録されているか、
// もしくはローカルでインポートする必要がある点に注意してください。
// この例ではグローバル登録されている前提です。

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: '検索キーワードを入力...'
  },
  buttonText: {
    type: String,
    default: '検索'
  },
  loading: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['update:modelValue', 'search']);

const internalSearchQuery = ref(props.modelValue);

watch(() => props.modelValue, (newValue) => {
  internalSearchQuery.value = newValue;
});

function onInput(value) {
  internalSearchQuery.value = value;
  emit('update:modelValue', value);
}

function onSearch() {
  emit('search', internalSearchQuery.value);
}
</script>

<style scoped>
.search-bar-container {
  display: flex;
  gap: 0.75rem;
  align-items: center;
}

.search-bar-container :deep(input) {
  font-size: 1rem; /* 必要に応じて調整 */
}
</style>
