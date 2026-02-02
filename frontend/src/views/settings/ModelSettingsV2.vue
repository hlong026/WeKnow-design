<template>
  <div class="model-settings-v2">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>{{ $t('modelSettings.title') }}</h2>
      <p class="page-description">{{ $t('modelSettings.description') }}</p>
    </div>

    <!-- Tab 切换 -->
    <t-tabs v-model="activeTab" class="settings-tabs">
      <!-- 厂商管理 Tab -->
      <t-tab-panel value="providers" :label="$t('modelSettings.tabs.providers')">
        <div class="tab-content">
          <div class="section-header">
            <h3>已配置的厂商</h3>
            <t-button theme="primary" @click="showProviderSelector = true">
              <template #icon><t-icon name="add" /></template>
              添加厂商
            </t-button>
          </div>
          
          <!-- 已配置厂商列表 -->
          <div class="provider-grid" v-if="configuredProviders.length > 0">
            <ProviderCard
              v-for="provider in configuredProviders"
              :key="provider.name"
              :provider="{
                name: provider.name,
                displayName: provider.display_name,
                description: provider.description,
                icon: provider.icon
              }"
              :has-credential="true"
              :credential-status="getCredentialStatus(provider.name)"
              :model-count="getModelCount(provider.name)"
              @manage="openProviderManage(provider)"
              @edit="openCredentialEdit(provider)"
              @delete="deleteProviderCredential(provider)"
              @test="testProviderCredential(provider)"
            />
          </div>
          
          <!-- 空状态 -->
          <div v-else class="empty-state">
            <t-icon name="cloud" size="48px" class="empty-icon" />
            <p>还没有配置任何模型厂商</p>
            <t-button theme="primary" @click="showProviderSelector = true">
              添加第一个厂商
            </t-button>
          </div>
        </div>
      </t-tab-panel>

      <!-- 模型列表 Tab -->
      <t-tab-panel value="models" :label="$t('modelSettings.tabs.models')">
        <div class="tab-content">
          <!-- 筛选栏 -->
          <div class="filter-bar">
            <t-select
              v-model="filterType"
              :options="modelTypeOptions"
              placeholder="模型类型"
              clearable
              style="width: 150px"
            />
            <t-select
              v-model="filterProvider"
              :options="providerOptions"
              placeholder="厂商"
              clearable
              style="width: 150px"
            />
            <t-input
              v-model="searchKeyword"
              placeholder="搜索模型名称"
              clearable
              style="width: 200px"
            >
              <template #prefix-icon><t-icon name="search" /></template>
            </t-input>
            <div class="filter-spacer"></div>
            <t-button theme="primary" @click="openAddModel">
              <template #icon><t-icon name="add" /></template>
              添加模型
            </t-button>
          </div>

          <!-- 模型列表 -->
          <div class="model-list" v-loading="loading">
            <div
              v-for="model in filteredModels"
              :key="model.id"
              class="model-item"
              :class="{ 'is-builtin': model.is_builtin }"
            >
              <div class="model-info">
                <div class="model-name">
                  {{ model.name }}
                  <t-tag v-if="model.is_builtin" theme="primary" size="small">内置</t-tag>
                  <t-tag v-if="model.is_default" theme="success" size="small">默认</t-tag>
                </div>
                <div class="model-meta">
                  <span class="meta-item">
                    <t-icon name="layers" size="14px" />
                    {{ getModelTypeLabel(model.type) }}
                  </span>
                  <span class="meta-item">
                    <t-icon name="cloud" size="14px" />
                    {{ getProviderLabel(model.parameters?.provider) }}
                  </span>
                </div>
              </div>
              <div class="model-actions">
                <t-button
                  variant="text"
                  size="small"
                  :disabled="model.is_builtin"
                  @click="editModel(model)"
                >
                  编辑
                </t-button>
                <t-popconfirm
                  content="确定要删除这个模型吗？"
                  @confirm="handleDeleteModel(model.id)"
                >
                  <t-button
                    variant="text"
                    theme="danger"
                    size="small"
                    :disabled="model.is_builtin"
                  >
                    删除
                  </t-button>
                </t-popconfirm>
              </div>
            </div>
            
            <!-- 空状态 -->
            <div v-if="filteredModels.length === 0 && !loading" class="empty-state">
              <p>没有找到匹配的模型</p>
            </div>
          </div>
        </div>
      </t-tab-panel>
    </t-tabs>

    <!-- 厂商选择器弹窗 -->
    <t-dialog
      v-model:visible="showProviderSelector"
      header="选择模型厂商"
      width="720px"
      :footer="false"
    >
      <div class="provider-selector">
        <div
          v-for="provider in availableProviders"
          :key="provider.name"
          class="provider-option"
          @click="selectProvider(provider)"
        >
          <div class="provider-icon-placeholder">
            {{ provider.display_name.charAt(0) }}
          </div>
          <div class="provider-info">
            <div class="provider-name">{{ provider.display_name }}</div>
            <div class="provider-desc">{{ provider.description }}</div>
          </div>
          <t-icon name="chevron-right" class="arrow-icon" />
        </div>
      </div>
    </t-dialog>

    <!-- 凭证配置弹窗 -->
    <CredentialForm
      v-model="showCredentialForm"
      :provider="selectedProvider"
      :credential="editingCredential"
      @save="handleCredentialSave"
      @test="handleCredentialTest"
    />

    <!-- 模型编辑弹窗 -->
    <ModelEditorDialog
      v-model:visible="showModelEditor"
      :model-type="currentModelType"
      :editing-model="editingModel"
      @save="handleModelSave"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { useI18n } from 'vue-i18n'
