<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="dialogVisible" class="model-editor-overlay" @click.self="handleCancel">
        <div class="model-editor-modal">
          <!-- 关闭按钮 -->
          <button class="close-btn" @click="handleCancel" :aria-label="$t('common.close')">
            <svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
              <path d="M15 5L5 15M5 5L15 15" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
            </svg>
          </button>

          <!-- 标题区域 -->
          <div class="modal-header">
            <h2 class="modal-title">{{ isEdit ? $t('model.editor.editTitle') : $t('model.editor.addTitle') }}</h2>
            <p class="modal-desc">{{ getModalDescription() }}</p>
          </div>

          <!-- 表单内容区域 -->
          <div class="modal-body">
            <t-form ref="formRef" :data="formData" :rules="rules" layout="vertical">
              <!-- 模型来源（隐藏，默认使用remote） -->
              <input type="hidden" v-model="formData.source" />

              <!-- 厂商选择器 -->
              <div class="form-item">
                <label class="form-label">{{ $t('model.editor.providerLabel') }}</label>
                <t-select 
                  v-model="formData.provider" 
                  :placeholder="$t('model.editor.providerPlaceholder')"
                  @change="handleProviderChange"
                >
                  <t-option 
                    v-for="opt in providerOptions" 
                    :key="opt.value" 
                    :value="opt.value" 
                    :label="opt.label"
                  >
                    <div class="provider-option">
                      <span class="provider-name">{{ opt.label }}</span>
                      <span class="provider-desc">{{ opt.description }}</span>
                    </div>
                  </t-option>
                </t-select>
              </div>

              <!-- 模型名称 -->
              <div class="form-item">
                <label class="form-label required">{{ $t('model.modelName') }}</label>
                <t-input 
                  v-model="formData.modelName" 
                  :placeholder="getModelNamePlaceholder()"
                />
              </div>

              <div class="form-item">
                <label class="form-label required">{{ $t('model.editor.baseUrlLabel') }}</label>
                <t-input 
                  v-model="formData.baseUrl" 
                  :placeholder="getBaseUrlPlaceholder()"
                />
              </div>

              <div class="form-item">
                <label class="form-label">{{ $t('model.editor.apiKeyOptional') }}</label>
                <t-input 
                  v-model="formData.apiKey" 
                  type="password"
                  autocomplete="off"
                  :placeholder="$t('model.editor.apiKeyPlaceholder')"
                />
              </div>

              <!-- Remote API 校验 -->
              <div class="form-item">
                <label class="form-label">{{ $t('model.editor.connectionTest') }}</label>
                <div class="api-test-section">
                  <t-button 
                    variant="outline" 
                    @click="checkRemoteAPI"
                    :loading="checking"
                    :disabled="!formData.modelName || !formData.baseUrl"
                  >
                    <template #icon>
                      <t-icon 
                        v-if="!checking && remoteChecked && remoteAvailable"
                        name="check-circle-filled" 
                        class="status-icon available"
                      />
                      <t-icon 
                        v-else-if="!checking && remoteChecked && !remoteAvailable"
                        name="close-circle-filled" 
                        class="status-icon unavailable"
                      />
                    </template>
                    {{ checking ? $t('model.editor.testing') : $t('model.editor.testConnection') }}
                  </t-button>
                  <span v-if="remoteChecked" :class="['test-message', remoteAvailable ? 'success' : 'error']">
                    {{ remoteMessage }}
                  </span>
                </div>
              </div>

              <!-- Embedding 专用：维度 -->
              <div v-if="modelType === 'embedding'" class="form-item">
                <label class="form-label">{{ $t('model.editor.dimensionLabel') }}</label>
                <div class="dimension-control">
                  <t-input 
                    v-model.number="formData.dimension" 
                    type="number"
                    :min="128"
                    :max="4096"
                    :placeholder="$t('model.editor.dimensionPlaceholder')"
                  />
                </div>
                <p v-if="dimensionChecked && dimensionMessage" class="dimension-hint" :class="{ success: dimensionSuccess }">
                  {{ dimensionMessage }}
                </p>
              </div>
            </t-form>
          </div>

          <!-- 底部按钮区域 -->
          <div class="modal-footer">
            <t-button theme="default" variant="outline" @click="handleCancel">
              {{ $t('common.cancel') }}
            </t-button>
            <t-button theme="primary" @click="handleConfirm" :loading="saving">
              {{ $t('common.save') }}
            </t-button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { checkRemoteModel, testEmbeddingModel, checkRerankModel, listModelProviders, type ModelProviderOption } from '@/api/initialization'
