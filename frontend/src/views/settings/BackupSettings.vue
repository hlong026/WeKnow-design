<template>
  <div class="backup-settings">
    <div class="section-header">
      <h3 class="section-title">{{ $t('settings.backup.title') }}</h3>
      <p class="section-desc">{{ $t('settings.backup.description') }}</p>
    </div>

    <!-- Export Section -->
    <div class="backup-section">
      <div class="section-label">
        <t-icon name="download" class="section-icon" />
        <span>{{ $t('settings.backup.export') }}</span>
      </div>
      
      <div class="export-options">
        <div class="options-header">
          <span>{{ $t('settings.backup.selectData') }}</span>
          <div class="select-actions">
            <t-button variant="text" size="small" @click="selectAll">
              {{ $t('settings.backup.selectAll') }}
            </t-button>
            <t-button variant="text" size="small" @click="deselectAll">
              {{ $t('settings.backup.deselectAll') }}
            </t-button>
          </div>
        </div>
        
        <div class="options-grid">
          <div 
            v-for="option in exportOptions" 
            :key="option.key"
            class="option-item"
          >
            <t-checkbox 
              v-model="selectedOptions[option.key]"
              :disabled="option.count === 0"
            >
              <span class="option-label">{{ option.label }}</span>
              <span class="option-count">({{ option.count }})</span>
            </t-checkbox>
          </div>
        </div>
      </div>

      <div class="action-row">
        <t-button 
          theme="primary" 
          :loading="exporting"
          :disabled="!hasSelectedOptions"
          @click="handleExport"
        >
          <template #icon><t-icon name="download" /></template>
          {{ $t('settings.backup.exportBtn') }}
        </t-button>
      </div>
    </div>

    <!-- Import Section -->
    <div class="backup-section">
      <div class="section-label">
        <t-icon name="upload" class="section-icon" />
        <span>{{ $t('settings.backup.import') }}</span>
      </div>

      <div class="import-area">
        <t-upload
          ref="uploadRef"
          v-model="fileList"
          theme="custom"
          accept=".zip"
          :auto-upload="false"
          :max="1"
          @change="handleFileChange"
        >
          <template #default>
            <div class="upload-trigger" :class="{ 'has-file': fileList.length > 0 }">
              <t-icon v-if="fileList.length === 0" name="upload" class="upload-icon" />
              <t-icon v-else name="file-zip" class="upload-icon file-icon" />
              <div class="upload-text">
                <template v-if="fileList.length === 0">
                  <p class="main-text">{{ $t('settings.backup.dropFile') }}</p>
                  <p class="sub-text">{{ $t('settings.backup.supportFormat') }}</p>
                </template>
                <template v-else>
                  <p class="main-text">{{ fileList[0].name }}</p>
                  <p class="sub-text">{{ formatFileSize(fileList[0].size || 0) }}</p>
                </template>
              </div>
            </div>
          </template>
        </t-upload>

        <div v-if="fileList.length > 0" class="import-options">
          <t-checkbox v-model="skipExisting">
            {{ $t('settings.backup.skipExisting') }}
          </t-checkbox>
        </div>
      </div>

      <div class="action-row">
        <t-button 
          v-if="fileList.length > 0"
          variant="outline"
          @click="clearFile"
        >
          {{ $t('settings.backup.clearFile') }}
        </t-button>
        <t-button 
          theme="primary" 
          :loading="importing"
          :disabled="fileList.length === 0"
          @click="handleImport"
        >
          <template #icon><t-icon name="upload" /></template>
          {{ $t('settings.backup.importBtn') }}
        </t-button>
      </div>
    </div>

    <!-- Import Result Dialog -->
    <t-dialog
      v-model:visible="showResultDialog"
      :header="$t('settings.backup.importResult')"
      :footer="false"
      width="500px"
    >
      <div class="import-result">
        <div class="result-summary">
          <t-icon name="check-circle-filled" class="success-icon" />
          <span>{{ $t('settings.backup.importSuccess') }}</span>
        </div>
        
        <div class="result-details">
          <div v-if="importResult.knowledge_bases_imported > 0" class="result-item">
            <span>{{ $t('settings.backup.knowledgeBases') }}</span>
            <span class="count">+{{ importResult.knowledge_bases_imported }}</span>
          </div>
          <div v-if="importResult.knowledge_imported > 0" class="result-item">
            <span>{{ $t('settings.backup.knowledge') }}</span>
            <span class="count">+{{ importResult.knowledge_imported }}</span>
          </div>
          <div v-if="importResult.chunks_imported > 0" class="result-item">
            <span>{{ $t('settings.backup.chunks') }}</span>
            <span class="count">+{{ importResult.chunks_imported }}</span>
          </div>
          <div v-if="importResult.sessions_imported > 0" class="result-item">
            <span>{{ $t('settings.backup.sessions') }}</span>
            <span class="count">+{{ importResult.sessions_imported }}</span>
          </div>
          <div v-if="importResult.messages_imported > 0" class="result-item">
            <span>{{ $t('settings.backup.messages') }}</span>
            <span class="count">+{{ importResult.messages_imported }}</span>
          </div>
          <div v-if="importResult.models_imported > 0" class="result-item">
            <span>{{ $t('settings.backup.models') }}</span>
            <span class="count">+{{ importResult.models_imported }}</span>
          </div>
          <div v-if="importResult.agents_imported > 0" class="result-item">
            <span>{{ $t('settings.backup.agents') }}</span>
            <span class="count">+{{ importResult.agents_imported }}</span>
          </div>
          <div v-if="importResult.tags_imported > 0" class="result-item">
            <span>{{ $t('settings.backup.tags') }}</span>
            <span class="count">+{{ importResult.tags_imported }}</span>
          </div>
        </div>

        <div v-if="importResult.errors && importResult.errors.length > 0" class="result-errors">
          <div class="errors-header">
            <t-icon name="error-circle" class="error-icon" />
            <span>{{ $t('settings.backup.importErrors') }}</span>
          </div>
          <div class="errors-list">
            <div v-for="(error, index) in importResult.errors" :key="index" class="error-item">
              {{ error }}
            </div>
          </div>
        </div>

        <div class="result-actions">
          <t-button theme="primary" @click="showResultDialog = false">
            {{ $t('general.confirm') }}
          </t-button>
        </div>
      </div>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { useI18n } from 'vue-i18n'
