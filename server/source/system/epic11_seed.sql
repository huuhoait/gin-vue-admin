-- =============================================================================
-- Epic 9 + 11: Full SkyAgent Menu, API & Casbin seed
-- Run against the GVA admin database for EXISTING deployments.
-- Idempotent: safe to re-run (DELETE before INSERT).
-- =============================================================================

-- Step 0: Find skyagent parent menu ID
--   SELECT id FROM sys_base_menus WHERE name = 'skyagent';
-- Replace <SKYAGENT_ID> below with the actual value.

-- =============================================================================
-- MENUS (sys_base_menus)
-- =============================================================================

DELETE FROM sys_authority_menus WHERE sys_base_menu_id IN (
  SELECT id FROM sys_base_menus WHERE name IN (
    'agentList','agentPending','catalogProducts','catalogSuppliers',
    'orderList','onboardingTickets','onboardingCreate','onboardingReview','onboardingAgentL0'
  )
);
DELETE FROM sys_base_menus WHERE name IN (
  'agentList','agentPending','catalogProducts','catalogSuppliers',
  'orderList','onboardingTickets','onboardingCreate','onboardingReview','onboardingAgentL0'
);

INSERT INTO sys_base_menus (created_at, updated_at, menu_level, parent_id, path, name, hidden, component, sort, title, icon)
VALUES
  -- Agent Management (Epic 9)
  (NOW(), NOW(), 1, <SKYAGENT_ID>, 'agent/list',           'agentList',          false, 'view/agent/agentList.vue',            4,  'admin.agent.list',                'user'),
  (NOW(), NOW(), 1, <SKYAGENT_ID>, 'agent/pending',        'agentPending',       false, 'view/agent/pendingReview.vue',         5,  'admin.agent.pending',             'bell'),
  -- Catalog (Epic 9)
  (NOW(), NOW(), 1, <SKYAGENT_ID>, 'catalog/products',     'catalogProducts',    false, 'view/catalog/productList.vue',         6,  'admin.catalog.products',          'goods'),
  (NOW(), NOW(), 1, <SKYAGENT_ID>, 'catalog/suppliers',    'catalogSuppliers',   false, 'view/catalog/supplierList.vue',        7,  'admin.catalog.suppliers',         'van'),
  -- Orders (Epic 9)
  (NOW(), NOW(), 1, <SKYAGENT_ID>, 'order/list',           'orderList',          false, 'view/order/orderList.vue',             8,  'admin.order.list',                'document'),
  -- Onboarding Flow (Epic 11)
  (NOW(), NOW(), 1, <SKYAGENT_ID>, 'onboarding/tickets',   'onboardingTickets',  false, 'view/onboarding/ticketList.vue',       9,  'admin.onboarding.ticket_list',    'tickets'),
  (NOW(), NOW(), 1, <SKYAGENT_ID>, 'onboarding/create',    'onboardingCreate',   false, 'view/onboarding/createTicket.vue',     10, 'admin.onboarding.create_ticket',  'circle-plus'),
  (NOW(), NOW(), 1, <SKYAGENT_ID>, 'onboarding/review',    'onboardingReview',   false, 'view/onboarding/reviewQueue.vue',      11, 'admin.onboarding.review_queue',   'checked'),
  (NOW(), NOW(), 1, <SKYAGENT_ID>, 'onboarding/agent-l0',  'onboardingAgentL0',  false, 'view/onboarding/createAgentL0.vue',    12, 'admin.onboarding.create_agent_l0','user-filled');

-- =============================================================================
-- APIs (sys_apis)
-- =============================================================================

DELETE FROM sys_apis WHERE api_group LIKE 'SkyAgent%';

