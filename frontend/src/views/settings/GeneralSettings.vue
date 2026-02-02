<template>
  <div class="general-settings">
    <div class="section-header">
      <h2>{{ $t('general.title') }}</h2>
      <p class="section-description">{{ $t('general.description') }}</p>
    </div>

    <div class="settings-group">
      <!-- 语言选择 -->
      <div class="setting-row">
        <div class="setting-info">
          <label>{{ $t('language.language') }}</label>
          <p class="desc">{{ $t('language.languageDescription') }}</p>
        </div>
        <div class="setting-control">
          <t-select
            v-model="localLanguage"
            :placeholder="$t('language.selectLanguage')"
            @change="handleLanguageChange"
            style="width: 280px;"
          >
            <t-option value="zh-CN" :label="$t('language.zhCN')">{{ $t('language.zhCN') }}</t-option>
            <t-option value="en-US" :label="$t('language.enUS')">{{ $t('language.enUS') }}</t-option>
            <t-option value="ru-RU" :label="$t('language.ruRU')">{{ $t('language.ruRU') }}</t-option>
            <t-option value="ko-KR" :label="$t('language.koKR')">{{ $t('language.koKR') }}</t-option>
          </t-select>
        </div>
      </div>
    </div>

    <!-- 品牌配置 -->
    <div class="section-header" style="margin-top: 40px;">
      <h2>品牌配置</h2>
      <p class="section-description">自定义应用的 Logo、名称和主题</p>
    </div>

    <div class="settings-group">
      <!-- 应用名称 -->
      <div class="setting-row">
        <div class="setting-info">
          <label>应用名称</label>
          <p class="desc">自定义显示的应用名称</p>
        </div>
        <div class="setting-control">
          <t-input
            v-model="brandConfig.app_name"
            placeholder="WeKnora"
            style="width: 280px;"
          />
        </div>
      </div>

      <!-- 显示 Logo -->
      <div class="setting-row">
        <div class="setting-info">
          <label>显示 Logo</label>
          <p class="desc">控制登录页面左上角是否显示 Logo</p>
        </div>
        <div class="setting-control">
          <t-switch
            v-model="brandConfig.show_logo"
            size="large"
          />
        </div>
      </div>

      <!-- Logo 上传 -->
      <div class="setting-row">
        <div class="setting-info">
          <label>Logo 图片</label>
          <p class="desc">上传自定义 Logo 图片（支持 PNG、JPG、SVG，建议尺寸 200x60）</p>
        </div>
        <div class="setting-control">
          <div class="image-upload-group">
            <div class="upload-area" @click="triggerLogoUpload" @dragover.prevent @drop.prevent="handleLogoDrop">
              <input
                ref="logoInputRef"
                type="file"
                accept="image/png,image/jpeg,image/svg+xml,image/gif"
                style="display: none"
                @change="handleLogoUpload"
              />
              <div v-if="brandConfig.logo_url" class="image-preview">
                <img :src="brandConfig.logo_url" alt="Logo Preview" @error="handleLogoError" />
                <div class="preview-overlay">
                  <span class="overlay-text">点击更换</span>
                </div>
                <t-button
                  theme="danger"
                  size="small"
                  shape="circle"
                  class="remove-btn"
                  @click.stop="removeLogo"
                >
                  <template #icon><t-icon name="close" /></template>
                </t-button>
              </div>
              <div v-else class="upload-placeholder">
                <t-icon name="upload" size="24px" />
                <span>点击或拖拽上传 Logo</span>
              </div>
            </div>
            <t-input
              v-model="brandConfig.logo_url"
              placeholder="或输入图片 URL"
              style="width: 280px; margin-top: 8px;"
            />
          </div>
        </div>
      </div>

      <!-- Favicon 上传 -->
      <div class="setting-row">
        <div class="setting-info">
          <label>Favicon 图标</label>
          <p class="desc">上传浏览器标签页图标（支持 ICO、PNG，建议尺寸 32x32）</p>
        </div>
        <div class="setting-control">
          <div class="image-upload-group">
            <div class="upload-area favicon-area" @click="triggerFaviconUpload" @dragover.prevent @drop.prevent="handleFaviconDrop">
              <input
                ref="faviconInputRef"
                type="file"
                accept="image/x-icon,image/png,image/ico"
                style="display: none"
                @change="handleFaviconUpload"
              />
              <div v-if="brandConfig.favicon_url" class="image-preview favicon-preview">
                <img :src="brandConfig.favicon_url" alt="Favicon Preview" />
                <div class="preview-overlay">
                  <span class="overlay-text">更换</span>
                </div>
                <t-button
                  theme="danger"
                  size="small"
                  shape="circle"
                  class="remove-btn"
                  @click.stop="removeFavicon"
                >
                  <template #icon><t-icon name="close" /></template>
                </t-button>
              </div>
              <div v-else class="upload-placeholder">
                <t-icon name="upload" size="20px" />
                <span>上传图标</span>
              </div>
            </div>
            <t-input
              v-model="brandConfig.favicon_url"
              placeholder="或输入图标 URL"
              style="width: 280px; margin-top: 8px;"
            />
          </div>
        </div>
      </div>

      <!-- 主题色 -->
      <div class="setting-row">
        <div class="setting-info">
          <label>主题色</label>
          <p class="desc">自定义应用的主题颜色</p>
        </div>
        <div class="setting-control">
          <div class="color-input-group">
            <t-input
              v-model="brandConfig.primary_color"
              placeholder="#07C05F"
              style="width: 200px;"
            />
            <input
              type="color"
              v-model="brandConfig.primary_color"
              class="color-picker"
            />
          </div>
        </div>
      </div>

      <!-- 欢迎语 -->
      <div class="setting-row">
        <div class="setting-info">
          <label>欢迎语</label>
          <p class="desc">用户进入对话页面时显示的欢迎消息</p>
        </div>
        <div class="setting-control">
          <t-textarea
            v-model="brandConfig.welcome_message"
            placeholder="你好！我是智能助手，有什么可以帮助你的？"
            :autosize="{ minRows: 2, maxRows: 4 }"
            style="width: 280px;"
          />
        </div>
      </div>

      <!-- 页脚文字 -->
      <div class="setting-row">
        <div class="setting-info">
          <label>页脚文字</label>
          <p class="desc">显示在页面底部的自定义文字</p>
        </div>
        <div class="setting-control">
          <t-input
            v-model="brandConfig.footer_text"
            placeholder="Powered by WeKnora"
            style="width: 280px;"
          />
        </div>
      </div>

      <!-- 版权信息 -->
      <div class="setting-row">
        <div class="setting-info">
          <label>版权信息</label>
          <p class="desc">显示在页面底部的版权声明</p>
        </div>
        <div class="setting-control">
          <t-input
            v-model="brandConfig.copyright_text"
            placeholder="© 2024 Your Company. All rights reserved."
            style="width: 280px;"
          />
        </div>
      </div>

      <!-- 保存按钮 -->
      <div class="setting-row save-row">
        <div class="setting-info"></div>
        <div class="setting-control">
          <t-button theme="primary" :loading="saving" @click="handleSaveBrandConfig">
            保存品牌配置
          </t-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { MessagePlugin } from 'tdesign-vue-next'
