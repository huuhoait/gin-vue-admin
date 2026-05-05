-- =============================================================================
-- Audit columns backfill — for existing deployments that can't restart the
-- Go server right now to let AutoMigrate add CreatedBy/UpdatedBy/DeletedBy.
--
-- Idempotent: every ALTER uses ADD COLUMN IF NOT EXISTS. Re-running the
-- script is a no-op once the columns exist.
--
-- Postgres syntax. For MySQL replace `BIGINT` with `BIGINT UNSIGNED` and
-- drop the IF NOT EXISTS clauses (MySQL 5.7 doesn't support them on ADD
-- COLUMN; use a stored procedure or check information_schema first).
--
-- Tables: every model that embeds global.GVA_MODEL. Sourced from
-- server/initialize/gorm.go (RegisterTables) plus the tenant + onlineusers
-- + oauth2server + sysmonitor plugins.
-- =============================================================================

DO $$
DECLARE
  tbl text;
  tables text[] := ARRAY[
    -- system core
    'sys_apis',
    'sys_ignore_apis',
    'sys_users',
    'sys_base_menus',
    'jwt_blacklists',
    'sys_authorities',
    'sys_dictionaries',
    'sys_operation_records',
    'sys_auto_code_histories',
    'sys_dictionary_details',
    'sys_base_menu_parameters',
    'sys_base_menu_btns',
    'sys_authority_btns',
    'sys_auto_code_packages',
    'sys_export_templates',
    'conditions',
    'join_templates',
    'sys_params',
    'sys_versions',
    'sys_errors',
    'sys_api_tokens',
    'sys_login_logs',
    'sys_policy_change_logs',
    'sys_data_change_logs',
    -- example
    'exa_files',
    'exa_customers',
    'exa_file_chunks',
    'exa_file_upload_and_downloads',
    'exa_attachment_categories',
    -- plugins
    'gva_announcements_info',
    'gva_oauth2_clients',
    'gva_oauth2_auth_codes',
    'gva_oauth2_tokens',
    'gva_tenants',
    'gva_tenant_packages'
  ];
BEGIN
  FOREACH tbl IN ARRAY tables LOOP
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = tbl) THEN
      EXECUTE format('ALTER TABLE %I ADD COLUMN IF NOT EXISTS created_by BIGINT DEFAULT 0', tbl);
      EXECUTE format('ALTER TABLE %I ADD COLUMN IF NOT EXISTS updated_by BIGINT DEFAULT 0', tbl);
      EXECUTE format('ALTER TABLE %I ADD COLUMN IF NOT EXISTS deleted_by BIGINT DEFAULT 0', tbl);
    END IF;
  END LOOP;
END $$;

-- After running: verify with
--   SELECT column_name FROM information_schema.columns
--   WHERE table_name = 'sys_users' AND column_name IN ('created_by','updated_by','deleted_by');
