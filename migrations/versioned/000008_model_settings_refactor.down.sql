-- Rollback: Model Settings Refactor

DO $$ BEGIN RAISE NOTICE '[Migration 000003] Rolling back model settings refactor'; END $$;

-- Drop tables in reverse order
DROP TABLE IF EXISTS model_usage_stats;
DROP TABLE IF EXISTS model_configurations;
DROP TABLE IF EXISTS provider_credentials;

DO $$ BEGIN RAISE NOTICE '[Migration 000003] Rollback completed'; END $$;