import ProviderCard from '@/components/model/ProviderCard.vue'
import CredentialForm from '@/components/model/CredentialForm.vue'
import ModelEditorDialog from '@/components/ModelEditorDialog.vue'
import { 
  listModels, 
  deleteModel, 
  createModel, 
  updateModel,
  listProviders,
  type ModelConfig,
  type ProviderDetail
} from '@/api/model'
import {
  listCredentials,
  createCredential,
  updateCredential,
  deleteCredential,
  testCredential,
  type ProviderCredential
} from '@/api/credential'

const { t } = useI18n()

// 状态
const activeTab = ref('providers')
const loading = ref(false)
const allModels = ref<ModelConfig[]>([])
const allProviders = ref<ProviderDetail[]>([])
const allCredentials = ref<ProviderCredential[]>([])

// 筛选
const filterType = ref('')
const filterProvider = ref('')
const searchKeyword = ref('')

// 弹窗状态
const showProviderSelector = ref(false)
const showCredentialForm = ref(false)
const showModelEditor = ref(false)
const selectedProvider = ref<ProviderDetail | null>(null)
const editingCredential = ref<ProviderCredential | null>(null)
const editingModel = ref<any>(null)
const currentModelType = ref<'chat' | 'embedding' | 'rerank' | 'vllm'>('chat')

// 计算属性
const configuredProviders = computed(() => {
  const credentialProviders = new Set(allCredentials.value.map(c => c.provider))
  return allProviders.value.filter(p => credentialProviders.has(p.name))
})

const availableProviders = computed(() => {
  const credentialProviders = new Set(allCredentials.value.map(c => c.provider))
  return allProviders.value.filter(p => !credentialProviders.has(p.name))
})

const filteredModels = computed(() => {
  let result = allModels.value
  
  if (filterType.value) {
    result = result.filter(m => m.type === filterType.value)
  }
  
  if (filterProvider.value) {
    result = result.filter(m => m.parameters?.provider === filterProvider.value)
  }
  
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(m => m.name.toLowerCase().includes(keyword))
  }
  
  return result
})