import { useI18n } from 'vue-i18n'
import { getCurrentTenant, updateBrandConfig, type BrandConfig } from '@/api/tenant'

const { t, locale } = useI18n()

// 本地状态
const localLanguage = ref('zh-CN')
const saving = ref(false)

// 文件上传引用
const logoInputRef = ref<HTMLInputElement | null>(null)
const faviconInputRef = ref<HTMLInputElement | null>(null)

// 品牌配置
const brandConfig = reactive<BrandConfig>({
  app_name: '',
  logo_url: '',
  favicon_url: '',
  primary_color: '#07C05F',
  welcome_message: '',
  footer_text: '',
  copyright_text: '',
  show_logo: true
})

// 初始化加载
onMounted(async () => {
  const savedLocale = localStorage.getItem('locale')
  if (savedLocale) {
    localLanguage.value = savedLocale
    locale.value = savedLocale
  } else {
    localLanguage.value = locale.value
  }
  await loadBrandConfig()
})

// 加载品牌配置
const loadBrandConfig = async () => {
  try {
    const res = await getCurrentTenant()
    if (res.success && res.data?.brand_config) {
      Object.assign(brandConfig, res.data.brand_config)
    }
  } catch (error) {
    console.error('Failed to load brand config:', error)
  }
}

// 处理语言变化
const handleLanguageChange = () => {
  locale.value = localLanguage.value
  localStorage.setItem('locale', localLanguage.value)
  MessagePlugin.success(t('language.languageSaved'))
}

