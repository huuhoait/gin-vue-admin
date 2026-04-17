package system

import (
	"context"

	. "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderMenu = initOrderAuthority + 1

type initMenu struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenu, &initMenu{})
}

func (i *initMenu) InitializerName() string {
	return SysBaseMenu{}.TableName()
}

func (i *initMenu) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(
		&SysBaseMenu{},
		&SysBaseMenuParameter{},
		&SysBaseMenuBtn{},
	)
}

func (i *initMenu) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	m := db.Migrator()
	return m.HasTable(&SysBaseMenu{}) &&
		m.HasTable(&SysBaseMenuParameter{}) &&
		m.HasTable(&SysBaseMenuBtn{})
}

func (i *initMenu) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}

	// Define all menus
	allMenus := []SysBaseMenu{
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "dashboard", Name: "dashboard", Component: "view/dashboard/index.vue", Sort: 1, Meta: Meta{Title: "Dashboard", Icon: "odometer"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "about", Name: "about", Component: "view/about/index.vue", Sort: 9, Meta: Meta{Title: "About Us", Icon: "info-filled"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "admin", Name: "superAdmin", Component: "view/superAdmin/index.vue", Sort: 3, Meta: Meta{Title: "Super Admin", Icon: "user"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "skyagent", Name: "skyagent", Component: "view/routerHolder.vue", Sort: 2, Meta: Meta{Title: "admin.menu.skyagent", Icon: "grid"}},
		{MenuLevel: 0, Hidden: true, ParentId: 0, Path: "person", Name: "person", Component: "view/person/person.vue", Sort: 4, Meta: Meta{Title: "Personal Info", Icon: "message"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "example", Name: "example", Component: "view/example/index.vue", Sort: 7, Meta: Meta{Title: "Examples", Icon: "management"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "systemTools", Name: "systemTools", Component: "view/systemTools/index.vue", Sort: 5, Meta: Meta{Title: "System Tools", Icon: "tools"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "https://www.gin-vue-admin.com", Name: "https://www.gin-vue-admin.com", Component: "/", Sort: 0, Meta: Meta{Title: "Official Website", Icon: "customer-gva"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "state", Name: "state", Component: "view/system/state.vue", Sort: 8, Meta: Meta{Title: "Server Status", Icon: "cloudy"}},
		{MenuLevel: 0, Hidden: false, ParentId: 0, Path: "plugin", Name: "plugin", Component: "view/routerHolder.vue", Sort: 6, Meta: Meta{Title: "Plugin System", Icon: "cherry"}},
	}

	// Create parent menus first(ParentId = 0 ofMenu)
	if err = db.Create(&allMenus).Error; err != nil {
		return ctx, errors.Wrap(err, SysBaseMenu{}.TableName()+"parent menu initialization failed!")
	}

	// Build menu mapping - ApprovedNameLookupAlreadyCreateofMenuAnd ItsID
	menuNameMap := make(map[string]uint)
	for _, menu := range allMenus {
		menuNameMap[menu.Name] = menu.ID
	}

	// Define child menus, AndsetCorrectofParentId
	childMenus := []SysBaseMenu{
		// superAdminchild menus
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "authority", Name: "authority", Component: "view/superAdmin/authority/authority.vue", Sort: 1, Meta: Meta{Title: "Role Management", Icon: "avatar"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "menu", Name: "menu", Component: "view/superAdmin/menu/menu.vue", Sort: 2, Meta: Meta{Title: "Menu Management", Icon: "tickets", KeepAlive: true}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "api", Name: "api", Component: "view/superAdmin/api/api.vue", Sort: 3, Meta: Meta{Title: "API Management", Icon: "platform", KeepAlive: true}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "user", Name: "user", Component: "view/superAdmin/user/user.vue", Sort: 4, Meta: Meta{Title: "User Management", Icon: "coordinate"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "dictionary", Name: "dictionary", Component: "view/superAdmin/dictionary/sysDictionary.vue", Sort: 5, Meta: Meta{Title: "Dictionary Management", Icon: "notebook"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "operation", Name: "operation", Component: "view/superAdmin/operation/sysOperationRecord.vue", Sort: 6, Meta: Meta{Title: "Operation History", Icon: "pie-chart"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "sysParams", Name: "sysParams", Component: "view/superAdmin/params/sysParams.vue", Sort: 7, Meta: Meta{Title: "Parameter Management", Icon: "compass"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "system", Name: "system", Component: "view/systemTools/system/system.vue", Sort: 8, Meta: Meta{Title: "System Configuration", Icon: "operation"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "apiToken", Name: "apiToken", Component: "view/systemTools/apiToken/index.vue", Sort: 9, Meta: Meta{Title: "API Token", Icon: "key"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "loginLog", Name: "loginLog", Component: "view/systemTools/loginLog/index.vue", Sort: 10, Meta: Meta{Title: "Login Log", Icon: "monitor"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "sysVersion", Name: "sysVersion", Component: "view/systemTools/version/version.vue", Sort: 11, Meta: Meta{Title: "Version Management", Icon: "server"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["superAdmin"], Path: "sysError", Name: "sysError", Component: "view/systemTools/sysError/sysError.vue", Sort: 12, Meta: Meta{Title: "Error Log", Icon: "warn"}},

		// SkyAgent domain menus (Epic 8)
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["skyagent"], Path: "kyc", Name: "kyc", Component: "view/routerHolder.vue", Sort: 1, Meta: Meta{Title: "admin.menu.kyc", Icon: "document"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["skyagent"], Path: "kyc-cases", Name: "kycCases", Component: "view/routerHolder.vue", Sort: 2, Meta: Meta{Title: "admin.menu.kyc_cases", Icon: "list"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["skyagent"], Path: "commission", Name: "commission", Component: "view/routerHolder.vue", Sort: 3, Meta: Meta{Title: "admin.menu.commission", Icon: "coin"}},

		// SkyAgent Epic 9 — Admin Portal integration menus
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["skyagent"], Path: "agent/list", Name: "agentList", Component: "view/agent/agentList.vue", Sort: 4, Meta: Meta{Title: "admin.agent.list", Icon: "user"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["skyagent"], Path: "agent/pending", Name: "agentPending", Component: "view/agent/pendingReview.vue", Sort: 5, Meta: Meta{Title: "admin.agent.pending", Icon: "bell"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["skyagent"], Path: "catalog/products", Name: "catalogProducts", Component: "view/catalog/productList.vue", Sort: 6, Meta: Meta{Title: "admin.catalog.products", Icon: "goods"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["skyagent"], Path: "catalog/suppliers", Name: "catalogSuppliers", Component: "view/catalog/supplierList.vue", Sort: 7, Meta: Meta{Title: "admin.catalog.suppliers", Icon: "van"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["skyagent"], Path: "order/list", Name: "orderList", Component: "view/order/orderList.vue", Sort: 8, Meta: Meta{Title: "admin.order.list", Icon: "document"}},
		// Onboarding flow (Epic 11) — grouped under skyagent
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["skyagent"], Path: "onboarding/tickets", Name: "onboardingTickets", Component: "view/onboarding/ticketList.vue", Sort: 9, Meta: Meta{Title: "admin.onboarding.ticket_list", Icon: "tickets"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["skyagent"], Path: "onboarding/create", Name: "onboardingCreate", Component: "view/onboarding/createTicket.vue", Sort: 10, Meta: Meta{Title: "admin.onboarding.create_ticket", Icon: "circle-plus"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["skyagent"], Path: "onboarding/review", Name: "onboardingReview", Component: "view/onboarding/reviewQueue.vue", Sort: 11, Meta: Meta{Title: "admin.onboarding.review_queue", Icon: "checked"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["skyagent"], Path: "onboarding/agent-l0", Name: "onboardingAgentL0", Component: "view/onboarding/createAgentL0.vue", Sort: 12, Meta: Meta{Title: "admin.onboarding.create_agent_l0", Icon: "user-filled"}},

		// examplechild menus
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["example"], Path: "upload", Name: "upload", Component: "view/example/upload/upload.vue", Sort: 5, Meta: Meta{Title: "Media Library (Upload/Download)", Icon: "upload"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["example"], Path: "breakpoint", Name: "breakpoint", Component: "view/example/breakpoint/breakpoint.vue", Sort: 6, Meta: Meta{Title: "Resumable Upload", Icon: "upload-filled"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["example"], Path: "customer", Name: "customer", Component: "view/example/customer/customer.vue", Sort: 7, Meta: Meta{Title: "Customer List (Resource Example)", Icon: "avatar"}},

		// systemToolschild menus
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["systemTools"], Path: "autoPkg", Name: "autoPkg", Component: "view/systemTools/autoPkg/autoPkg.vue", Sort: 0, Meta: Meta{Title: "Template Configuration", Icon: "folder"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["systemTools"], Path: "autoCode", Name: "autoCode", Component: "view/systemTools/autoCode/index.vue", Sort: 1, Meta: Meta{Title: "Code Generator", Icon: "cpu", KeepAlive: true}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["systemTools"], Path: "autoCodeAdmin", Name: "autoCodeAdmin", Component: "view/systemTools/autoCodeAdmin/index.vue", Sort: 2, Meta: Meta{Title: "Auto Code Management", Icon: "magic-stick"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["systemTools"], Path: "formCreate", Name: "formCreate", Component: "view/systemTools/formCreate/index.vue", Sort: 3, Meta: Meta{Title: "Form Generator", Icon: "magic-stick", KeepAlive: true}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["systemTools"], Path: "aiWorkflow", Name: "aiWorkflow", Component: "view/systemTools/aiWrokflow/index.vue", Sort: 4, Meta: Meta{Title: "AIrequirement workflow", Icon: "magic-stick", KeepAlive: true}},
		{MenuLevel: 1, Hidden: true, ParentId: menuNameMap["systemTools"], Path: "autoCodeEdit/:id", Name: "autoCodeEdit", Component: "view/systemTools/autoCode/index.vue", Sort: 0, Meta: Meta{Title: "automation code-${id}", Icon: "magic-stick"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["systemTools"], Path: "exportTemplate", Name: "exportTemplate", Component: "view/systemTools/exportTemplate/exportTemplate.vue", Sort: 5, Meta: Meta{Title: "Export Template", Icon: "reading"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["systemTools"], Path: "mcpTest", Name: "mcpTest", Component: "view/systemTools/autoCode/mcpTest.vue", Sort: 6, Meta: Meta{Title: "Mcp Toolsmanagement", Icon: "partly-cloudy"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["systemTools"], Path: "mcpTool", Name: "mcpTool", Component: "view/systemTools/autoCode/mcp.vue", Sort: 7, Meta: Meta{Title: "Mcp ToolsTemplate", Icon: "magnet"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["systemTools"], Path: "skills", Name: "skills", Component: "view/systemTools/skills/index.vue", Sort: 8, Meta: Meta{Title: "Skillsmanagement", Icon: "document"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["systemTools"], Path: "picture", Name: "picture", Component: "view/systemTools/autoCode/picture.vue", Sort: 9, Meta: Meta{Title: "AIPageDraw", Icon: "picture-filled"}},

		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["plugin"], Path: "https://plugin.gin-vue-admin.com/", Name: "https://plugin.gin-vue-admin.com/", Component: "https://plugin.gin-vue-admin.com/", Sort: 0, Meta: Meta{Title: "Plugin Market", Icon: "shop"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["plugin"], Path: "installPlugin", Name: "installPlugin", Component: "view/systemTools/installPlugin/index.vue", Sort: 1, Meta: Meta{Title: "Plugin Installation", Icon: "box"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["plugin"], Path: "pubPlug", Name: "pubPlug", Component: "view/systemTools/pubPlug/pubPlug.vue", Sort: 3, Meta: Meta{Title: "Package Plugin", Icon: "files"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["plugin"], Path: "plugin-email", Name: "plugin-email", Component: "plugin/email/view/index.vue", Sort: 4, Meta: Meta{Title: "Email Plugin", Icon: "message"}},
		{MenuLevel: 1, Hidden: false, ParentId: menuNameMap["plugin"], Path: "anInfo", Name: "anInfo", Component: "plugin/announcement/view/info.vue", Sort: 5, Meta: Meta{Title: "Announcement Management [Example]", Icon: "scaleToOriginal"}},
	}

	// Create child menus
	if err = db.Create(&childMenus).Error; err != nil {
		return ctx, errors.Wrap(err, SysBaseMenu{}.TableName()+"child menu initialization failed!")
	}

	// Combine all menus as return result
	allEntities := append(allMenus, childMenus...)
	next = context.WithValue(ctx, i.InitializerName(), allEntities)
	return next, nil
}

func (i *initMenu) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ?", "autoPkg").First(&SysBaseMenu{}).Error, gorm.ErrRecordNotFound) { // Check if data exists
		return false
	}
	return true
}
