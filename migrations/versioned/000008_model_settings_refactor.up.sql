-- Migration: Model Settings Refactor
-- Description: 重构模型设置，分离凭证和模型配置

DO $$ BEGIN RAISE NOTICE '[Migration 000003] Starting model settings refactor'; END $$;

-- 1. 创建厂商凭证表
DO $$ BEGIN RAISE NOTICE '[Migration 000003] Creating table: provider_credentials'; END $$;
CREATE TABLE IF NOT EXISTS provider_credentials (
    id VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4()::varchar,
    tenant_id INTEGER NOT NULL,
    provider VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    credentials JSONB NOT NULL DEFAULT '{}',
    base_url VARCHAR(500),
    is_default BOOLEAN DEFAULT false,
    status VARCHAR(50) DEFAULT 'active',
    quota_config JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT uq_credentials_tenant_provider_name UNIQUE(tenant_id, provider, name)
);

CREATE INDEX IF NOT EXISTS idx_credentials_tenant_provider ON provider_credentials(tenant_id, provider);
CREATE INDEX IF NOT EXISTS idx_credentials_deleted_at ON provider_credentials(deleted_at);

-- 2. 创建新的模型配置表
DO $$ BEGIN RAISE NOTICE '[Migration 000003] Creating table: model_configurations'; END $$;
CREATE TABLE IF NOT EXISTS model_configurations (
    id VARCHAR(36) PRIMARY KEY DEFAULT uuid_generate_v4()::varchar,
    tenant_id INTEGER NOT NULL,
    credential_id VARCHAR(36) REFERENCES provider_credentials(id) ON DELETE SET NULL,
    provider VARCHAR(50) NOT NULL,
    model_id VARCHAR(255) NOT NULL,
    display_name VARCHAR(255) NOT NULL,
    model_type VARCHAR(50) NOT NULL,
    capabilities JSONB DEFAULT '[]',
    parameters JSONB DEFAULT '{}',
    pricing_info JSONB,
    is_enabled BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,
    is_builtin BOOLEAN DEFAULT false,
    tags JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT uq_models_tenant_provider_model UNIQUE(tenant_id, provider, model_id)
);

CREATE INDEX IF NOT EXISTS idx_model_configs_tenant_type ON model_configurations(tenant_id, model_type);
CREATE INDEX IF NOT EXISTS idx_model_configs_credential ON model_configurations(credential_id);
CREATE INDEX IF NOT EXISTS idx_model_configs_deleted_at ON model_configurations(deleted_at);

-- 3. 创建使用统计表
DO $$ BEGIN RAISE NOTICE '[Migration 000003] Creating table: model_usage_stats'; END $$;
CREATE TABLE IF NOT EXISTS model_usage_stats (
    id BIGSERIAL PRIMARY KEY,
    tenant_id INTEGER NOT NULL,
    model_id VARCHAR(36) NOT NULL,
    credential_id VARCHAR(36),
    date DATE NOT NULL,
    request_count INTEGER DEFAULT 0,
    input_tokens BIGINT DEFAULT 0,
    output_tokens BIGINT DEFAULT 0,
    total_cost DECIMAL(10, 4) DEFAULT 0,
    error_count INTEGER DEFAULT 0,
    avg_latency_ms INTEGER DEFAULT 0,
    
    CONSTRAINT uq_usage_tenant_model_date UNIQUE(tenant_id, model_id, date)
);

CREATE INDEX IF NOT EXISTS idx_usage_tenant_date ON model_usage_stats(tenant_id, date);
CREATE INDEX IF NOT EXISTS idx_usage_model_date ON model_usage_stats(model_id, date);

-- 4. 数据迁移：从旧 models 表迁移凭证
DO $$ BEGIN RAISE NOTICE '[Migration 000003] Migrating credentials from models table'; END $$;
INSERT INTO provider_credentials (id, tenant_id, provider, name, credentials, base_url, is_default, created_at, updated_at)
SELECT DISTINCT ON (tenant_id, COALESCE(parameters->>'provider', 'generic'))
    uuid_generate_v4()::varchar,
    tenant_id,
    COALESCE(parameters->>'provider', 'generic'),
    COALESCE(parameters->>'provider', 'generic') || ' 默认凭证',
    jsonb_build_object('api_key', parameters->>'api_key'),
    parameters->>'base_url',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
FROM models
WHERE deleted_at IS NULL
  AND parameters->>'api_key' IS NOT NULL
  AND parameters->>'api_key' != ''
ON CONFLICT (tenant_id, provider, name) DO NOTHING;

-- 5. 数据迁移：从旧 models 表迁移模型配置
DO $$ BEGIN RAISE NOTICE '[Migration 000003] Migrating model configurations from models table'; END $$;
INSERT INTO model_configurations (
    id, tenant_id, credential_id, provider, model_id, display_name,
    model_type, parameters, is_default, is_builtin, created_at, updated_at
)
SELECT 
    m.id,
    m.tenant_id,
    pc.id as credential_id,
    COALESCE(m.parameters->>'provider', 'generic') as provider,
    m.name as model_id,
    m.name as display_name,
    m.type as model_type,
    m.parameters,
    m.is_default,
    COALESCE(m.is_builtin, false),
    m.created_at,
    m.updated_at
FROM models m
LEFT JOIN provider_credentials pc ON 
    pc.tenant_id = m.tenant_id 
    AND pc.provider = COALESCE(m.parameters->>'provider', 'generic')
    AND pc.is_default = true
WHERE m.deleted_at IS NULL
ON CONFLICT (tenant_id, provider, model_id) DO NOTHING;

DO $$ BEGIN RAISE NOTICE '[Migration 000003] Model settings refactor completed'; END $$;
