import { get, put } from '@/utils/request'

// 品牌配置接口
export interface BrandConfig {
  app_name?: string
  logo_url?: string
  favicon_url?: string
  primary_color?: string
  welcome_message?: string
  footer_text?: string
  copyright_text?: string
  show_logo?: boolean
}

// 租户信息接口
export interface TenantInfo {
  id: number
  name: string
  description?: string
  api_key?: string
  status?: string
  business?: string
  storage_quota?: number
  storage_used?: number
  brand_config?: BrandConfig
  created_at: string
  updated_at: string
}

// 搜索租户参数
export interface SearchTenantsParams {
  keyword?: string
  tenant_id?: number
  page?: number
  page_size?: number
}

// 搜索租户响应
export interface SearchTenantsResponse {
  success: boolean
  data?: {
    items: TenantInfo[]
    total: number
    page: number
    page_size: number
  }
  message?: string
}

/**
 * 获取所有租户列表（需要跨租户访问权限）
 * @deprecated 建议使用 searchTenants 代替，支持分页和搜索
 */
export async function listAllTenants(): Promise<{ success: boolean; data?: { items: TenantInfo[] }; message?: string }> {
  try {
    const response = await get('/api/v1/tenants/all')
    return response as unknown as { success: boolean; data?: { items: TenantInfo[] }; message?: string }
  } catch (error: any) {
    return {
      success: false,
      message: error.message || '获取租户列表失败'
    }
  }
}

/**
 * 搜索租户（支持分页、关键词搜索和租户ID过滤）
 */
export async function searchTenants(params: SearchTenantsParams = {}): Promise<SearchTenantsResponse> {
  try {
    const queryParams = new URLSearchParams()
    if (params.keyword) {
      queryParams.append('keyword', params.keyword)
    }
    if (params.tenant_id) {
      queryParams.append('tenant_id', String(params.tenant_id))
    }
    if (params.page) {
      queryParams.append('page', String(params.page))
    }
    if (params.page_size) {
      queryParams.append('page_size', String(params.page_size))
    }
    
    const queryString = queryParams.toString()
    const url = `/api/v1/tenants/search${queryString ? '?' + queryString : ''}`
    const response = await get(url)
    return response as unknown as SearchTenantsResponse
  } catch (error: any) {
    return {
      success: false,
      message: error.message || '搜索租户失败'
    }
  }
}



/**
 * 获取当前租户信息（通过 auth/me 接口）
 */
export async function getCurrentTenant(): Promise<{ success: boolean; data?: TenantInfo; message?: string }> {
  try {
    const response = await get('/api/v1/auth/me')
    const result = response as unknown as { success: boolean; data?: { tenant: TenantInfo }; message?: string }
    if (result.success && result.data?.tenant) {
      return {
        success: true,
        data: result.data.tenant
      }
    }
    return {
      success: false,
      message: result.message || '获取租户信息失败'
    }
  } catch (error: any) {
    return {
      success: false,
      message: error.message || '获取租户信息失败'
    }
  }
}

/**
 * 更新品牌配置（通过更新租户接口）
 */
export async function updateBrandConfig(config: BrandConfig): Promise<{ success: boolean; data?: TenantInfo; message?: string }> {
  try {
    // 先获取当前租户信息
    const currentTenant = await getCurrentTenant()
    if (!currentTenant.success || !currentTenant.data) {
      return {
        success: false,
        message: '无法获取当前租户信息'
      }
    }
    
    // 更新租户的 brand_config
    const response = await put(`/api/v1/tenants/${currentTenant.data.id}`, {
      brand_config: config
    })
    return response as unknown as { success: boolean; data?: TenantInfo; message?: string }
  } catch (error: any) {
    return {
      success: false,
      message: error.message || '更新品牌配置失败'
    }
  }
}
