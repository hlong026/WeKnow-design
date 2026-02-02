<template>
  <div 
    class="provider-card" 
    :class="{ 
      'configured': hasCredential,
      'has-error': credentialStatus === 'invalid'
    }"
  >
    <div class="provider-header">
      <div class="provider-icon-wrapper">
        <img 
          v-if="provider.icon" 
          :src="provider.icon" 
          :alt="providerDisplayName"
          class="provider-icon"
        />
        <div v-else class="provider-icon-placeholder">
          {{ providerDisplayName.charAt(0) }}
        </div>
      </div>
      
      <div class="provider-info">
        <h4 class="provider-name">{{ providerDisplayName }}</h4>
        <p class="provider-description">{{ provider.description }}</p>
        <div class="provider-meta">
          <span class="model-count">
            <t-icon name="layers" size="14px" />
            {{ modelCount }} 个模型
          </span>
          <span v-if="hasCredential" class="credential-status" :class="credentialStatus">
            <t-icon :name="statusIcon" size="14px" />
            {{ statusText }}
          </span>
        </div>
      </div>
    </div>
    
    <div class="provider-actions">
      <t-button 
        v-if="!hasCredential" 
        theme="primary" 
        size="small"
        @click="$emit('configure')"
      >
        <template #icon><t-icon name="setting" /></template>
        配置
      </t-button>
      <template v-else>
        <t-button 
          theme="default" 
          variant="outline"
          size="small"
          @click="$emit('manage')"
        >
          管理模型
        </t-button>
        <t-dropdown :options="moreOptions" @click="handleMoreAction">
          <t-button variant="text" shape="square" size="small">
            <t-icon name="more" />
          </t-button>
        </t-dropdown>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Provider {
  name: string
  displayName: string
  display_name?: string
  description: string
  icon?: string
}

const props = defineProps<{
  provider: Provider
  hasCredential: boolean
  credentialStatus?: 'active' | 'invalid' | 'quota_exceeded'
  modelCount: number
}>()

const emit = defineEmits<{
  configure: []
  manage: []
  edit: []
  delete: []
  test: []
}>()

// 兼容 displayName 和 display_name
const providerDisplayName = computed(() => {
  return props.provider.displayName || props.provider.display_name || props.provider.name
})

const statusIcon = computed(() => {
  switch (props.credentialStatus) {
    case 'active': return 'check-circle'
    case 'invalid': return 'error-circle'
    case 'quota_exceeded': return 'info-circle'
    default: return 'check-circle'
  }
})

const statusText = computed(() => {
  switch (props.credentialStatus) {
    case 'active': return '已连接'
    case 'invalid': return '凭证无效'
    case 'quota_exceeded': return '配额超限'
    default: return '已连接'
  }
})

const moreOptions = [
  { content: '编辑凭证', value: 'edit' },
  { content: '测试连接', value: 'test' },
  { content: '删除', value: 'delete', theme: 'error' }
]

const handleMoreAction = (data: { value: string }) => {
  switch (data.value) {
    case 'edit': emit('edit'); break
    case 'test': emit('test'); break
    case 'delete': emit('delete'); break
  }
}
</script>

<style scoped lang="less">
.provider-card {
  background: var(--td-bg-color-container);
  border: 1px solid var(--td-border-level-1-color);
  border-radius: 8px;
  padding: 16px;
  transition: all 0.2s ease;
  
  &:hover {
    border-color: var(--td-brand-color);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  }
  
  &.configured {
    border-color: var(--td-success-color-3);
  }
  
  &.has-error {
    border-color: var(--td-error-color-3);
  }
}

.provider-header {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.provider-icon-wrapper {
  flex-shrink: 0;
}

.provider-icon {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  object-fit: contain;
}

.provider-icon-placeholder {
  width: 40px;
  height: 40px;
  border-radius: 8px;
  background: var(--td-brand-color-light);
  color: var(--td-brand-color);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 18px;
}

.provider-info {
  flex: 1;
  min-width: 0;
}

.provider-name {
  margin: 0 0 4px;
  font-size: 14px;
  font-weight: 600;
  color: var(--td-text-color-primary);
}

.provider-description {
  margin: 0 0 8px;
  font-size: 12px;
  color: var(--td-text-color-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.provider-meta {
  display: flex;
  gap: 12px;
  font-size: 12px;
  color: var(--td-text-color-placeholder);
}

.model-count, .credential-status {
  display: flex;
  align-items: center;
  gap: 4px;
}

.credential-status {
  &.active { color: var(--td-success-color); }
  &.invalid { color: var(--td-error-color); }
  &.quota_exceeded { color: var(--td-warning-color); }
}

.provider-actions {
  display: flex;
  gap: 8px;
  justify-content: flex-end;
}
</style>
