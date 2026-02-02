import { get, post, put, del } from '../../utils/request'

// 凭证类型定义
export interface ProviderCredential {
  id?: string
  tenant_id?: number
  provider: string
  name: string
  credentials: Record<string, string>
  base_url?: string
  is_default?: boolean
  status?: 'active' | 'invalid' | 'quota_exceeded'
  quota_config?: {
    daily_limit?: number
    monthly_limit?: number
    token_limit?: number
    alert_threshold?: number
  }
  created_at?: string
  updated_at?: string
}

// 厂商元数据
export interface ProviderMetadata {
  name: string
  display_name: string
  description: string
  icon?: string
  website?: string
  docs_url?: string
  auth_config: {
    type: string
    fields: Array<{
      key: string
      label: string
      type: string
      required: boolean
      placeholder: string
      help_text?: string
    }>
    help_text?: string
    help_url?: string
  }
  supported_types: string[]
  preset_models: Array<{
    model_id: string
    display_name: string
    model_type: string
    capabilities: string[]
    context_size: number
    pricing?: {
      input_price: number
      output_price: number
      currency: string
    }
  }>
  endpoints: Record<string, string>
  features: {
    supports_streaming: boolean
    supports_function_call: boolean
    supports_vision: boolean
    supports_json_mode: boolean
    supports_custom_model: boolean
  }
}

// 获取所有支持的厂商列表
export function listProviders(): Promise<ProviderMetadata[]> {
  return new Promise((resolve, reject) => {
    get('/api/v1/providers')
      .then((response: any) => {
        resolve(response.data || [])
      })
      .catch((error: any) => {
        console.error('获取厂商列表失败:', error)
        resolve([])
      })
  })
}

// 获取厂商详情
export function getProvider(provider: string): Promise<ProviderMetadata> {
  return new Promise((resolve, reject) => {
    get(`/api/v1/providers/${provider}`)
      .then((response: any) => {
        if (response.success && response.data) {
          resolve(response.data)
        } else {
          reject(new Error(response.message || '获取厂商详情失败'))
        }
      })
      .catch(reject)
  })
}

// 创建凭证
export function createCredential(data: ProviderCredential): Promise<ProviderCredential> {
  return new Promise((resolve, reject) => {
    post('/api/v1/credentials', data)
      .then((response: any) => {
        if (response.success && response.data) {
          resolve(response.data)
        } else {
          reject(new Error(response.message || '创建凭证失败'))
        }
      })
      .catch(reject)
  })
}

// 获取凭证列表
export function listCredentials(provider?: string): Promise<ProviderCredential[]> {
  return new Promise((resolve, reject) => {
    const url = provider 
      ? `/api/v1/credentials?provider=${encodeURIComponent(provider)}`
      : '/api/v1/credentials'
    get(url)
      .then((response: any) => {
        resolve(response.data || [])
      })
      .catch((error: any) => {
        console.error('获取凭证列表失败:', error)
        resolve([])
      })
  })
}

// 更新凭证
export function updateCredential(id: string, data: Partial<ProviderCredential>): Promise<ProviderCredential> {
  return new Promise((resolve, reject) => {
    put(`/api/v1/credentials/${id}`, data)
      .then((response: any) => {
        if (response.success && response.data) {
          resolve(response.data)
        } else {
          reject(new Error(response.message || '更新凭证失败'))
        }
      })
      .catch(reject)
  })
}

// 删除凭证
export function deleteCredential(id: string): Promise<void> {
  return new Promise((resolve, reject) => {
    del(`/api/v1/credentials/${id}`)
      .then((response: any) => {
        if (response.success) {
          resolve()
        } else {
          reject(new Error(response.message || '删除凭证失败'))
        }
      })
      .catch(reject)
  })
}

// 测试凭证连接
export function testCredential(id: string): Promise<{ success: boolean; message: string }> {
  return new Promise((resolve, reject) => {
    post(`/api/v1/credentials/${id}/test`)
      .then((response: any) => {
        resolve(response.data || { success: response.success, message: response.message })
      })
      .catch(reject)
  })
}
