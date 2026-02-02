<template>
  <t-dialog
    v-model:visible="visible"
    :header="$t('socialMedia.importTitle')"
    :confirm-btn="$t('common.import')"
    :cancel-btn="$t('common.cancel')"
    width="600px"
    @confirm="handleImport"
    @close="handleClose"
  >
    <div class="social-media-import-dialog">
      <div class="platform-selector">
        <label class="form-label">{{ $t('socialMedia.selectPlatform') }}</label>
        <t-select
          v-model="selectedPlatform"
          :options="platformOptions"
          :placeholder="$t('socialMedia.selectPlatformPlaceholder')"
          style="width: 100%"
        />
      </div>

      <div class="url-input">
        <label class="form-label">{{ $t('socialMedia.videoUrl') }}</label>
        <t-textarea
          v-model="videoUrl"
          :placeholder="$t('socialMedia.videoUrlPlaceholder')"
          :autosize="{ minRows: 3, maxRows: 6 }"
        />
        <div class="input-hint">{{ $t('socialMedia.urlHint') }}</div>
      </div>

      <div v-if="importing" class="importing-status">
        <t-loading size="small" />
        <span>{{ $t('socialMedia.importing') }}</span>
      </div>
    </div>
  </t-dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { useI18n } from 'vue-i18n'
import { extractSocialMediaContent } from '@/api/social-media'

const { t } = useI18n()

const props = defineProps<{
  modelValue: boolean
  kbId: string
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const visible = ref(props.modelValue)
const selectedPlatform = ref('')
const videoUrl = ref('')
const importing = ref(false)

const platformOptions = [
  { label: t('socialMedia.platforms.xiaohongshu'), value: 'xiaohongshu' },
  { label: t('socialMedia.platforms.douyin'), value: 'douyin' },
]

watch(() => props.modelValue, (val) => {
  visible.value = val
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleImport = async () => {
  if (!selectedPlatform.value) {
    MessagePlugin.warning(t('socialMedia.selectPlatformFirst'))
    return
  }

  if (!videoUrl.value.trim()) {
    MessagePlugin.warning(t('socialMedia.enterUrlFirst'))
    return
  }

  importing.value = true

  try {
    const response = await extractSocialMediaContent({
      platform: selectedPlatform.value,
      videoUrl: videoUrl.value,
      kbId: props.kbId
    })

    if (response.success) {
      MessagePlugin.success(t('socialMedia.importSuccess'))
      emit('success')
      handleClose()
    } else {
      MessagePlugin.error(response.message || t('socialMedia.importFailed'))
    }
  } catch (error: any) {
    MessagePlugin.error(error.message || t('socialMedia.importFailed'))
  } finally {
    importing.value = false
  }
}

const handleClose = () => {
  visible.value = false
  selectedPlatform.value = ''
  videoUrl.value = ''
  importing.value = false
}

</script>

<style scoped lang="less">
.social-media-import-dialog {
  padding: 20px 0;

  .platform-selector,
  .url-input {
    margin-bottom: 24px;
  }

  .form-label {
    display: block;
    margin-bottom: 8px;
    color: #000000e6;
    font-size: 14px;
    font-weight: 500;
  }

  .input-hint {
    margin-top: 8px;
    color: #00000066;
    font-size: 12px;
    line-height: 20px;
  }

  .importing-status {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px;
    background: #f7f9fc;
    border-radius: 6px;
    color: #00000099;
    font-size: 14px;
  }
}
</style>
