import axios from 'axios'
import { put } from '@/utils/request'

export interface ExtractContentRequest {
  platform: string
  videoUrl: string
  kbId: string
}

export interface ExtractContentResponse {
  success: boolean
  message: string
  data?: {
    content: string
    knowledgeId: string
  }
}

// API基础URL
const BASE_URL = import.meta.env.VITE_IS_DOCKER ? "" : "http://localhost:8080";

/**
 * 从自媒体平台提取文案内容
 * 设置 8 分钟超时，因为 Coze Workflow 执行时间较长
 */
export async function extractSocialMediaContent(data: ExtractContentRequest): Promise<ExtractContentResponse> {
  const response = await axios.post(`${BASE_URL}/api/v1/social-media/extract`, data, {
    timeout: 8 * 60 * 1000, // 8 分钟
    headers: {
      'Content-Type': 'application/json'
    }
  })
  return response.data as ExtractContentResponse
}


/**
 * 更新知识库的阿里云 API Key
 */
export async function updateAliyunAPIKey(kbId: string, aliyunApiKey: string): Promise<{ success: boolean; message: string }> {
  const response: any = await put(`/api/v1/initialization/kb/${kbId}/aliyun-api-key`, {
    aliyunApiKey
  })
  return response as { success: boolean; message: string }
}