import { useI18n } from 'vue-i18n'

interface ModelFormData {
  id: string
  name: string
  source: 'local' | 'remote'
  provider?: string // Provider identifier: openai, aliyun, zhipu, generic, etc.
  modelName: string
  baseUrl?: string
  apiKey?: string
  dimension?: number
  interfaceType?: 'ollama' | 'openai'
  isDefault: boolean
}

interface Props {
  visible: boolean
  modelType: 'chat' | 'embedding' | 'rerank' | 'vllm'
  modelData?: ModelFormData | null
}

const { t } = useI18n()

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  modelData: null
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
  'confirm': [data: ModelFormData]
}>()

// API 返回的 Provider 列表
const apiProviderOptions = ref<ModelProviderOption[]>([])
const loadingProviders = ref(false)

// 硬编码的后备 Provider 配置 (当 API 不可用时使用)
const fallbackProviderOptions = computed(() => [
  { 
    value: 'openai', 
    label: t('model.editor.providers.openai.label'), 
    defaultUrls: {
      chat: 'https://api.openai.com/v1',
      embedding: 'https://api.openai.com/v1',
      rerank: 'https://api.openai.com/v1',
      vllm: 'https://api.openai.com/v1'
    },
    description: t('model.editor.providers.openai.description'),
    modelTypes: ['chat', 'embedding', 'vllm']
  },
  { 
    value: 'aliyun', 
    label: t('model.editor.providers.aliyun.label'), 
    defaultUrls: {
      chat: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
      embedding: 'https://dashscope.aliyuncs.com/compatible-mode/v1',
      rerank: 'https://dashscope.aliyuncs.com/api/v1/services/rerank/text-rerank/text-rerank',
      vllm: 'https://dashscope.aliyuncs.com/compatible-mode/v1'
    },
    description: t('model.editor.providers.aliyun.description'),
    modelTypes: ['chat', 'embedding', 'rerank', 'vllm']
  },
  { 
    value: 'zhipu', 
    label: t('model.editor.providers.zhipu.label'), 
    defaultUrls: {
      chat: 'https://open.bigmodel.cn/api/paas/v4',
      embedding: 'https://open.bigmodel.cn/api/paas/v4/embeddings',
      vllm: 'https://open.bigmodel.cn/api/paas/v4'
    },
    description: t('model.editor.providers.zhipu.description'),
    modelTypes: ['chat', 'embedding', 'vllm']
  },
  { 
    value: 'openrouter', 
    label: t('model.editor.providers.openrouter.label'), 
    defaultUrls: { chat: 'https://openrouter.ai/api/v1' },
    description: t('model.editor.providers.openrouter.description'),
    modelTypes: ['chat']
  },
  { 
    value: 'siliconflow', 
    label: t('model.editor.providers.siliconflow.label'), 
    defaultUrls: {
      chat: 'https://api.siliconflow.cn/v1',
      embedding: 'https://api.siliconflow.cn/v1',
      rerank: 'https://api.siliconflow.cn/v1'
    },
    description: t('model.editor.providers.siliconflow.description'),
    modelTypes: ['chat', 'embedding', 'rerank']
  },
  { 
    value: 'jina', 
    label: t('model.editor.providers.jina.label'), 
    defaultUrls: {
      embedding: 'https://api.jina.ai/v1',
      rerank: 'https://api.jina.ai/v1'
    },
    description: t('model.editor.providers.jina.description'),
    modelTypes: ['embedding', 'rerank']
  },
  { 
    value: 'generic', 
    label: t('model.editor.providers.generic.label'), 
    defaultUrls: {},
    description: t('model.editor.providers.generic.description'),
    modelTypes: ['chat', 'embedding', 'rerank', 'vllm']
  },
])

