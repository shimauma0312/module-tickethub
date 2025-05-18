<template>
  <div>
    <div class="px-4 py-5 sm:px-6 flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-semibold text-gray-900">Issues</h1>
        <p class="mt-1 text-sm text-gray-500">課題管理・追跡</p>
      </div>
      <button 
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
      >
        新規Issue作成
      </button>
    </div>

    <div class="bg-white shadow overflow-hidden sm:rounded-md">
      <ul class="divide-y divide-gray-200">
        <li v-if="issues.length === 0" class="p-4 text-center text-gray-500">
          まだIssueはありません。「新規Issue作成」ボタンから登録してください。
        </li>
        <li v-for="issue in issues" :key="issue.id" class="px-4 py-4 sm:px-6 hover:bg-gray-50 cursor-pointer">
          <div class="flex items-center justify-between">
            <p class="text-sm font-medium text-indigo-600 truncate">
              {{ issue.title }}
            </p>
            <div class="ml-2 flex-shrink-0 flex">
              <p class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full" 
                 :class="getStatusClass(issue.status)">
                {{ issue.status }}
              </p>
            </div>
          </div>
          <div class="mt-2 sm:flex sm:justify-between">
            <div class="sm:flex">
              <p class="flex items-center text-sm text-gray-500">
                #{{ issue.id }} が {{ issue.created_by }} により {{ issue.created_at }} に作成
              </p>
            </div>
            <div class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0">
              <span>{{ issue.comments_count }} コメント</span>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'

// サンプルデータ
const issues = ref([
  {
    id: 1,
    title: 'アーキテクチャ図の改善',
    status: 'open',
    created_by: 'shimauma0312',
    created_at: '2025-05-01',
    comments_count: 3
  },
  {
    id: 2,
    title: 'プロジェクト初期構成とDockerセットアップ',
    status: 'in-progress',
    created_by: 'shimauma0312',
    created_at: '2025-05-15',
    comments_count: 2
  }
])

// ステータスに応じたクラスを返す
function getStatusClass(status) {
  switch (status) {
    case 'open':
      return 'bg-green-100 text-green-800'
    case 'in-progress':
      return 'bg-yellow-100 text-yellow-800'
    case 'closed':
      return 'bg-gray-100 text-gray-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

onMounted(async () => {
  // APIからデータを取得する（将来実装）
  // const config = useRuntimeConfig()
  // try {
  //   const response = await fetch(`${config.public.apiBaseUrl}/api/v1/issues`)
  //   issues.value = await response.json()
  // } catch (error) {
  //   console.error('APIからデータを取得できませんでした:', error)
  // }
})
</script>