import { getExportOptions, exportData, importData, type ExportOption, type ImportResult } from '@/api/backup'

const { t } = useI18n()

const exportOptions = ref<ExportOption[]>([])
const selectedOptions = reactive<Record<string, boolean>>({})
const exporting = ref(false)
const importing = ref(false)
const fileList = ref<any[]>([])
const skipExisting = ref(true)
const showResultDialog = ref(false)
const importResult = ref<ImportResult>({
  tenants_imported: 0,
  users_imported: 0,
  knowledge_bases_imported: 0,
  knowledge_imported: 0,
  chunks_imported: 0,
  sessions_imported: 0,
  messages_imported: 0,
  models_imported: 0,
  credentials_imported: 0,
  tags_imported: 0,
  agents_imported: 0,
  mcp_services_imported: 0,
})

const hasSelectedOptions = computed(() => {
  return Object.values(selectedOptions).some(v => v)
})

const loadExportOptions = async () => {
  try {
    const res = await getExportOptions()
    exportOptions.value = res.data
    // Initialize selected options
    res.data.forEach(opt => {
      selectedOptions[opt.key] = opt.count > 0
    })
  } catch (error) {
    console.error('Failed to load export options:', error)
  }
}

const selectAll = () => {
  exportOptions.value.forEach(opt => {
    if (opt.count > 0) {
      selectedOptions[opt.key] = true
    }
  })
}

const deselectAll = () => {
  Object.keys(selectedOptions).forEach(key => {
    selectedOptions[key] = false
  })
}

const handleExport = async () => {
  exporting.value = true
  try {
    const options: Record<string, boolean> = {}
    Object.entries(selectedOptions).forEach(([key, value]) => {
      if (value) {
        options[key] = true
      }
    })
    
    const blob = await exportData(options)
    
    // Download file
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `weknora_backup_${new Date().toISOString().slice(0, 10)}.zip`
    document.body.appendChild(a)
    a.click()
    window.URL.revokeObjectURL(url)
    document.body.removeChild(a)
    
    MessagePlugin.success(t('settings.backup.exportSuccess'))
  } catch (error: any) {
    MessagePlugin.error(error.message || t('settings.backup.exportFailed'))
  } finally {
    exporting.value = false
  }
}

