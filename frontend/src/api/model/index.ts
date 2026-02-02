import { get, post, put, del } from '../../utils/request';

// 模型类型定义
export interface ModelConfig {
  id?: string;
  tenant_id?: number;
  name: string;
  type: 'KnowledgeQA' | 'Embedding' | 'Rerank' | 'VLLM';
  source: 'local' | 'remote';
  description?: string;
  parameters: {
    base_url?: string;
    api_key?: string;
    provider?: string; // Provider identifier: openai, aliyun, zhipu, generic
    embedding_parameters?: {
      dimension?: number;
      truncate_prompt_tokens?: number;
    };
    interface_type?: 'ollama' | 'openai'; // VLLM专用
    parameter_size?: string; // Ollama模型参数大小 (e.g., "7B", "13B", "70B")
    extra_config?: Record<string, string>; // Provider-specific configuration
  };
  is_default?: boolean;
  is_builtin?: boolean;
  status?: string;
  created_at?: string;
  updated_at?: string;
  deleted_at?: string | null;
}

// 模型厂商信息
export interface ModelProviderOption {
  value: string;
  label: string;
  description: string;
  defaultUrls: Record<string, string>;
  modelTypes: string[];
}

// 厂商详情（新版）
export interface ProviderDetail {
  name: string;
  display_name: string;
  description: string;
  icon?: string;
  website?: string;
  docs_url?: string;
  auth_config: {
    type: string;
    fields: Array<{
      key: string;
      label: string;
      type: string;
      required: boolean;
      placeholder: string;
      help_text?: string;
    }>;
    help_text?: string;
    help_url?: string;
  };
  supported_types: string[];
  preset_models: PresetModel[];
  endpoints: Record<string, string>;
  features: {
    supports_streaming: boolean;
    supports_function_call: boolean;
    supports_vision: boolean;
    supports_json_mode: boolean;
    supports_custom_model: boolean;
  };
}

// 预置模型
export interface PresetModel {
  model_id: string;
  display_name: string;
  model_type: string;
  capabilities: string[];
  context_size: number;
  pricing?: {
    input_price: number;
    output_price: number;
    currency: string;
  };
  deprecated?: boolean;
}

// 创建模型
export function createModel(data: ModelConfig): Promise<ModelConfig> {
  return new Promise((resolve, reject) => {
    post('/api/v1/models', data)
      .then((response: any) => {
        if (response.success && response.data) {
          resolve(response.data);
        } else {
          reject(new Error(response.message || '创建模型失败'));
        }
      })
      .catch((error: any) => {
        console.error('创建模型失败:', error);
        reject(error);
      });
  });
}

// 获取模型列表
export function listModels(type?: string): Promise<ModelConfig[]> {
  return new Promise((resolve, reject) => {
    const url = `/api/v1/models`;
    get(url)
      .then((response: any) => {
        if (response.success && response.data) {
          if (type) {
            response.data = response.data.filter((item: ModelConfig) => item.type === type);
          }
          resolve(response.data);
        } else {
          resolve([]);
        }
      })
      .catch((error: any) => {
        console.error('获取模型列表失败:', error);
        resolve([]);
      });
  });
}

// 获取单个模型
export function getModel(id: string): Promise<ModelConfig> {
  return new Promise((resolve, reject) => {
    get(`/api/v1/models/${id}`)
      .then((response: any) => {
        if (response.success && response.data) {
          resolve(response.data);
        } else {
          reject(new Error(response.message || '获取模型失败'));
        }
      })
      .catch((error: any) => {
        console.error('获取模型失败:', error);
        reject(error);
      });
  });
}

// 更新模型
export function updateModel(id: string, data: Partial<ModelConfig>): Promise<ModelConfig> {
  return new Promise((resolve, reject) => {
    put(`/api/v1/models/${id}`, data)
      .then((response: any) => {
        if (response.success && response.data) {
          resolve(response.data);
        } else {
          reject(new Error(response.message || '更新模型失败'));
        }
      })
      .catch((error: any) => {
        console.error('更新模型失败:', error);
        reject(error);
      });
  });
}

// 删除模型
export function deleteModel(id: string): Promise<void> {
  return new Promise((resolve, reject) => {
    del(`/api/v1/models/${id}`)
      .then((response: any) => {
        if (response.success) {
          resolve();
        } else {
          reject(new Error(response.message || '删除模型失败'));
        }
      })
      .catch((error: any) => {
        console.error('删除模型失败:', error);
        reject(error);
      });
  });
}

// 获取模型厂商列表（兼容旧接口）
export function listModelProviders(modelType?: string): Promise<ModelProviderOption[]> {
  return new Promise((resolve, reject) => {
    const url = modelType
      ? `/api/v1/models/providers?model_type=${encodeURIComponent(modelType)}`
      : '/api/v1/models/providers';
    get(url)
      .then((response: any) => {
        resolve(response.data || []);
      })
      .catch((error: any) => {
        console.error('获取模型厂商列表失败:', error);
        resolve([]);
      });
  });
}

// 获取所有厂商列表（新接口）
export function listProviders(): Promise<ProviderDetail[]> {
  return new Promise((resolve, reject) => {
    get('/api/v1/providers')
      .then((response: any) => {
        resolve(response.data || []);
      })
      .catch((error: any) => {
        console.error('获取厂商列表失败:', error);
        resolve([]);
      });
  });
}

// 获取厂商详情
export function getProvider(provider: string): Promise<ProviderDetail> {
  return new Promise((resolve, reject) => {
    get(`/api/v1/providers/${provider}`)
      .then((response: any) => {
        if (response.success && response.data) {
          resolve(response.data);
        } else {
          reject(new Error(response.message || '获取厂商详情失败'));
        }
      })
      .catch(reject);
  });
}

// 获取厂商预置模型
export function getProviderModels(provider: string, modelType?: string): Promise<PresetModel[]> {
  return new Promise((resolve, reject) => {
    const url = modelType
      ? `/api/v1/providers/${provider}/models?model_type=${encodeURIComponent(modelType)}`
      : `/api/v1/providers/${provider}/models`;
    get(url)
      .then((response: any) => {
        resolve(response.data || []);
      })
      .catch((error: any) => {
        console.error('获取预置模型失败:', error);
        resolve([]);
      });
  });
}

