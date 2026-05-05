-- =============================================================================
-- EXPORT: Security menu + OAuth2 clients + Online sessions
--
-- Idempotent SQL for existing deployments (MySQL / MariaDB syntax).
-- Mirrors plugin bootstrap (oauth2server + onlineusers) and Casbin paths used
-- when system.router-prefix is /admin-api (Casbin "obj" = path with that
-- prefix stripped — see middleware/casbin_rbac.go).
--
-- BEFORE RUNNING:
--   1) Replace @admin_role_id if your super-admin authority_id is not 888.
--   2) If your router-prefix is NOT /admin-api, adjust casbin_rule.v1 and
--      verify sys_apis.path values against your sys_apis table (must match
--      what the UI / Casbin sync expects).
--
-- AFTER RUNNING:
--   Restart the admin server so Casbin reloads policies from casbin_rule.
-- =============================================================================

SET @admin_role_id := 888;

-- ---------------------------------------------------------------------------
-- 1) Top-level Security parent (router holder)
-- ---------------------------------------------------------------------------
INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 0, 0, 'security', 'security', false, 'view/routerHolder.vue', 6, 'admin.menu.security', 'lock'
WHERE NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE name = 'security');

-- ---------------------------------------------------------------------------
-- 2) OAuth2 Clients menu row (create if missing, then reparent + normalise)
-- ---------------------------------------------------------------------------
INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1,
       (SELECT id FROM sys_base_menus WHERE name = 'security' LIMIT 1),
       'oauth2Clients', 'oauth2Clients', false, 'plugin/oauth2server/view/clients.vue', 1,
       'admin.plugin.oauth2.menu_title', 'key'
WHERE NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE name = 'oauth2Clients');

UPDATE sys_base_menus
SET parent_id = (SELECT id FROM (SELECT id FROM sys_base_menus WHERE name = 'security') AS s),
    title     = 'admin.plugin.oauth2.menu_title',
    sort      = 1
WHERE name = 'oauth2Clients';

-- ---------------------------------------------------------------------------
-- 3) Online Sessions menu row
-- ---------------------------------------------------------------------------
INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 1,
       (SELECT id FROM sys_base_menus WHERE name = 'security' LIMIT 1),
       'onlineUsers', 'onlineUsers', false, 'plugin/onlineusers/view/online.vue', 2,
       'admin.plugin.online_users.menu_title', 'user'
WHERE NOT EXISTS (SELECT 1 FROM sys_base_menus WHERE name = 'onlineUsers');

UPDATE sys_base_menus
SET parent_id = (SELECT id FROM (SELECT id FROM sys_base_menus WHERE name = 'security') AS s),
    title     = 'admin.plugin.online_users.menu_title',
    sort      = 2
WHERE name = 'onlineUsers';

-- ---------------------------------------------------------------------------
-- 4) Optional: normalise other plugin menu titles (same as security_menu_seed.sql)
-- ---------------------------------------------------------------------------
UPDATE sys_base_menus SET title = 'admin.plugin.sysmonitor.menu_title' WHERE name = 'sysmonitor';
UPDATE sys_base_menus SET title = 'admin.plugin.tenant.menu_title'       WHERE name = 'tenants';

-- ---------------------------------------------------------------------------
-- 5) sys_apis — OAuth2 client admin + online session APIs (if not present)
-- ---------------------------------------------------------------------------
INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/oauth2Client/create', 'create OAuth2 client', 'OAuth2Client', 'POST'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/oauth2Client/create' AND method = 'POST');
INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/oauth2Client/update', 'update OAuth2 client', 'OAuth2Client', 'PUT'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/oauth2Client/update' AND method = 'PUT');
INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/oauth2Client/delete', 'delete OAuth2 client', 'OAuth2Client', 'DELETE'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/oauth2Client/delete' AND method = 'DELETE');
INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/oauth2Client/find', 'find OAuth2 client', 'OAuth2Client', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/oauth2Client/find' AND method = 'GET');
INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/oauth2Client/list', 'list OAuth2 clients', 'OAuth2Client', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/oauth2Client/list' AND method = 'GET');
INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/oauth2Client/regenerateSecret', 'rotate client secret', 'OAuth2Client', 'POST'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/oauth2Client/regenerateSecret' AND method = 'POST');
INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/oauth2/authorize', 'OAuth2 authorize endpoint', 'OAuth2', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/oauth2/authorize' AND method = 'GET');

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/onlineUsers/list', 'list online sessions', 'OnlineUsers', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/onlineUsers/list' AND method = 'GET');
INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
SELECT NOW(), NOW(), '/onlineUsers/kick', 'kick an online session', 'OnlineUsers', 'POST'
WHERE NOT EXISTS (SELECT 1 FROM sys_apis WHERE path = '/onlineUsers/kick' AND method = 'POST');

