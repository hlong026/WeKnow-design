<template>
  <t-drawer
    v-model:visible="drawerVisible"
    :header="$t('globalSearch.title')"
    :footer="false"
    size="600px"
    placement="right"
    :close-on-overlay-click="true"
    @close="handleClose"
  >
    <div class="global-search-panel">
      <!-- 搜索输入框 -->
      <div class="search-input-wrapper">
        <t-input
          v-model="searchKeyword"
          :placeholder="$t('globalSearch.placeholder')"
          clearable
          size="large"
          @enter="handleSearch"
          @clear="handleClear"
        >
          <template #prefix-icon>
            <t-icon name="search" />
          </template>
        </t-input>
        <t-button theme="primary" @click="handleSearch" :loading="loading">
          {{ $t('common.search') }}
        </t-button>
      </div>

      <!-- 搜索结果 -->
      <div class="search-results">
        <div v-if="loading" class="loading-state">
          <t-loading />
          <span>{{ $t('globalSearch.searching') }}</span>
        </div>

        <div v-else-if="hasSearched && results.length === 0" class="empty-state">
          <t-icon name="search" size="48px" />
          <span>{{ $t('globalSearch.noResults') }}</span>
        </div>

        <div v-else-if="results.length > 0" class="results-list">
          <div class="results-header">
            <span>{{ $t('globalSearch.resultsCount', { count: total }) }}</span>
          </div>
          
          <div 
            v-for="item in results" 
            :key="item.id" 
            class="result-item"
            @click="handleResultClick(item)"
          >
            <div class="result-header">
              <span class="kb-name">{{ item.knowledge_base_name }}</span>
              <span class="score-badge" :class="getScoreClass(item.score)">
                {{ formatScore(item.score) }}
              </span>
            </div>
            <div class="result-title">{{ item.knowledge_title || $t('globalSearch.untitled') }}</div>
            <div class="result-content" v-html="highlightKeyword(item.content)"></div>
          </div>

          <!-- 分页 -->
          <div v-if="total > pageSize" class="pagination-wrapper">
            <t-pagination
              v-model:current="currentPage"
              :total="total"
              :page-size="pageSize"
              :show-page-size="false"
              @current-change="handlePageChange"
            />
          </div>
        </div>

        <div v-else class="initial-state">
          <t-icon name="search" size="48px" />
          <span>{{ $t('globalSearch.hint') }}</span>
        </div>
      </div>
    </div>
  </t-drawer>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { globalSearchKnowledge } from '@/api/knowledge-base'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
  visible: boolean
}>()

const emit = defineEmits<{
  (e: 'update:visible', value: boolean): void
}>()

const router = useRouter()
const { t } = useI18n()

const drawerVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

interface SearchResult {
  id: string
  content: string
  knowledge_id: string
  knowledge_title: string
  knowledge_base_id: string
  knowledge_base_name: string
  score: number
  match_type: string
}

const searchKeyword = ref('')
const results = ref<SearchResult[]>([])
const loading = ref(false)
const hasSearched = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = 20

let debounceTimer: ReturnType<typeof setTimeout> | null = null

const handleSearch = async () => {
  if (!searchKeyword.value.trim()) return
  
  loading.value = true
  hasSearched.value = true
  
  try {
    const res = await globalSearchKnowledge({
      keyword: searchKeyword.value.trim(),
      page: currentPage.value,
      page_size: pageSize
    }) as any
    
    if (res.success) {
      results.value = res.data || []
      total.value = res.total || 0
    }
  } catch (error) {
    console.error('Search failed:', error)
    results.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handleClear = () => {
  results.value = []
  total.value = 0
  hasSearched.value = false
  currentPage.value = 1
}

const handleClose = () => {
  emit('update:visible', false)
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  handleSearch()
}

const handleResultClick = (item: SearchResult) => {
  emit('update:visible', false)
  router.push(`/platform/knowledge-bases/${item.knowledge_base_id}`)
}

const highlightKeyword = (content: string) => {
  if (!searchKeyword.value.trim() || !content) return content
  const keyword = searchKeyword.value.trim()
  const regex = new RegExp(`(${keyword.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')})`, 'gi')
  const truncated = content.length > 200 ? content.substring(0, 200) + '...' : content
  return truncated.replace(regex, '<mark>$1</mark>')
}

const formatScore = (score: number) => {
  return (score * 100).toFixed(0) + '%'
}

const getScoreClass = (score: number) => {
  if (score >= 0.8) return 'high'
  if (score >= 0.5) return 'medium'
  return 'low'
}

// 监听关键词变化，防抖搜索
watch(searchKeyword, (newVal) => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
  if (newVal.trim()) {
    debounceTimer = setTimeout(() => {
      currentPage.value = 1
      handleSearch()
    }, 300)
  }
})

// 重置状态当抽屉关闭
watch(() => props.visible, (newVal) => {
  if (!newVal) {
    // 保留搜索关键词，但可以选择清空
    // searchKeyword.value = ''
    // results.value = []
    // hasSearched.value = false
  }
})
</script>

<style scoped lang="less">
.global-search-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.search-input-wrapper {
  display: flex;
  gap: 12px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e7ebf0;
}

.search-results {
  flex: 1;
  overflow-y: auto;
  padding-top: 16px;
}

.loading-state,
.empty-state,
.initial-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 300px;
  color: #999;
  gap: 12px;
}

.results-header {
  color: #666;
  font-size: 14px;
  margin-bottom: 12px;
}

.results-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.result-item {
  padding: 16px;
  border: 1px solid #e7ebf0;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;

  &:hover {
    border-color: #07c05f;
    box-shadow: 0 2px 8px rgba(7, 192, 95, 0.1);
  }
}

.result-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.kb-name {
  font-size: 12px;
  color: #07c05f;
  background: rgba(7, 192, 95, 0.1);
  padding: 2px 8px;
  border-radius: 4px;
}

.score-badge {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;

  &.high {
    color: #07c05f;
    background: rgba(7, 192, 95, 0.1);
  }

  &.medium {
    color: #f59e0b;
    background: rgba(245, 158, 11, 0.1);
  }

  &.low {
    color: #999;
    background: #f5f5f5;
  }
}

.result-title {
  font-size: 15px;
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.result-content {
  font-size: 13px;
  color: #666;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
  overflow: hidden;

  :deep(mark) {
    background: #fff3cd;
    color: #856404;
    padding: 0 2px;
    border-radius: 2px;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding-top: 16px;
  margin-top: 8px;
  border-top: 1px solid #e7ebf0;
}
</style>
