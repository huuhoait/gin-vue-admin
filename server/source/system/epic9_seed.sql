-- =============================================================================
-- Epic 9: SkyAgent Admin Portal — Menu, API & Casbin seed
-- Run against the GVA admin database for EXISTING deployments.
-- Fresh deployments will auto-seed via source/system/menu.go & api.go.
-- =============================================================================

-- Step 1: Find the skyagent parent menu ID
-- (adjust if your skyagent menu has a different ID)
-- SELECT id FROM sys_base_menus WHERE name = 'skyagent';
-- Assume @skyagent_id below. Replace with actual ID.

-- ---------------------------------------------------------------------------
-- MENUS (sys_base_menus)
-- ---------------------------------------------------------------------------
-- NOTE: Replace <SKYAGENT_PARENT_ID> with the actual ID of the "skyagent" menu.
--       Run: SELECT id FROM sys_base_menus WHERE name = 'skyagent';

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
VALUES
  (NOW(), NOW(), 1, <SKYAGENT_PARENT_ID>, 'agent/list',        'agentList',        false, 'plugin/skyagent/view/agent/agentList.vue',        4,  'admin.agent.list',        'user'),
  (NOW(), NOW(), 1, <SKYAGENT_PARENT_ID>, 'agent/pending',     'agentPending',     false, 'plugin/skyagent/view/agent/pendingReview.vue',     5,  'admin.agent.pending',     'bell'),
  (NOW(), NOW(), 1, <SKYAGENT_PARENT_ID>, 'catalog/products',  'catalogProducts',  false, 'plugin/skyagent/view/catalog/productList.vue',     6,  'admin.catalog.products',  'goods'),
  (NOW(), NOW(), 1, <SKYAGENT_PARENT_ID>, 'catalog/suppliers', 'catalogSuppliers', false, 'plugin/skyagent/view/catalog/supplierList.vue',    7,  'admin.catalog.suppliers', 'van'),
  (NOW(), NOW(), 1, <SKYAGENT_PARENT_ID>, 'order/list',        'orderList',        false, 'plugin/skyagent/view/order/orderList.vue',         8,  'admin.order.list',        'document');

-- ---------------------------------------------------------------------------
-- APIs (sys_apis)
-- ---------------------------------------------------------------------------

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
VALUES
  (NOW(), NOW(), '/admin-api/v1/agents',            'List agents',             'SkyAgent Agent',     'GET'),
  (NOW(), NOW(), '/admin-api/v1/agents/:id',        'Get agent detail',        'SkyAgent Agent',     'GET'),
  (NOW(), NOW(), '/admin-api/v1/agents',            'Create agent',            'SkyAgent Agent',     'POST'),
  (NOW(), NOW(), '/admin-api/v1/agents/:id',        'Update agent',            'SkyAgent Agent',     'PUT'),
  (NOW(), NOW(), '/admin-api/v1/agents/:id/status', 'Update agent status',     'SkyAgent Agent',     'PUT'),
  (NOW(), NOW(), '/admin-api/v1/orders',            'List orders',             'SkyAgent Order',     'GET'),
  (NOW(), NOW(), '/admin-api/v1/orders/:id',        'Get order detail',        'SkyAgent Order',     'GET'),
  (NOW(), NOW(), '/admin-api/v1/products',          'List products',           'SkyAgent Catalog',   'GET'),
  (NOW(), NOW(), '/admin-api/v1/suppliers',         'List suppliers',          'SkyAgent Catalog',   'GET'),
  (NOW(), NOW(), '/admin-api/v1/dashboard/overview','Dashboard overview',      'SkyAgent Dashboard', 'GET');

-- ---------------------------------------------------------------------------
-- CASBIN RULES (casbin_rule)
-- Grant all SkyAgent APIs to admin role (authority_id = 888)
-- Adjust 888 to your admin role ID if different.
-- ---------------------------------------------------------------------------

INSERT INTO casbin_rule (ptype, v0, v1, v2)
VALUES
  ('p', '888', '/admin-api/v1/agents',            'GET'),
  ('p', '888', '/admin-api/v1/agents/:id',        'GET'),
  ('p', '888', '/admin-api/v1/agents',            'POST'),
  ('p', '888', '/admin-api/v1/agents/:id',        'PUT'),
  ('p', '888', '/admin-api/v1/agents/:id/status', 'PUT'),
  ('p', '888', '/admin-api/v1/orders',            'GET'),
  ('p', '888', '/admin-api/v1/orders/:id',        'GET'),
  ('p', '888', '/admin-api/v1/products',          'GET'),
  ('p', '888', '/admin-api/v1/suppliers',         'GET'),
  ('p', '888', '/admin-api/v1/dashboard/overview','GET');

-- ---------------------------------------------------------------------------
-- AUTHORITY-MENU MAPPING (sys_authority_menus)
-- Grant all Epic 9 menus to admin role (sys_authority_authority_id = 888)
-- NOTE: Replace <MENU_ID_*> with actual IDs from the INSERT above.
--       Or run this after INSERT and use subquery:
-- ---------------------------------------------------------------------------

INSERT INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id)
SELECT 888, id FROM sys_base_menus
WHERE name IN ('agentList', 'agentPending', 'catalogProducts', 'catalogSuppliers', 'orderList');

-- =============================================================================
-- DONE. After running:
-- 1. Restart the admin service (Casbin policies reload on boot)
-- 2. Hard-refresh the browser (Ctrl+Shift+R) to reload menus
-- =============================================================================