INSERT INTO sys_apis (created_at, updated_at, path, description, api_group, method)
VALUES
  -- Agent CRUD
  (NOW(), NOW(), '/admin-api/v1/agents',              'List agents',              'SkyAgent Agent',      'GET'),
  (NOW(), NOW(), '/admin-api/v1/agents/:id',          'Get agent detail',         'SkyAgent Agent',      'GET'),
  (NOW(), NOW(), '/admin-api/v1/agents',              'Create agent',             'SkyAgent Agent',      'POST'),
  (NOW(), NOW(), '/admin-api/v1/agents/:id',          'Update agent',             'SkyAgent Agent',      'PUT'),
  (NOW(), NOW(), '/admin-api/v1/agents/:id/status',   'Update agent status',      'SkyAgent Agent',      'PUT'),
  (NOW(), NOW(), '/admin-api/v1/agents/:id/full',     'Get agent admin PII',      'SkyAgent Agent',      'GET'),
  -- Orders
  (NOW(), NOW(), '/admin-api/v1/orders',              'List orders',              'SkyAgent Order',      'GET'),
  (NOW(), NOW(), '/admin-api/v1/orders/:id',          'Get order detail',         'SkyAgent Order',      'GET'),
  -- Catalog
  (NOW(), NOW(), '/admin-api/v1/products',            'List products',            'SkyAgent Catalog',    'GET'),
  (NOW(), NOW(), '/admin-api/v1/suppliers',           'List suppliers',           'SkyAgent Catalog',    'GET'),
  -- Dashboard
  (NOW(), NOW(), '/admin-api/v1/dashboard/overview',  'Dashboard overview',       'SkyAgent Dashboard',  'GET'),
  -- Onboarding Tickets
  (NOW(), NOW(), '/admin-api/v1/onboarding/tickets',              'List tickets',         'SkyAgent Onboarding', 'GET'),
  (NOW(), NOW(), '/admin-api/v1/onboarding/tickets',              'Create ticket',        'SkyAgent Onboarding', 'POST'),
  (NOW(), NOW(), '/admin-api/v1/onboarding/tickets/:ticket_id',   'Get ticket detail',    'SkyAgent Onboarding', 'GET'),
  (NOW(), NOW(), '/admin-api/v1/onboarding/tickets/:ticket_id/attachments', 'Upload attachment', 'SkyAgent Onboarding', 'POST'),
  (NOW(), NOW(), '/admin-api/v1/onboarding/tickets/:ticket_id/submit',      'Submit ticket',     'SkyAgent Onboarding', 'PUT'),
  (NOW(), NOW(), '/admin-api/v1/onboarding/tickets/:ticket_id/review',      'Review ticket',     'SkyAgent Onboarding', 'PUT'),
  -- Onboarding Agent L0 (Story 11.8)
  (NOW(), NOW(), '/admin-api/v1/onboarding/agents',   'Create Agent L0',          'SkyAgent Onboarding', 'POST');

-- =============================================================================
-- CASBIN RULES (casbin_rule) — role 888 (admin)
-- =============================================================================

DELETE FROM casbin_rule WHERE v1 LIKE '/admin-api/v1/%';

INSERT INTO casbin_rule (ptype, v0, v1, v2)
VALUES
  ('p', '888', '/admin-api/v1/agents',              'GET'),
  ('p', '888', '/admin-api/v1/agents/:id',          'GET'),
  ('p', '888', '/admin-api/v1/agents',              'POST'),
  ('p', '888', '/admin-api/v1/agents/:id',          'PUT'),
  ('p', '888', '/admin-api/v1/agents/:id/status',   'PUT'),
  ('p', '888', '/admin-api/v1/agents/:id/full',     'GET'),
  ('p', '888', '/admin-api/v1/orders',              'GET'),
  ('p', '888', '/admin-api/v1/orders/:id',          'GET'),
  ('p', '888', '/admin-api/v1/products',            'GET'),
  ('p', '888', '/admin-api/v1/suppliers',           'GET'),
  ('p', '888', '/admin-api/v1/dashboard/overview',  'GET'),
  ('p', '888', '/admin-api/v1/onboarding/tickets',              'GET'),
  ('p', '888', '/admin-api/v1/onboarding/tickets',              'POST'),
  ('p', '888', '/admin-api/v1/onboarding/tickets/:ticket_id',   'GET'),
  ('p', '888', '/admin-api/v1/onboarding/tickets/:ticket_id/attachments', 'POST'),
  ('p', '888', '/admin-api/v1/onboarding/tickets/:ticket_id/submit',      'PUT'),
  ('p', '888', '/admin-api/v1/onboarding/tickets/:ticket_id/review',      'PUT'),
  ('p', '888', '/admin-api/v1/onboarding/agents',   'POST');

-- =============================================================================
-- AUTHORITY-MENU MAPPING — role 888
-- =============================================================================

INSERT INTO sys_authority_menus (sys_authority_authority_id, sys_base_menu_id)
SELECT 888, id FROM sys_base_menus
WHERE name IN (
  'agentList', 'agentPending',
  'catalogProducts', 'catalogSuppliers',
  'orderList',
  'onboardingTickets', 'onboardingCreate', 'onboardingReview', 'onboardingAgentL0'
);

-- =============================================================================
-- DONE. After running:
-- 1. Replace <SKYAGENT_ID> with: SELECT id FROM sys_base_menus WHERE name = 'skyagent';
-- 2. Restart admin service (Casbin reloads on boot)
-- 3. Hard-refresh browser (Ctrl+Shift+R)
-- =============================================================================