-- ---------------------------------------------------------------------------
-- 6) casbin_rule — grant private plugin APIs to super admin (v0 = authority id)
--    Paths must match CasbinHandler "obj" (URL path minus system.router-prefix).
-- ---------------------------------------------------------------------------
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', CAST(@admin_role_id AS CHAR), '/oauth2Client/create', 'POST'
WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = CAST(@admin_role_id AS CHAR) AND v1 = '/oauth2Client/create' AND v2 = 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', CAST(@admin_role_id AS CHAR), '/oauth2Client/update', 'PUT'
WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = CAST(@admin_role_id AS CHAR) AND v1 = '/oauth2Client/update' AND v2 = 'PUT');
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', CAST(@admin_role_id AS CHAR), '/oauth2Client/delete', 'DELETE'
WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = CAST(@admin_role_id AS CHAR) AND v1 = '/oauth2Client/delete' AND v2 = 'DELETE');
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', CAST(@admin_role_id AS CHAR), '/oauth2Client/find', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = CAST(@admin_role_id AS CHAR) AND v1 = '/oauth2Client/find' AND v2 = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', CAST(@admin_role_id AS CHAR), '/oauth2Client/list', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = CAST(@admin_role_id AS CHAR) AND v1 = '/oauth2Client/list' AND v2 = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', CAST(@admin_role_id AS CHAR), '/oauth2Client/regenerateSecret', 'POST'
WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = CAST(@admin_role_id AS CHAR) AND v1 = '/oauth2Client/regenerateSecret' AND v2 = 'POST');
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', CAST(@admin_role_id AS CHAR), '/oauth2/authorize', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = CAST(@admin_role_id AS CHAR) AND v1 = '/oauth2/authorize' AND v2 = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', CAST(@admin_role_id AS CHAR), '/onlineUsers/list', 'GET'
WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = CAST(@admin_role_id AS CHAR) AND v1 = '/onlineUsers/list' AND v2 = 'GET');
INSERT INTO casbin_rule (ptype, v0, v1, v2)
SELECT 'p', CAST(@admin_role_id AS CHAR), '/onlineUsers/kick', 'POST'
WHERE NOT EXISTS (SELECT 1 FROM casbin_rule WHERE ptype = 'p' AND v0 = CAST(@admin_role_id AS CHAR) AND v1 = '/onlineUsers/kick' AND v2 = 'POST');

-- ---------------------------------------------------------------------------
-- 7) sys_authority_menus — sidebar visibility for default admin
-- ---------------------------------------------------------------------------
INSERT INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id)
SELECT @admin_role_id, m.id
FROM sys_base_menus m
WHERE m.name IN ('security', 'oauth2Clients', 'onlineUsers', 'sysmonitor', 'tenants')
  AND NOT EXISTS (
    SELECT 1 FROM sys_authority_menus am
    WHERE am.sys_authority_authority_id = @admin_role_id AND am.sys_base_menu_id = m.id
  );

-- =============================================================================
-- Rollback snippets (manual):
--   DELETE FROM casbin_rule WHERE ptype='p' AND v1 LIKE '/oauth2%' OR v1 LIKE '/onlineUsers%';
--   DELETE FROM sys_apis WHERE path LIKE '/oauth2%' OR path LIKE '/onlineUsers%';
--   (Menus: see security_menu_seed.sql rollback block.)
-- =============================================================================