const modelTypeOptions = [
  { value: 'KnowledgeQA', label: '对话模型' },
  { value: 'Embedding', label: 'Embedding 模型' },
  { value: 'Rerank', label: 'Rerank 模型' },
  { value: 'VLLM', label: 'VLLM 模型' }
]

const providerOptions = computed(() => {
  return allProviders.value.map(p => ({
    value: p.name,
    label: p.display_name
  }))
})

// 方法
const loadData = async () => {
  loading.value = true
  try {
    const [models, providers, credentials] = await Promise.all([
      listModels(),
      listProviders(),
      listCredentials()
    ])
    allModels.value = models
    allProviders.value = providers
    allCredentials.value = credentials
  } catch (error: any) {
    MessagePlugin.error(error.message || '加载数据失败')
  } finally {
    loading.value = false
  }
}

const getCredentialStatus = (providerName: string) => {
  const credential = allCredentials.value.find(c => c.provider === providerName)
  return credential?.status || 'active'
}

const getModelCount = (providerName: string) => {
  return allModels.value.filter(m => m.parameters?.provider === providerName).length
}

const getModelTypeLabel = (type: string) => {
  const option = modelTypeOptions.find(o => o.value === type)
  return option?.label || type
}

const getProviderLabel = (providerName?: string) => {
  if (!providerName) return '未知'
  const provider = allProviders.value.find(p => p.name === providerName)
  return provider?.display_name || providerName
}

const selectProvider = (provider: ProviderDetail) => {
  selectedProvider.value = provider
  editingCredential.value = null
  showProviderSelector.value = false
  showCredentialForm.value = true
}

const openProviderManage = (provider: ProviderDetail) => {
  filterProvider.value = provider.name
  activeTab.value = 'models'
}

const openCredentialEdit = (provider: ProviderDetail) => {
  selectedProvider.value = provider
  const credential = allCredentials.value.find(c => c.provider === provider.name)
  editingCredential.value = credential || null
  showCredentialForm.value = true
}

const deleteProviderCredential = async (provider: ProviderDetail) => {
  const credential = allCredentials.value.find(c => c.provider === provider.name)
  if (!credential) return
  
  try {
    await deleteCredential(credential.id!)
    MessagePlugin.success('删除成功')
    loadData()
  } catch (error: any) {
    MessagePlugin.error(error.message || '删除失败')
  }
}

const testProviderCredential = async (provider: ProviderDetail) => {
  const credential = allCredentials.value.find(c => c.provider === provider.name)
  if (!credential) return
  
  try {
    const result = await testCredential(credential.id!)
    if (result.success) {
      MessagePlugin.success('连接测试成功')
    } else {
      MessagePlugin.error(result.message || '连接测试失败')
    }
    loadData()
  } catch (error: any) {
    MessagePlugin.error(error.message || '测试失败')
  }
}

const handleCredentialSave = async (data: any) => {
  try {
    if (data.id) {
      await updateCredential(data.id, data)
      MessagePlugin.success('更新成功')
    } else {
      await createCredential({
        ...data,
        provider: selectedProvider.value?.name
      })
      MessagePlugin.success('创建成功')
    }
    showCredentialForm.value = false
    loadData()
  } catch (error: any) {
    MessagePlugin.error(error.message || '保存失败')
  }
}

const handleCredentialTest = async (data: any) => {
  // 测试逻辑在 CredentialForm 组件内处理
}

const openAddModel = () => {
  editingModel.value = null
  showModelEditor.value = true
}

