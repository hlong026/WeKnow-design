<template>
  <div class="aliyun-api-key-config">
    <t-form :data="formData" @submit="handleSubmit">
      <t-form-item label="阿里云 API Key" name="apiKey">
        <t-input
          v-model="formData.apiKey"
          type="password"
          placeholder="请输入阿里云 API Key (用于自媒体文案提取)"
          clearable
        />
        <div class="hint">
          用于调用 Coze workflow API 提取自媒体文案内容
        </div>
      </t-form-item>

      <t-form-item>
        <t-space>
          <t-button theme="primary" type="submit" :loading="saving">
            保存配置
          </t-button>
          <t-button theme="default" @click="handleCancel">
            取消
          </t-button>
        </t-space>
      </t-form-item>
    </t-form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { updateAliyunAPIKey } from '@/api/social-media'

const props = defineProps<{
  kbId: string
  currentApiKey?: string
}>()

const emit = defineEmits<{
  (e: 'success'): void
  (e: 'cancel'): void
}>()

const formData = ref({
  apiKey: props.currentApiKey || ''
})

const saving = ref(false)

onMounted(() => {
  if (props.currentApiKey) {
    formData.value.apiKey = props.currentApiKey
  }
})

const handleSubmit = async () => {
  if (!formData.value.apiKey.trim()) {
    MessagePlugin.warning('请输入阿里云 API Key')
    return
  }

  saving.value = true

  try {
    const response = await updateAliyunAPIKey(props.kbId, formData.value.apiKey)
    
    if (response.success) {
      MessagePlugin.success('配置保存成功')
      emit('success')
    } else {
      MessagePlugin.error(response.message || '配置保存失败')
    }
  } catch (error: any) {
    MessagePlugin.error(error.message || '配置保存失败')
  } finally {
    saving.value = false
  }
}

const handleCancel = () => {
  emit('cancel')
}
</script>

<style scoped lang="less">
.aliyun-api-key-config {
  padding: 20px;

  .hint {
    margin-top: 8px;
    color: #00000066;
    font-size: 12px;
    line-height: 20px;
  }
}
</style>