// 从 API 获取 Provider 列表
const loadProviders = async () => {
  console.log('[ModelEditor] Loading providers for type:', props.modelType)
  loadingProviders.value = true
  try {
    const providers = await listModelProviders(props.modelType)
    console.log('[ModelEditor] API returned providers:', providers?.length || 0, providers)
    if (providers && providers.length > 0) {
      apiProviderOptions.value = providers
      console.log('[ModelEditor] Set API providers:', apiProviderOptions.value.length)
    } else {
      console.warn('[ModelEditor] API returned empty providers, will use fallback')
    }
  } catch (error) {
    console.error('[ModelEditor] Failed to load providers from API, using fallback', error)
  } finally {
    loadingProviders.value = false
  }
}

// 根据当前模型类型过滤的 Provider 列表
const providerOptions = computed(() => {
  // 优先使用 API 返回的数据
  if (apiProviderOptions.value.length > 0) {
    return apiProviderOptions.value
  }
  // 回退到硬编码值，按 modelTypes 过滤
  const filtered = fallbackProviderOptions.value.filter(p => 
    p.modelTypes.includes(props.modelType)
  )
  console.log('[ModelEditor] Using fallback providers for type', props.modelType, ':', filtered.length, filtered)
  return filtered
})

const dialogVisible = computed({
  get: () => props.visible,
  set: (val) => emit('update:visible', val)
})

const isEdit = computed(() => !!props.modelData)

const formRef = ref()
const saving = ref(false)
const checking = ref(false)
const remoteChecked = ref(false)
const remoteAvailable = ref(false)
const remoteMessage = ref('')
const dimensionChecked = ref(false)
const dimensionSuccess = ref(false)
const dimensionMessage = ref('')

const formData = ref<ModelFormData>({
  id: '',
  name: '',
  source: 'remote',
  provider: 'openai',
  modelName: '',
  baseUrl: '',
  apiKey: '',
  dimension: undefined,
  interfaceType: undefined,
  isDefault: false
})

const rules = computed(() => {
  const baseRules: any = {
    modelName: [
      { required: true, message: t('model.editor.validation.modelNameRequired') },
      { 
        validator: (val: string) => {
          if (!val || !val.trim()) {
            return { result: false, message: t('model.editor.validation.modelNameEmpty') }
          }
          if (val.trim().length > 100) {
            return { result: false, message: t('model.editor.validation.modelNameMax') }
          }
          return { result: true }
        },
        trigger: 'blur'
      }
    ],
    baseUrl: [
      { 
        required: true, 
        message: t('model.editor.validation.baseUrlRequired'),
        trigger: 'blur'
      },
      {
        validator: (val: string) => {
          if (!val || !val.trim()) {
            return { result: false, message: t('model.editor.validation.baseUrlEmpty') }
          }
          // 简单的 URL 格式校验
          try {
            new URL(val.trim())
            return { result: true }
          } catch {
            return { result: false, message: t('model.editor.validation.baseUrlInvalid') }
          }
        },
        trigger: 'blur'
      }
    ]
  }
  
  // Embedding 模型需要额外验证维度
  if (props.modelType === 'embedding') {
    baseRules.dimension = [
      { 
        required: true, 
        message: t('model.editor.validation.dimensionRequired') || 'Dimension is required',
        trigger: 'blur'
      },
      {
        validator: (val: number) => {
          if (!val || val < 128 || val > 4096) {
            return { 
              result: false, 
              message: t('model.editor.validation.dimensionRange') || 'Dimension must be between 128 and 4096'
            }
          }
          return { result: true }
        },
        trigger: 'blur'
      }
    ]
  }
  
  return baseRules
})

// 获取弹窗描述文字
const getModalDescription = () => {
  const key = `model.editor.description.${props.modelType}` as const
  return t(key) || t('model.editor.description.default')
}

