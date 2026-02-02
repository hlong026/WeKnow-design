<template>
  <div class="user-menu" ref="menuRef">
    <!-- 用户按钮 -->
    <div class="user-button" @click="toggleMenu">
      <div class="user-avatar">
        <img v-if="userAvatar" :src="userAvatar" :alt="$t('common.avatar')" />
        <span v-else class="avatar-placeholder">{{ userInitial }}</span>
      </div>
      <div class="user-info">
        <div class="user-name">{{ userName }}</div>
        <div class="user-email">{{ userEmail }}</div>
      </div>
      <t-icon :name="menuVisible ? 'chevron-up' : 'chevron-down'" class="dropdown-icon" />
    </div>

    <!-- 下拉菜单 -->
    <Transition name="dropdown">
      <div v-if="menuVisible" class="user-dropdown" @click.stop>
        <div class="menu-item" @click="handleQuickNav('models')">
          <t-icon name="control-platform" class="menu-icon" />
          <span>{{ $t('settings.modelManagement') }}</span>
        </div>
        <div class="menu-item" @click="handleModelSettings">
          <t-icon name="server" class="menu-icon" />
          <span>{{ $t('modelSettings.title') }}</span>
        </div>
        <div class="menu-item" @click="handleQuickNav('websearch')">
          <svg 
            width="16" 
            height="16" 
            viewBox="0 0 18 18" 
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            class="menu-icon svg-icon"
          >
            <circle cx="9" cy="9" r="7" stroke="currentColor" stroke-width="1.2" fill="none"/>
            <path d="M 9 2 A 3.5 7 0 0 0 9 16" stroke="currentColor" stroke-width="1.2" fill="none"/>
            <path d="M 9 2 A 3.5 7 0 0 1 9 16" stroke="currentColor" stroke-width="1.2" fill="none"/>
            <line x1="2.94" y1="5.5" x2="15.06" y2="5.5" stroke="currentColor" stroke-width="1.2" stroke-linecap="round"/>
            <line x1="2.94" y1="12.5" x2="15.06" y2="12.5" stroke="currentColor" stroke-width="1.2" stroke-linecap="round"/>
          </svg>
          <span>{{ $t('settings.webSearchConfig') }}</span>
        </div>
        <div class="menu-divider"></div>
        <div class="menu-item" @click="handleSettings">
          <t-icon name="setting" class="menu-icon" />
          <span>{{ $t('general.allSettings') }}</span>
        </div>
        <div class="menu-divider"></div>
        <div class="menu-item" @click="openWebsite">
          <t-icon name="home" class="menu-icon" />
          <span>{{ $t('common.website') }}</span>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUIStore } from '@/stores/ui'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const router = useRouter()
const uiStore = useUIStore()
const authStore = useAuthStore()

const menuRef = ref<HTMLElement>()
const menuVisible = ref(false)

// 单机版：固定的用户信息
const userInfo = ref({
  username: '本地用户',
  email: 'local@weknora.local',
  avatar: ''
})

const userName = computed(() => userInfo.value.username)
const userEmail = computed(() => userInfo.value.email)
const userAvatar = computed(() => userInfo.value.avatar)

// 用户名首字母（用于无头像时显示）
const userInitial = computed(() => {
  return userName.value.charAt(0).toUpperCase()
})

// 切换菜单显示
const toggleMenu = () => {
  menuVisible.value = !menuVisible.value
}

// 快捷导航到设置的特定部分
const handleQuickNav = (section: string) => {
  menuVisible.value = false
  uiStore.openSettings()
  router.push('/platform/settings')
  
  // 延迟一下，确保设置页面已经渲染
  setTimeout(() => {
    // 触发设置页面切换到对应section
    const event = new CustomEvent('settings-nav', { detail: { section } })
    window.dispatchEvent(event)
  }, 100)
}

// 跳转到模型设置页面
const handleModelSettings = () => {
  menuVisible.value = false
  router.push('/platform/model-settings')
}

// 打开设置
const handleSettings = () => {
  menuVisible.value = false
  uiStore.openSettings()
  router.push('/platform/settings')
}

// 打开官网
const openWebsite = () => {
  menuVisible.value = false
  // 链接暂时为空
  // window.open('', '_blank')
}

// 单机版：不需要从服务器加载用户信息
const loadUserInfo = async () => {
  // 单机版：已经在 ref 中设置了固定值，无需加载
}

// 点击外部关闭菜单
const handleClickOutside = (e: MouseEvent) => {
  if (menuRef.value && !menuRef.value.contains(e.target as Node)) {
    menuVisible.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  loadUserInfo()
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style lang="less" scoped>
.user-menu {
  position: relative;
  width: 100%;
}

.user-button {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s;
  background: transparent;

  &:hover {
    background: #f5f7fa;
  }

  &:active {
    transform: scale(0.98);
  }
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  flex-shrink: 0;
  background: linear-gradient(135deg, #07C05F 0%, #05A34E 100%);
  display: flex;
  align-items: center;
  justify-content: center;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .avatar-placeholder {
    color: #ffffff;
    font-size: 16px;
    font-weight: 600;
  }
}

.user-info {
  flex: 1;
  min-width: 0;
  text-align: left;

  .user-name {
    font-size: 14px;
    font-weight: 500;
    color: #333333;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .user-email {
    font-size: 12px;
    color: #666666;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
}

.dropdown-icon {
  font-size: 16px;
  color: #666666;
  flex-shrink: 0;
  transition: transform 0.2s;
}

.user-dropdown {
  position: absolute;
  bottom: 100%;
  left: 8px;
  right: 8px;
  margin-bottom: 8px;
  background: #ffffff;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
  border: 1px solid #e5e7eb;
  overflow: hidden;
  z-index: 1000;
}

.menu-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
  transition: all 0.2s;
  font-size: 14px;
  color: #333333;

  &:hover {
    background: #f5f7fa;
  }

  &.danger {
    color: #e34d59;

    &:hover {
      background: #fef0f0;
    }

    .menu-icon {
      color: #e34d59;
    }
  }

  .menu-icon {
    font-size: 16px;
    color: #666666;
    
    &.svg-icon {
      width: 16px;
      height: 16px;
      flex-shrink: 0;
    }
  }

  .menu-text-with-icon {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 6px;
    color: inherit;
    min-width: 0;

    span {
      display: inline-flex;
      align-items: center;
      min-width: 0;
    }
  }

  .menu-external-icon {
    width: 14px;
    height: 14px;
    color: #9ca3af;
    flex-shrink: 0;
    transition: color 0.2s ease;
    pointer-events: none;
  }

  &:hover .menu-external-icon {
    color: #07c05f;
  }
}

.menu-divider {
  height: 1px;
  background: #e5e7eb;
  margin: 4px 0;
}

// 下拉动画
.dropdown-enter-active,
.dropdown-leave-active {
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: translateY(8px);
}

.dropdown-enter-to,
.dropdown-leave-from {
  opacity: 1;
  transform: translateY(0);
}
</style>
