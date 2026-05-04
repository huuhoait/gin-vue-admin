-- =============================================================================
-- Security menu grouping — seed for existing deployments
--
-- Creates a top-level "Security" menu and reparents the OAuth2 Clients +
-- Online Sessions plugin menus underneath it. Also normalises the menu
-- titles of all four new plugins (online-users, oauth2, sysmonitor, tenant)
-- to their i18n bundle keys so the FE sidebar resolves them via vue-i18n.
--
-- Fresh deployments do NOT need this script — the plugins' initialize/menu.go
-- now seeds i18n-keyed titles via FirstOrCreate on every server start, and
-- this script can additionally re-parent OAuth + Online Sessions under
-- Security.
--
-- Run order: AFTER the server has booted at least once (so the plugin
-- menu rows exist) and the SQL bundle from server/resource/i18n/*.toml is
-- in sync with this commit.
--
-- Idempotent: safe to re-run. INSERTs are guarded by NOT EXISTS, UPDATEs
-- are by-name and converge on the same target state.
-- =============================================================================

-- ---------------------------------------------------------------------------
-- 1. Create the "Security" parent menu (top-level, parent_id = 0)
-- ---------------------------------------------------------------------------
INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
SELECT NOW(), NOW(), 0, 0, 'security', 'security', false, 'view/routerHolder.vue', 6, 'admin.menu.security', 'lock'
WHERE NOT EXISTS (
  SELECT 1 FROM sys_base_menus WHERE name = 'security'
);

-- ---------------------------------------------------------------------------
-- 2. Re-parent OAuth2 Clients under Security and switch its title to the
--    i18n key. Matches the row by Name so it works regardless of the
--    auto-incremented ID assigned at first plugin boot.
-- ---------------------------------------------------------------------------
UPDATE sys_base_menus
SET parent_id = (SELECT id FROM (SELECT id FROM sys_base_menus WHERE name = 'security') AS s),
    title     = 'admin.plugin.oauth2.menu_title',
    sort      = 1
WHERE name = 'oauth2Clients';

-- ---------------------------------------------------------------------------
-- 3. Re-parent Online Sessions (online users) under Security.
-- ---------------------------------------------------------------------------
UPDATE sys_base_menus
SET parent_id = (SELECT id FROM (SELECT id FROM sys_base_menus WHERE name = 'security') AS s),
    title     = 'admin.plugin.online_users.menu_title',
    sort      = 2
WHERE name = 'onlineUsers';

-- ---------------------------------------------------------------------------
-- 4. Normalise titles of the other two plugin menus that stay where they
--    are (sysmonitor + tenants remain under their existing parent).
-- ---------------------------------------------------------------------------
UPDATE sys_base_menus
SET title = 'admin.plugin.sysmonitor.menu_title'
WHERE name = 'sysmonitor';

UPDATE sys_base_menus
SET title = 'admin.plugin.tenant.menu_title'
WHERE name = 'tenants';

-- ---------------------------------------------------------------------------
-- 5. (Optional) Grant the Security parent and its children to the default
--    admin role so they show up in the sidebar without manual assignment.
--    Adjust 888 to your admin authority_id if different.
-- ---------------------------------------------------------------------------
INSERT INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id)
SELECT 888, m.id
FROM sys_base_menus m
WHERE m.name IN ('security', 'oauth2Clients', 'onlineUsers', 'sysmonitor', 'tenants')
  AND NOT EXISTS (
    SELECT 1 FROM sys_authority_menus am
    WHERE am.sys_authority_authority_id = 888 AND am.sys_base_menu_id = m.id
  );

-- =============================================================================
-- Rollback (run manually if you need to revert the grouping):
--   UPDATE sys_base_menus SET parent_id = 9 WHERE name IN ('oauth2Clients','onlineUsers');
--   DELETE FROM sys_authority_menus WHERE sys_base_menu_id IN
--     (SELECT id FROM sys_base_menus WHERE name = 'security');
--   DELETE FROM sys_base_menus WHERE name = 'security';
-- =============================================================================
