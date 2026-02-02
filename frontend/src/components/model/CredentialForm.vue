<template>
  <t-dialog
    v-model:visible="visible"
    :header="isEdit ? `编辑 ${providerDisplayName} 凭证` : `配置 ${providerDisplayName}`"
    :confirm-btn="{ content: '保存', loading: saving }"
    :cancel-btn="{ content: '取消' }"
    width="520px"
    @confirm="handleSave"
    @close="handleClose"
  >
    <t-form
      ref="formRef"
      :data="formData"
      :rules="formRules"
      label-width="100px"
      class="credential-form"
    >
      <!-- API Key 申请引导 -->
      <div v-if="helpUrl" class="api-key-guide">
        <t-icon name="info-circle" />
        <span>{{ helpText }}</span>
        <a :href="helpUrl" target="_blank" rel="noopener noreferrer" class="guide-link">
          前往申请 <t-icon name="jump" size="14px" />
        </a>
      </div>

      <!-- 凭证名称 -->
      <t-form-item label="凭证名称" name="name">
        <t-input
          v-model="formData.name"
          placeholder="为此凭证起个名字，便于识别"
        />
      </t-form-item>
      
      <!-- 动态认证字段 -->
      <t-form-item
        v-for="field in authFields"
        :key="field.key"
        :label="field.label"
        :name="`credentials.${field.key}`"
      >
        <t-input
          v-if="field.type === 'password'"
          v-model="formData.credentials[field.key]"
          type="password"
          :placeholder="field.placeholder"
        />
        <t-input
          v-else
          v-model="formData.credentials[field.key]"
          :placeholder="field.placeholder"
        />
        <template #help v-if="field.helpText || field.help_text">
          <span class="field-help">{{ field.helpText || field.help_text }}</span>
        </template>
      </t-form-item>
      
      <!-- 自定义端点 -->
      <t-form-item label="API 端点" name="baseUrl">
        <t-input
          v-model="formData.baseUrl"
          :placeholder="defaultEndpoint || '留空使用默认端点'"
        />
        <template #help>
          <span class="field-help">自定义 API 端点，留空使用官方默认地址</span>
        </template>
      </t-form-item>
      
      <!-- 配额设置 -->
      <t-collapse class="quota-collapse">
        <t-collapse-panel header="配额限制 (可选)" value="quota">
          <t-form-item label="每日限制" name="quotaDailyLimit">
            <t-input-number
              v-model="formData.quotaConfig.dailyLimit"
              :min="0"
              placeholder="不限制"
              suffix="次"
            />
          </t-form-item>
          <t-form-item label="每月限制" name="quotaMonthlyLimit">
            <t-input-number
              v-model="formData.quotaConfig.monthlyLimit"
              :min="0"
              placeholder="不限制"
              suffix="次"
            />
          </t-form-item>
        </t-collapse-panel>
      </t-collapse>
    </t-form>
    
    <template #footer>
      <div class="dialog-footer">
        <t-button
          theme="default"
          variant="outline"
          :loading="testing"
          @click="handleTest"
        >
          <template #icon><t-icon name="link" /></template>
          测试连接
        </t-button>
        <div class="footer-right">
          <t-button theme="default" @click="handleClose">取消</t-button>
          <t-button theme="primary" :loading="saving" @click="handleSave">保存</t-button>
        </div>
      </div>
    </template>
  </t-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'

interface AuthField {
  key: string
  label: string
  type: string
  required: boolean
  placeholder: string
  helpText?: string
  help_text?: string
}

interface Provider {
  name: string
  displayName?: string
  display_name?: string
  authConfig?: {
    fields: AuthField[]
    helpUrl?: string
    helpText?: string
  }
  auth_config?: {
    fields: AuthField[]
    help_url?: string
    help_text?: string
  }
  endpoints?: Record<string, string>
}

interface Credential {
  id?: string
  name: string
  credentials: Record<string, string>
  baseUrl: string
  quotaConfig: {
    dailyLimit?: number
    monthlyLimit?: number
  }
}

const props = defineProps<{
  modelValue: boolean
  provider?: Provider
  credential?: Credential
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'save': [data: Credential]
  'test': [data: Credential]
}>()

const visible = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val)
})

// 兼容 displayName 和 display_name
const providerDisplayName = computed(() => {
  return props.provider?.displayName || props.provider?.display_name || props.provider?.name || ''
})

