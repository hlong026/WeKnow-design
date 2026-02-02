-- Add brand_config column to tenants table
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS brand_config JSONB DEFAULT NULL;

-- Add comment
COMMENT ON COLUMN tenants.brand_config IS 'Brand configuration for custom logo and app name';