// 触发 Logo 上传
const triggerLogoUpload = () => {
  logoInputRef.value?.click()
}

// 触发 Favicon 上传
const triggerFaviconUpload = () => {
  faviconInputRef.value?.click()
}

// 将文件转换为 Base64
const fileToBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onload = () => resolve(reader.result as string)
    reader.onerror = reject
    reader.readAsDataURL(file)
  })
}

// 处理 Logo 上传
const handleLogoUpload = async (e: Event) => {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  // 检查文件大小（限制 500KB）
  if (file.size > 500 * 1024) {
    MessagePlugin.warning('Logo 图片大小不能超过 500KB')
    return
  }

  try {
    const base64 = await fileToBase64(file)
    brandConfig.logo_url = base64
    MessagePlugin.success('Logo 上传成功')
  } catch (error) {
    MessagePlugin.error('Logo 上传失败')
  }
  
  // 清空 input 以便重复上传同一文件
  input.value = ''
}

// 处理 Logo 拖拽上传
const handleLogoDrop = async (e: DragEvent) => {
  const file = e.dataTransfer?.files?.[0]
  if (!file || !file.type.startsWith('image/')) {
    MessagePlugin.warning('请上传图片文件')
    return
  }

  if (file.size > 500 * 1024) {
    MessagePlugin.warning('Logo 图片大小不能超过 500KB')
    return
  }

  try {
    const base64 = await fileToBase64(file)
    brandConfig.logo_url = base64
    MessagePlugin.success('Logo 上传成功')
  } catch (error) {
    MessagePlugin.error('Logo 上传失败')
  }
}

// 处理 Favicon 上传
const handleFaviconUpload = async (e: Event) => {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return

  // 检查文件大小（限制 100KB）
  if (file.size > 100 * 1024) {
    MessagePlugin.warning('Favicon 图标大小不能超过 100KB')
    return
  }

  try {
    const base64 = await fileToBase64(file)
    brandConfig.favicon_url = base64
    MessagePlugin.success('Favicon 上传成功')
  } catch (error) {
    MessagePlugin.error('Favicon 上传失败')
  }
  
  input.value = ''
}

// 处理 Favicon 拖拽上传
const handleFaviconDrop = async (e: DragEvent) => {
  const file = e.dataTransfer?.files?.[0]
  if (!file || !file.type.startsWith('image/')) {
    MessagePlugin.warning('请上传图片文件')
    return
  }

  if (file.size > 100 * 1024) {
    MessagePlugin.warning('Favicon 图标大小不能超过 100KB')
    return
  }

  try {
    const base64 = await fileToBase64(file)
    brandConfig.favicon_url = base64
    MessagePlugin.success('Favicon 上传成功')
  } catch (error) {
    MessagePlugin.error('Favicon 上传失败')
  }
}

// 移除 Logo
const removeLogo = () => {
  brandConfig.logo_url = ''
}

// 移除 Favicon
const removeFavicon = () => {
  brandConfig.favicon_url = ''
}

