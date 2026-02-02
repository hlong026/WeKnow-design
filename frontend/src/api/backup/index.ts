import { get, post, postUpload, getDown } from '@/utils/request'

export interface ExportOption {
  key: string
  label: string
  count: number
}

export interface ExportRequest {
  include_tenants?: boolean
  include_users?: boolean
  include_knowledge_bases?: boolean
  include_knowledge?: boolean
  include_chunks?: boolean
  include_sessions?: boolean
  include_messages?: boolean
  include_models?: boolean
  include_credentials?: boolean
  include_tags?: boolean
  include_agents?: boolean
  include_mcp_services?: boolean
}

export interface ImportResult {
  tenants_imported: number
  users_imported: number
  knowledge_bases_imported: number
  knowledge_imported: number
  chunks_imported: number
  sessions_imported: number
  messages_imported: number
  models_imported: number
  credentials_imported: number
  tags_imported: number
  agents_imported: number
  mcp_services_imported: number
  errors?: string[]
}

// Get export options with data counts
export function getExportOptions(): Promise<{ data: ExportOption[] }> {
  return get('/api/v1/system/backup/options')
}

// Export database data
export async function exportData(options: ExportRequest): Promise<Blob> {
  const response = await fetch('/api/v1/system/backup/export', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(options),
  })
  
  if (!response.ok) {
    throw new Error('Export failed')
  }
  
  return response.blob()
}

// Import database data
export function importData(
  file: File, 
  skipExisting: boolean = false,
  onUploadProgress?: (progressEvent: any) => void
): Promise<{ data: ImportResult }> {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('skip_existing', skipExisting.toString())
  return postUpload('/api/v1/system/backup/import', formData, onUploadProgress)
}