const handleFileChange = (files: any[]) => {
  fileList.value = files
}

const clearFile = () => {
  fileList.value = []
}

const handleImport = async () => {
  if (fileList.value.length === 0) return
  
  importing.value = true
  try {
    const file = fileList.value[0].raw || fileList.value[0]
    const res = await importData(file, skipExisting.value)
    importResult.value = res.data
    showResultDialog.value = true
    fileList.value = []
  } catch (error: any) {
    MessagePlugin.error(error.message || t('settings.backup.importFailed'))
  } finally {
    importing.value = false
  }
}

const formatFileSize = (bytes: number): string => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

onMounted(() => {
  loadExportOptions()
})
</script>

<style lang="less" scoped>
.backup-settings {
  .section-header {
    margin-bottom: 24px;
  }

  .section-title {
    font-size: 16px;
    font-weight: 600;
    color: #333;
    margin: 0 0 8px 0;
  }

  .section-desc {
    font-size: 13px;
    color: #666;
    margin: 0;
  }
}

.backup-section {
  background: #f8f9fa;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;

  .section-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    font-weight: 500;
    color: #333;
    margin-bottom: 16px;

    .section-icon {
      font-size: 18px;
      color: #07C05F;
    }
  }
}

.export-options {
  background: #fff;
  border-radius: 6px;
  padding: 16px;
  margin-bottom: 16px;

  .options-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
    font-size: 13px;
    color: #666;

    .select-actions {
      display: flex;
      gap: 8px;
    }
  }

  .options-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }

  .option-item {
    .option-label {
      color: #333;
    }

    .option-count {
      color: #999;
      font-size: 12px;
      margin-left: 4px;
    }
  }
}

.import-area {
  background: #fff;
  border-radius: 6px;
  padding: 16px;
  margin-bottom: 16px;

  .upload-trigger {
    border: 2px dashed #d9d9d9;
    border-radius: 6px;
    padding: 32px;
    text-align: center;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      border-color: #07C05F;
    }

    &.has-file {
      border-color: #07C05F;
      background: rgba(7, 192, 95, 0.05);
    }

    .upload-icon {
      font-size: 40px;
      color: #999;
      margin-bottom: 12px;

      &.file-icon {
        color: #07C05F;
      }
    }

    .upload-text {
      .main-text {
        font-size: 14px;
        color: #333;
        margin: 0 0 4px 0;
      }

      .sub-text {
        font-size: 12px;
        color: #999;
        margin: 0;
      }
    }
  }

  .import-options {
    margin-top: 12px;
    padding-top: 12px;
    border-top: 1px solid #eee;
  }
}

.action-row {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.import-result {
  .result-summary {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 16px;
    font-weight: 500;
    color: #333;
    margin-bottom: 20px;

    .success-icon {
      font-size: 24px;
      color: #07C05F;
    }
  }

  .result-details {
    background: #f8f9fa;
    border-radius: 6px;
    padding: 16px;
    margin-bottom: 16px;

    .result-item {
      display: flex;
      justify-content: space-between;
      padding: 8px 0;
      border-bottom: 1px solid #eee;

      &:last-child {
        border-bottom: none;
      }

      .count {
        color: #07C05F;
        font-weight: 500;
      }
    }
  }

  .result-errors {
    background: #fff5f5;
    border-radius: 6px;
    padding: 16px;
    margin-bottom: 16px;

    .errors-header {
      display: flex;
      align-items: center;
      gap: 8px;
      font-weight: 500;
      color: #e34d59;
      margin-bottom: 12px;

      .error-icon {
        font-size: 18px;
      }
    }

    .errors-list {
      max-height: 150px;
      overflow-y: auto;

      .error-item {
        font-size: 12px;
        color: #666;
        padding: 4px 0;
      }
    }
  }

  .result-actions {
    display: flex;
    justify-content: flex-end;
  }
}
</style>
