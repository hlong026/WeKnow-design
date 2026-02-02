import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

// 单机版：简化的 store，无需真实认证
export const useAuthStore = defineStore('auth', () => {
  // 固定的单机用户信息
  const user = ref({
    id: 'local-user',
    username: '本地用户',
    email: 'local@weknora.local',
    tenant_id: 1,
    can_access_all_tenants: true
  })

  const tenant = ref({
    id: 1,
    name: '本地租户',
    api_key: 'local-api-key',
    status: 'active'
  })

  const token = ref('local-mode-token')
  const knowledgeBases = ref<any[]>([])
  const currentKnowledgeBase = ref<any | null>(null)

  // 计算属性 - 单机版始终返回已登录状态
  const isLoggedIn = computed(() => true)
  const hasValidTenant = computed(() => true)
  const currentTenantId = computed(() => 1)
  const currentUserId = computed(() => 'local-user')
  const canAccessAllTenants = computed(() => true)
  const effectiveTenantId = computed(() => 1)

  // 简化的方法 - 单机版不需要实际操作
  const setUser = (userData: any) => {
    // 单机版：不需要设置用户
  }

  const setTenant = (tenantData: any) => {
    // 单机版：不需要设置租户
  }

  const setToken = (tokenValue: string) => {
    // 单机版：不需要设置 token
  }

  const setRefreshToken = (refreshTokenValue: string) => {
    // 单机版：不需要刷新 token
  }

  const setKnowledgeBases = (kbList: any[]) => {
    knowledgeBases.value = Array.isArray(kbList) ? kbList : []
  }

  const setCurrentKnowledgeBase = (kb: any | null) => {
    currentKnowledgeBase.value = kb
  }

  const setSelectedTenant = (tenantId: number | null, tenantName: string | null = null) => {
    // 单机版：不需要切换租户
  }

  const setAllTenants = (tenants: any[]) => {
    // 单机版：不需要管理多租户
  }

  const getSelectedTenant = () => {
    return 1
  }

  const logout = () => {
    // 单机版：不需要登出
  }

  const initFromStorage = () => {
    // 单机版：不需要从存储恢复
  }

  return {
    // 状态
    user,
    tenant,
    token,
    knowledgeBases,
    currentKnowledgeBase,
    
    // 计算属性
    isLoggedIn,
    hasValidTenant,
    currentTenantId,
    currentUserId,
    canAccessAllTenants,
    effectiveTenantId,
    
    // 方法
    setUser,
    setTenant,
    setToken,
    setRefreshToken,
    setKnowledgeBases,
    setCurrentKnowledgeBase,
    setSelectedTenant,
    setAllTenants,
    getSelectedTenant,
    logout,
    initFromStorage
  }
})