// 获取模型名称占位符
const getModelNamePlaceholder = () => {
  if (props.modelType === 'vllm') {
    return formData.value.source === 'local'
      ? t('model.editor.modelNamePlaceholder.localVllm')
      : t('model.editor.modelNamePlaceholder.remoteVllm')
  }
  return formData.value.source === 'local'
    ? t('model.editor.modelNamePlaceholder.local')
    : t('model.editor.modelNamePlaceholder.remote')
}

const getBaseUrlPlaceholder = () => {
  return props.modelType === 'vllm'
    ? t('model.editor.baseUrlPlaceholderVllm')
    : t('model.editor.baseUrlPlaceholder')
}

// 监听 visible 变化，初始化表单
watch(() => props.visible, (val) => {
  console.log('[ModelEditor] Dialog visible changed:', val, 'Model type:', props.modelType)
  if (val) {
    // 锁定背景滚动
    document.body.style.overflow = 'hidden'

    // 从 API 加载 Model Provider 列表
    loadProviders()

    if (props.modelData) {
      console.log('[ModelEditor] Editing existing model:', props.modelData)
      formData.value = { ...props.modelData }
    } else {
      console.log('[ModelEditor] Creating new model')
      resetForm()
    }

    // ReRank 模型强制使用 remote 来源
    if (props.modelType === 'rerank') {
      formData.value.source = 'remote'
    }
    
    console.log('[ModelEditor] Form data initialized:', formData.value)
  } else {
    // 恢复背景滚动
    document.body.style.overflow = ''
  }
})

// 重置表单
const resetForm = () => {
  console.log('[ModelEditor] Resetting form')
  formData.value = {
    id: generateId(),
    name: '',
    source: 'remote',
    provider: 'generic',
    modelName: '',
    baseUrl: '',
    apiKey: '',
    dimension: undefined,
    interfaceType: undefined,
    isDefault: false
  }
  remoteChecked.value = false
  remoteAvailable.value = false
  remoteMessage.value = ''
  dimensionChecked.value = false
  dimensionSuccess.value = false
  dimensionMessage.value = ''
  console.log('[ModelEditor] Form reset complete')
}

// 处理厂商选择变化 (自动填充默认 URL)
const handleProviderChange = (value: string) => {
  console.log('[ModelEditor] Provider changed:', value)
  const provider = providerOptions.value.find(opt => opt.value === value)
  console.log('[ModelEditor] Found provider config:', provider)
  if (provider && provider.defaultUrls) {
    // 根据当前模型类型获取对应的默认 URL
    const defaultUrl = provider.defaultUrls[props.modelType]
    console.log('[ModelEditor] Default URL for', props.modelType, ':', defaultUrl)
    if (defaultUrl) {
      formData.value.baseUrl = defaultUrl
      console.log('[ModelEditor] Set baseUrl to:', defaultUrl)
    }
    // 重置校验状态
    remoteChecked.value = false
    remoteAvailable.value = false
    remoteMessage.value = ''
  }
}

// 监听来源变化，重置校验状态（已合并到下面的 watch）

