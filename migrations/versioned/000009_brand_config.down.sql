-- Remove brand_config column from tenants table
ALTER TABLE tenants DROP COLUMN IF EXISTS brand_config;