// 处理 Logo 加载错误
const handleLogoError = (e: Event) => {
  const img = e.target as HTMLImageElement
  img.style.display = 'none'
}

// 保存品牌配置
const handleSaveBrandConfig = async () => {
  saving.value = true
  try {
    const res = await updateBrandConfig(brandConfig)
    if (res.success) {
      MessagePlugin.success('品牌配置保存成功')
      applyBrandConfig()
    } else {
      MessagePlugin.error(res.message || '保存失败')
    }
  } catch (error: any) {
    MessagePlugin.error(error.message || '保存失败')
  } finally {
    saving.value = false
  }
}

// 应用品牌配置
const applyBrandConfig = () => {
  if (brandConfig.app_name) {
    document.title = brandConfig.app_name
  }
  
  if (brandConfig.favicon_url) {
    const link = document.querySelector("link[rel*='icon']") as HTMLLinkElement || document.createElement('link')
    link.type = 'image/x-icon'
    link.rel = 'shortcut icon'
    link.href = brandConfig.favicon_url
    document.head.appendChild(link)
  }
  
  if (brandConfig.primary_color) {
    document.documentElement.style.setProperty('--td-brand-color', brandConfig.primary_color)
  }
  
  localStorage.setItem('WeKnora_brand_config', JSON.stringify(brandConfig))
}
</script>

<style lang="less" scoped>
.general-settings {
  width: 100%;
}

.section-header {
  margin-bottom: 32px;

  h2 {
    font-size: 20px;
    font-weight: 600;
    color: #333333;
    margin: 0 0 8px 0;
  }

  .section-description {
    font-size: 14px;
    color: #666666;
    margin: 0;
    line-height: 1.5;
  }
}

.settings-group {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.setting-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 20px 0;
  border-bottom: 1px solid #e5e7eb;

  &:last-child {
    border-bottom: none;
  }

  &.save-row {
    border-bottom: none;
    padding-top: 24px;
  }
}

.setting-info {
  flex: 1;
  max-width: 65%;
  padding-right: 24px;

  label {
    font-size: 15px;
    font-weight: 500;
    color: #333333;
    display: block;
    margin-bottom: 4px;
  }

  .desc {
    font-size: 13px;
    color: #666666;
    margin: 0;
    line-height: 1.5;
  }
}

.setting-control {
  flex-shrink: 0;
  min-width: 280px;
  display: flex;
  justify-content: flex-end;
  align-items: flex-start;
}

.image-upload-group {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.upload-area {
  width: 280px;
  height: 80px;
  border: 2px dashed #d9d9d9;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #fafafa;
  position: relative;
  overflow: hidden;

  &:hover {
    border-color: #07C05F;
    background: #f0fdf4;
  }

  &.favicon-area {
    width: 100px;
    height: 100px;
  }
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  color: #999;
  font-size: 13px;

  .t-icon {
    color: #bbb;
  }
}

.image-preview {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;

  img {
    max-width: 90%;
    max-height: 90%;
    object-fit: contain;
  }

  &.favicon-preview img {
    max-width: 60%;
    max-height: 60%;
  }
}

.preview-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;

  .overlay-text {
    color: #fff;
    font-size: 13px;
  }
}

.upload-area:hover .preview-overlay {
  opacity: 1;
}

.remove-btn {
  position: absolute;
  top: 4px;
  right: 4px;
  opacity: 0;
  transition: opacity 0.3s ease;
  z-index: 10;
}

.upload-area:hover .remove-btn {
  opacity: 1;
}

.color-input-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.color-picker {
  width: 40px;
  height: 32px;
  padding: 0;
  border: 1px solid #d9d9d9;
  border-radius: 4px;
  cursor: pointer;
  
  &::-webkit-color-swatch-wrapper {
    padding: 2px;
  }
  
  &::-webkit-color-swatch {
    border: none;
    border-radius: 2px;
  }
}
</style>
