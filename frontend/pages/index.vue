<template>
  <div>
    <div class="px-4 py-5 sm:px-6">
      <h3 class="text-lg leading-6 font-medium text-gray-900">
        TicketHub ダッシュボード
      </h3>
      <p class="mt-1 max-w-2xl text-sm text-gray-500">
        チケット管理システム
      </p>
    </div>
    
    <div class="bg-white shadow overflow-hidden sm:rounded-lg mt-5">
      <div class="px-4 py-5 sm:p-6">
        <div class="flex justify-between">
          <h3 class="text-lg leading-6 font-medium text-gray-900">
            システムステータス
          </h3>
          <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
            オンライン
          </span>
        </div>
        
        <div class="mt-6 grid grid-cols-1 gap-5 sm:grid-cols-2">
          <div class="bg-gray-50 overflow-hidden shadow rounded-lg">
            <div class="px-4 py-5 sm:p-6">
              <div class="flex items-center">
                <div class="ml-5 w-0 flex-1">
                  <dl>
                    <dt class="text-sm font-medium text-gray-500 truncate">
                      API サーバー
                    </dt>
                    <dd class="flex items-center">
                      <div class="text-lg font-medium text-gray-900">
                        {{ apiStatus }}
                      </div>
                      <button 
                        class="ml-2 px-2 py-1 text-xs text-blue-600 hover:text-blue-800"
                        @click="checkApiStatus"
                      >
                        確認
                      </button>
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          <div class="bg-gray-50 overflow-hidden shadow rounded-lg">
            <div class="px-4 py-5 sm:p-6">
              <div class="flex items-center">
                <div class="ml-5 w-0 flex-1">
                  <dl>
                    <dt class="text-sm font-medium text-gray-500 truncate">
                      WebSocket 接続
                    </dt>
                    <dd class="flex items-center">
                      <div class="text-lg font-medium text-gray-900">
                        {{ wsStatus }}
                      </div>
                      <button 
                        class="ml-2 px-2 py-1 text-xs text-blue-600 hover:text-blue-800"
                        @click="checkWsStatus"
                      >
                        確認
                      </button>
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import axios from 'axios'
import { onMounted, ref } from 'vue'

const apiStatus = ref('確認中...')
const wsStatus = ref('未接続')

const config = useRuntimeConfig()

// APIステータス確認
async function checkApiStatus() {
  apiStatus.value = '確認中...'
  try {
    const response = await axios.get(`${config.public.apiBaseUrl}/health`)
    if (response.data && response.data.status === 'healthy') {
      apiStatus.value = '接続済み'
    } else {
      apiStatus.value = '応答エラー'
    }
  } catch (error) {
    console.error('API接続エラー:', error)
    apiStatus.value = '接続エラー'
  }
}

// WebSocket接続確認 (現在は単純なステータス表示のみ)
function checkWsStatus() {
  wsStatus.value = 'WebSocket未実装'
  // 実際のWebSocket実装はTODO
}

onMounted(() => {
  checkApiStatus()
})
</script>