const isEdit = computed(() => !!props.credential?.id)

const formRef = ref()
const saving = ref(false)
const testing = ref(false)

const formData = reactive<Credential>({
  name: '',
  credentials: {},
  baseUrl: '',
  quotaConfig: {}
})

// 初始化 credentials 对象，确保所有字段都有响应式属性
const initCredentials = () => {
  const authConfig = props.provider?.authConfig || props.provider?.auth_config
  const fields = authConfig?.fields || [
    { key: 'api_key', label: 'API Key', type: 'password', required: true, placeholder: '请输入 API Key' }
  ]
  const newCredentials: Record<string, string> = {}
  fields.forEach(field => {
    newCredentials[field.key] = formData.credentials[field.key] || ''
  })
  formData.credentials = newCredentials
}

// 监听 credential 变化，填充表单
watch(() => props.credential, (val) => {
  if (val) {
    formData.name = val.name
    formData.credentials = { ...val.credentials }
    formData.baseUrl = val.baseUrl
    formData.quotaConfig = { ...val.quotaConfig }
  } else {
    // 重置表单
    formData.name = providerDisplayName.value + ' 默认凭证'
    formData.credentials = {}
    formData.baseUrl = ''
    formData.quotaConfig = {}
  }
  // 初始化 credentials 对象
  initCredentials()
}, { immediate: true })

// 监听 provider 变化，重新初始化 credentials
watch(() => props.provider, () => {
  initCredentials()
}, { immediate: true })

const authFields = computed(() => {
  // 兼容 authConfig 和 auth_config
  const authConfig = props.provider?.authConfig || props.provider?.auth_config
  return authConfig?.fields || [
    { key: 'api_key', label: 'API Key', type: 'password', required: true, placeholder: '请输入 API Key' }
  ]
})

// 获取帮助链接
const helpUrl = computed(() => {
  const authConfig = props.provider?.authConfig || props.provider?.auth_config
  return authConfig?.helpUrl || authConfig?.help_url || ''
})

// 获取帮助文本
const helpText = computed(() => {
  const authConfig = props.provider?.authConfig || props.provider?.auth_config
  return authConfig?.helpText || authConfig?.help_text || '获取 API Key'
})

const defaultEndpoint = computed(() => {
  return props.provider?.endpoints?.['KnowledgeQA'] || ''
})

const formRules = computed(() => {
  const rules: Record<string, any[]> = {
    name: [{ required: true, message: '请输入凭证名称' }]
  }
  authFields.value.forEach(field => {
    if (field.required) {
      rules[`credentials.${field.key}`] = [{ required: true, message: `请输入 ${field.label}` }]
    }
  })
  return rules
})

const handleTest = async () => {
  testing.value = true
  try {
    emit('test', { ...formData })
    MessagePlugin.success('连接测试成功')
  } catch (error: any) {
    MessagePlugin.error(error.message || '连接测试失败')
  } finally {
    testing.value = false
  }
}

const handleSave = async () => {
  const valid = await formRef.value?.validate()
  if (valid !== true) return
  
  saving.value = true
  try {
    emit('save', { ...formData, id: props.credential?.id })
    visible.value = false
  } catch (error: any) {
    MessagePlugin.error(error.message || '保存失败')
  } finally {
    saving.value = false
  }
}

const handleClose = () => {
  visible.value = false
}
</script>

<style scoped lang="less">
.credential-form {
  padding: 16px 0;
}

.api-key-guide {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  margin-bottom: 16px;
  background: var(--td-brand-color-light);
  border-radius: 6px;
  font-size: 13px;
  color: var(--td-text-color-secondary);
  
  .t-icon {
    color: var(--td-brand-color);
    flex-shrink: 0;
  }
  
  .guide-link {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    color: var(--td-brand-color);
    text-decoration: none;
    margin-left: auto;
    white-space: nowrap;
    
    &:hover {
      text-decoration: underline;
    }
  }
}

.field-help {
  font-size: 12px;
  color: var(--td-text-color-placeholder);
}

.quota-collapse {
  margin-top: 16px;
  
  :deep(.t-collapse-panel__header) {
    padding: 12px 0;
  }
}

.dialog-footer {
  display: flex;
  justify-content: space-between;
  width: 100%;
}

.footer-right {
  display: flex;
  gap: 8px;
}
</style>