const editModel = (model: ModelConfig) => {
  // 转换为编辑器需要的格式
  editingModel.value = {
    id: model.id,
    name: model.name,
    source: model.source,
    modelName: model.name,
    baseUrl: model.parameters?.base_url || '',
    apiKey: model.parameters?.api_key || '',
    provider: model.parameters?.provider || '',
    dimension: model.parameters?.embedding_parameters?.dimension
  }
  
  // 设置模型类型
  switch (model.type) {
    case 'KnowledgeQA':
      currentModelType.value = 'chat'
      break
    case 'Embedding':
      currentModelType.value = 'embedding'
      break
    case 'Rerank':
      currentModelType.value = 'rerank'
      break
    case 'VLLM':
      currentModelType.value = 'vllm'
      break
  }
  
  showModelEditor.value = true
}

const handleModelSave = async (modelData: any) => {
  try {
    // 转换模型类型
    let modelType: ModelConfig['type'] = 'KnowledgeQA'
    switch (currentModelType.value) {
      case 'chat':
        modelType = 'KnowledgeQA'
        break
      case 'embedding':
        modelType = 'Embedding'
        break
      case 'rerank':
        modelType = 'Rerank'
        break
      case 'vllm':
        modelType = 'VLLM'
        break
    }
    
    const data: ModelConfig = {
      name: modelData.name || modelData.modelName,
      type: modelType,
      source: 'remote',
      parameters: {
        base_url: modelData.baseUrl,
        api_key: modelData.apiKey,
        provider: modelData.provider,
        embedding_parameters: modelData.dimension ? {
          dimension: modelData.dimension
        } : undefined
      }
    }
    
    if (modelData.id) {
      await updateModel(modelData.id, data)
      MessagePlugin.success('更新成功')
    } else {
      await createModel(data)
      MessagePlugin.success('创建成功')
    }
    
    showModelEditor.value = false
    loadData()
  } catch (error: any) {
    MessagePlugin.error(error.message || '保存失败')
  }
}

const handleDeleteModel = async (id: string) => {
  try {
    await deleteModel(id)
    MessagePlugin.success('删除成功')
    loadData()
  } catch (error: any) {
    MessagePlugin.error(error.message || '删除失败')
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="less">
.model-settings-v2 {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 24px;
  
  h2 {
    margin: 0 0 8px;
    font-size: 20px;
    font-weight: 600;
  }
  
  .page-description {
    margin: 0;
    color: var(--td-text-color-secondary);
    font-size: 14px;
  }
}

.settings-tabs {
  :deep(.t-tabs__nav) {
    margin-bottom: 24px;
  }
}

.tab-content {
  min-height: 400px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  
  h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 500;
  }
}

.provider-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  align-items: center;
  
  .filter-spacer {
    flex: 1;
  }
}

.model-list {
  border: 1px solid var(--td-border-level-1-color);
  border-radius: 8px;
  overflow: hidden;
}

.model-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid var(--td-border-level-1-color);
  
  &:last-child {
    border-bottom: none;
  }
  
  &:hover {
    background: var(--td-bg-color-container-hover);
  }
  
  &.is-builtin {
    background: var(--td-bg-color-secondarycontainer);
  }
}

.model-info {
  flex: 1;
}

.model-name {
  font-weight: 500;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.model-meta {
  display: flex;
  gap: 16px;
  font-size: 12px;
  color: var(--td-text-color-secondary);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.model-actions {
  display: flex;
  gap: 8px;
}

.empty-state {
  text-align: center;
  padding: 48px;
  color: var(--td-text-color-placeholder);
  
  .empty-icon {
    margin-bottom: 16px;
    opacity: 0.5;
  }
  
  p {
    margin: 0 0 16px;
  }
}

.provider-selector {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.provider-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border: 1px solid var(--td-border-level-1-color);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  
  &:hover {
    border-color: var(--td-brand-color);
    background: var(--td-bg-color-container-hover);
  }
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
  flex-shrink: 0;
}

.provider-info {
  flex: 1;
  min-width: 0;
}

.provider-name {
  font-weight: 500;
  margin-bottom: 2px;
}

.provider-desc {
  font-size: 12px;
  color: var(--td-text-color-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.arrow-icon {
  color: var(--td-text-color-placeholder);
}
</style>