// 生成唯一ID
const generateId = () => {
  return `model_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
}

// 检查 Remote API 连接（根据模型类型调用不同的接口）
const checkRemoteAPI = async () => {
  if (!formData.value.modelName || !formData.value.baseUrl) {
    MessagePlugin.warning(t('model.editor.fillModelAndUrl'))
    return
  }
  
  checking.value = true
  remoteChecked.value = false
  remoteMessage.value = ''
  
  try {
    let result: any
    
    // 根据模型类型调用不同的校验接口
    switch (props.modelType) {
      case 'chat':
        // 对话模型（KnowledgeQA）
        result = await checkRemoteModel({
          modelName: formData.value.modelName,
          baseUrl: formData.value.baseUrl,
          apiKey: formData.value.apiKey || ''
        })
        break
        
      case 'embedding':
        // Embedding 模型
        result = await testEmbeddingModel({
          source: 'remote',
          modelName: formData.value.modelName,
          baseUrl: formData.value.baseUrl,
          apiKey: formData.value.apiKey || '',
          dimension: formData.value.dimension,
          provider: formData.value.provider
        })
        // 如果测试成功且返回了维度，自动填充
        if (result.available && result.dimension) {
          formData.value.dimension = result.dimension
        MessagePlugin.info(t('model.editor.remoteDimensionDetected', { value: result.dimension }))
        }
        break
        
      case 'rerank':
        // Rerank 模型
        result = await checkRerankModel({
          modelName: formData.value.modelName,
          baseUrl: formData.value.baseUrl,
          apiKey: formData.value.apiKey || ''
        })
        break
        
      case 'vllm':
        // VLLM 模型（多模态）
        // VLLM 使用 checkRemoteModel 进行基础连接测试
        result = await checkRemoteModel({
          modelName: formData.value.modelName,
          baseUrl: formData.value.baseUrl,
          apiKey: formData.value.apiKey || ''
        })
        break
        
      default:
        MessagePlugin.error(t('model.editor.unsupportedModelType'))
        return
    }
    
    remoteChecked.value = true
    remoteAvailable.value = result.available || false
    remoteMessage.value = result.message || (result.available ? t('model.editor.connectionSuccess') : t('model.editor.connectionFailed'))
    
    if (result.available) {
      MessagePlugin.success(remoteMessage.value)
    } else {
      MessagePlugin.error(remoteMessage.value)
    }
  } catch (error: any) {
    console.error('Remote API 校验失败:', error)
    remoteChecked.value = true
    remoteAvailable.value = false
    remoteMessage.value = error.message || t('model.editor.connectionConfigError')
    MessagePlugin.error(remoteMessage.value)
  } finally {
    checking.value = false
  }
}

// 确认保存
const handleConfirm = async () => {
  console.log('[ModelEditor] handleConfirm called, formData:', formData.value)
  try {
    // 手动校验必填字段
    if (!formData.value.modelName || !formData.value.modelName.trim()) {
      console.log('[ModelEditor] Validation failed: modelName is empty')
      MessagePlugin.warning(t('model.editor.validation.modelNameRequired'))
      return
    }
    
    if (formData.value.modelName.trim().length > 100) {
      console.log('[ModelEditor] Validation failed: modelName too long')
      MessagePlugin.warning(t('model.editor.validation.modelNameMax'))
      return
    }
    
    // 如果是 remote 类型，必须填写 baseUrl 和 provider
    if (formData.value.source === 'remote') {
      if (!formData.value.baseUrl || !formData.value.baseUrl.trim()) {
        console.log('[ModelEditor] Validation failed: baseUrl is empty')
        MessagePlugin.warning(t('model.editor.remoteBaseUrlRequired'))
        return
      }
      
      // 校验 Base URL 格式
      try {
        new URL(formData.value.baseUrl.trim())
      } catch {
        console.log('[ModelEditor] Validation failed: baseUrl is invalid')
        MessagePlugin.warning(t('model.editor.validation.baseUrlInvalid'))
        return
      }

      // 校验 provider
      if (!formData.value.provider || !formData.value.provider.trim()) {
        console.log('[ModelEditor] Validation failed: provider is empty')
        MessagePlugin.warning(t('model.editor.validation.providerRequired') || 'Please select a provider')
        return
      }
    }

    // Embedding 模型必须填写维度
    if (props.modelType === 'embedding') {
      if (!formData.value.dimension || formData.value.dimension < 128 || formData.value.dimension > 4096) {
        console.log('[ModelEditor] Validation failed: dimension is invalid')
        MessagePlugin.warning(t('model.editor.validation.dimensionRange') || 'Dimension must be between 128 and 4096')
        return
      }
    }
    
    saving.value = true
    console.log('[ModelEditor] Starting form validation...')
    
    // 执行表单验证（如果验证失败会抛出异常）
    try {
      await formRef.value?.validate()
      console.log('[ModelEditor] Form validation passed')
    } catch (validationError) {
      console.error('[ModelEditor] Form validation failed:', validationError)
      MessagePlugin.warning(t('model.editor.validation.failed') || 'Please check the form fields')
      saving.value = false
      return
    }
    
    // 如果是新增且没有 id，生成一个
    if (!formData.value.id) {
      formData.value.id = generateId()
    }
    
    console.log('[ModelEditor] Emitting confirm event with data:', formData.value)
    emit('confirm', { ...formData.value })
    dialogVisible.value = false
    // 移除此处的成功提示，由父组件统一处理
  } catch (error) {
    console.error('[ModelEditor] Save failed:', error)
    MessagePlugin.error(t('model.editor.saveFailed') || 'Failed to save model')
  } finally {
    saving.value = false
  }
}

// 监听模型名称变化，清理维度检测状态
watch(() => formData.value.modelName, () => {
  dimensionChecked.value = false
  dimensionSuccess.value = false
  dimensionMessage.value = ''
})

// 监听 providerOptions 变化,确保有默认选择
watch(() => providerOptions.value, (options) => {
  console.log('[ModelEditor] Provider options changed:', options.length)
  if (options.length > 0) {
    // 如果当前没有选择 provider 或选择的 provider 不在列表中
    const currentProviderExists = options.some(opt => opt.value === formData.value.provider)
    if (!formData.value.provider || !currentProviderExists) {
      formData.value.provider = options[0].value
      console.log('[ModelEditor] Set default provider:', options[0].value)
      // 触发 provider 变化处理
      handleProviderChange(options[0].value)
    }
  }
}, { immediate: true })

// 取消
const handleCancel = () => {
  dialogVisible.value = false
}
</script>

<style lang="less" scoped>
// 遮罩层
.model-editor-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1200;
  backdrop-filter: blur(4px);
  overflow: hidden;
  padding: 20px;
}

// 弹窗主体
.model-editor-modal {
  position: relative;
  width: 100%;
  max-width: 560px;
  max-height: 90vh;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 6px 28px rgba(15, 23, 42, 0.08);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

// 关闭按钮
.close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 32px;
  height: 32px;
  border: none;
  background: transparent;
  border-radius: 6px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666666;
  transition: all 0.15s ease;
  z-index: 10;

  &:hover {
    background: #f5f5f5;
    color: #333333;
  }
}

// 标题区域
.modal-header {
  padding: 24px 24px 16px;
  border-bottom: 1px solid #e5e7eb;
  flex-shrink: 0;
}

.modal-title {
  margin: 0 0 6px 0;
  font-size: 18px;
  font-weight: 600;
  color: #333333;
}

.modal-desc {
  margin: 0;
  font-size: 13px;
  color: #666666;
  line-height: 1.5;
}

// 内容区域
.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  background: #ffffff;

  // 自定义滚动条
  &::-webkit-scrollbar {
    width: 6px;
  }

  &::-webkit-scrollbar-track {
    background: #f5f5f5;
    border-radius: 3px;
  }

  &::-webkit-scrollbar-thumb {
    background: #d0d0d0;
    border-radius: 3px;
    transition: background 0.15s;

    &:hover {
      background: #b0b0b0;
    }
  }

}

// 表单项样式
.form-item {
  margin-bottom: 20px;

  &:last-child {
    margin-bottom: 0;
  }
}

.form-label {
  display: block;
  margin-bottom: 8px;
  font-size: 14px;
  font-weight: 500;
  color: #333333;

  &.required::after {
    content: '*';
    color: #f56c6c;
    margin-left: 4px;
    font-weight: 600;
  }
}

// 输入框样式
:deep(.t-input),
:deep(.t-select),
:deep(.t-textarea),
:deep(.t-input-number) {
  width: 100%;
  font-size: 13px;

  .t-input__inner,
  .t-input__wrap,
  input,
  textarea {
    font-size: 13px;
    border-radius: 6px;
    border-color: #d9d9d9;
    transition: all 0.15s ease;
  }

  &:hover .t-input__inner,
  &:hover .t-input__wrap,
  &:hover input,
  &:hover textarea {
    border-color: #b3b3b3;
  }

  &.t-is-focused .t-input__inner,
  &.t-is-focused .t-input__wrap,
  &.t-is-focused input,
  &.t-is-focused textarea {
    border-color: #07C05F;
    box-shadow: 0 0 0 2px rgba(7, 192, 95, 0.1);
  }
}

// 厂商选择器样式
.provider-option {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 4px 0;

  .provider-name {
    font-size: 14px;
    font-weight: 500;
    color: #333333;
  }

  .provider-desc {
    font-size: 12px;
    color: #999999;
  }
}

// 单选按钮组
:deep(.t-radio-group) {
  display: flex;
  gap: 24px;

  .t-radio {
    margin-right: 0;
    font-size: 13px;

    &:hover {
      .t-radio__label {
        color: #07C05F;
      }
    }
  }

  .t-radio__label {
    font-size: 13px;
    color: #333333;
    transition: color 0.15s ease;
  }

  .t-radio__input:checked + .t-radio__label {
    color: #07C05F;
    font-weight: 500;
  }
}

// 复选框
:deep(.t-checkbox) {
  font-size: 13px;

  .t-checkbox__label {
    font-size: 13px;
    color: #333333;
  }
}

// 底部按钮区域
.modal-footer {
  padding: 16px 24px;
  border-top: 1px solid #e5e7eb;
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  flex-shrink: 0;
  background: #fafafa;

  :deep(.t-button) {
    min-width: 80px;
    height: 36px;
    font-weight: 500;
    font-size: 14px;
    border-radius: 6px;
    transition: all 0.15s ease;

    &.t-button--theme-primary {
      background: #07C05F;
      border-color: #07C05F;

      &:hover {
        background: #06b04d;
        border-color: #06b04d;
      }

      &:active {
        background: #059642;
        border-color: #059642;
      }
    }

    &.t-button--variant-outline {
      color: #666666;
      border-color: #d9d9d9;

      &:hover {
        border-color: #07C05F;
        color: #07C05F;
        background: rgba(7, 192, 95, 0.04);
      }
    }
  }
}

// 过渡动画
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;

  .model-editor-modal {
    transition: transform 0.2s ease, opacity 0.2s ease;
  }
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;

  .model-editor-modal {
    transform: scale(0.95);
    opacity: 0;
  }
}

// API 测试区域
.api-test-section {
  display: flex;
  align-items: center;
  gap: 12px;

  .test-message {
    font-size: 13px;
    line-height: 1.5;
    flex: 1;

    &.success {
      color: #059669;
    }

    &.error {
      color: #f56c6c;
    }
  }

  :deep(.t-button) {
    min-width: 88px;
    height: 32px;
    font-size: 13px;
    border-radius: 6px;
    flex-shrink: 0;
  }

  .status-icon {
    font-size: 16px;
    flex-shrink: 0;

    &.available {
      color: #07C05F;
    }

    &.unavailable {
      color: #f56c6c;
    }
  }
}

.model-select-row {
  display: flex;
  align-items: center;
  gap: 8px;

  .t-select {
    flex: 1;
  }

  :deep(.t-button) {
    height: 32px;
    font-size: 13px;
    border-radius: 6px;
    flex-shrink: 0;
  }
}

.refresh-btn {
  margin-top: 0;
  font-size: 13px;
  color: #666666;
  flex-shrink: 0;

  &:hover {
    color: #07C05F;
    background: rgba(7, 192, 95, 0.04);
  }
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

// 维度控制样式
.dimension-control {
  display: flex;
  align-items: center;
  gap: 8px;

  :deep(.t-input) {
    flex: 1;
  }
}

.dimension-check-btn {
  flex-shrink: 0;
  font-size: 12px;
}

.dimension-hint {
  margin: 8px 0 0 0;
  font-size: 13px;
  line-height: 1.5;
  color: #e34d59;

  &.success {
    color: #07C05F;
  }
}

</style>

