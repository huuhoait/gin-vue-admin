package system

import (
	"context"

	sysModel "github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/huuhoait/gin-vue-admin/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type initApi struct{}

const initOrderApi = system.InitOrderSystem + 1

// auto run
func init() {
	system.RegisterInit(initOrderApi, &initApi{})
}

func (i *initApi) InitializerName() string {
	return sysModel.SysApi{}.TableName()
}

func (i *initApi) MigrateTable(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	return ctx, db.AutoMigrate(&sysModel.SysApi{})
}

func (i *initApi) TableCreated(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	return db.Migrator().HasTable(&sysModel.SysApi{})
}

func (i *initApi) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	entities := []sysModel.SysApi{
		{ApiGroup: "jwt", Method: "POST", Path: "/jwt/jsonInBlacklist", Description: "add JWT to blacklist(Logout, required)"},

		{ApiGroup: "Login Log", Method: "DELETE", Path: "/sysLoginLog/deleteLoginLog", Description: "delete login log"},
		{ApiGroup: "Login Log", Method: "DELETE", Path: "/sysLoginLog/deleteLoginLogByIds", Description: "batch delete login logs"},
		{ApiGroup: "Login Log", Method: "GET", Path: "/sysLoginLog/findLoginLog", Description: "get by IDLogin Log"},
		{ApiGroup: "Login Log", Method: "GET", Path: "/sysLoginLog/getLoginLogList", Description: "get login log list"},

		{ApiGroup: "API Token", Method: "POST", Path: "/sysApiToken/createApiToken", Description: "SignsendAPI Token"},
		{ApiGroup: "API Token", Method: "POST", Path: "/sysApiToken/getApiTokenList", Description: "getAPI TokenList"},
		{ApiGroup: "API Token", Method: "POST", Path: "/sysApiToken/deleteApiToken", Description: "VoidAPI Token"},

		{ApiGroup: "System User", Method: "DELETE", Path: "/user/deleteUser", Description: "delete user"},
		{ApiGroup: "System User", Method: "POST", Path: "/user/admin_register", Description: "user register"},
		{ApiGroup: "System User", Method: "POST", Path: "/user/getUserList", Description: "get user list"},
		{ApiGroup: "System User", Method: "PUT", Path: "/user/setUserInfo", Description: "set user info"},
		{ApiGroup: "System User", Method: "PUT", Path: "/user/setSelfInfo", Description: "setown profile(required)"},
		{ApiGroup: "System User", Method: "GET", Path: "/user/getUserInfo", Description: "getown profile(required)"},
		{ApiGroup: "System User", Method: "POST", Path: "/user/setUserAuthorities", Description: "setPermissionGroup"},
		{ApiGroup: "System User", Method: "POST", Path: "/user/changePassword", Description: "Updatepassword(recommended selection)"},
		{ApiGroup: "System User", Method: "POST", Path: "/user/setUserAuthority", Description: "UpdateUserRole(required)"},
		{ApiGroup: "System User", Method: "POST", Path: "/user/resetPassword", Description: "reset user password"},
		{ApiGroup: "System User", Method: "PUT", Path: "/user/setSelfSetting", Description: "user interfaceconfiguration"},

		{ApiGroup: "api", Method: "POST", Path: "/api/createApi", Description: "Createapi"},
		{ApiGroup: "api", Method: "POST", Path: "/api/deleteApi", Description: "deleteApi"},
		{ApiGroup: "api", Method: "POST", Path: "/api/updateApi", Description: "updateApi"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getApiList", Description: "getapiList"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getAllApis", Description: "get all APIs"},
		{ApiGroup: "api", Method: "POST", Path: "/api/getApiById", Description: "getapiDetailsInformation"},
		{ApiGroup: "api", Method: "DELETE", Path: "/api/deleteApisByIds", Description: "Batchdelete API"},
		{ApiGroup: "api", Method: "GET", Path: "/api/syncApi", Description: "getPendingsync APIs"},
		{ApiGroup: "api", Method: "GET", Path: "/api/getApiGroups", Description: "getroute group"},
		{ApiGroup: "api", Method: "POST", Path: "/api/enterSyncApi", Description: "sureAcknowledgesync APIs"},
		{ApiGroup: "api", Method: "POST", Path: "/api/ignoreApi", Description: "IgnoreAPI"},

		{ApiGroup: "Role", Method: "POST", Path: "/authority/copyAuthority", Description: "copy role"},
		{ApiGroup: "Role", Method: "POST", Path: "/authority/createAuthority", Description: "create role"},
		{ApiGroup: "Role", Method: "POST", Path: "/authority/deleteAuthority", Description: "delete role"},
		{ApiGroup: "Role", Method: "PUT", Path: "/authority/updateAuthority", Description: "update role info"},
		{ApiGroup: "Role", Method: "POST", Path: "/authority/getAuthorityList", Description: "get role list"},
		{ApiGroup: "Role", Method: "POST", Path: "/authority/setDataAuthority", Description: "set role resource permissions"},
		{ApiGroup: "Role", Method: "GET", Path: "/authority/getUsersByAuthority", Description: "getRoleassociated userIDList"},
		{ApiGroup: "Role", Method: "POST", Path: "/authority/setRoleUsers", Description: "full overwriteRoleassociated user"},

		{ApiGroup: "casbin", Method: "POST", Path: "/casbin/updateCasbin", Description: "UpdateChangeRoleapiPermission"},
		{ApiGroup: "casbin", Method: "POST", Path: "/casbin/getPolicyPathByAuthorityId", Description: "get permissions list"},

		{ApiGroup: "Menu", Method: "POST", Path: "/menu/addBaseMenu", Description: "create menu"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/getMenu", Description: "get menu tree(required)"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/deleteBaseMenu", Description: "delete menu"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/updateBaseMenu", Description: "update menu"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/getBaseMenuById", Description: "According toidgetMenu"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/getMenuList", Description: "PaginationgetBasicmenuList"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/getBaseMenuTree", Description: "getuser dynamic routes"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/getMenuAuthority", Description: "get menus by role"},
		{ApiGroup: "Menu", Method: "POST", Path: "/menu/addMenuAuthority", Description: "add menu-role association"},

		{ApiGroup: "Chunked Upload", Method: "GET", Path: "/fileUploadAndDownload/findFile", Description: "SearchFindItemTagFile(SecondTransmit)"},
		{ApiGroup: "Chunked Upload", Method: "POST", Path: "/fileUploadAndDownload/breakpointContinue", Description: "Resumable Upload"},
		{ApiGroup: "Chunked Upload", Method: "POST", Path: "/fileUploadAndDownload/breakpointContinueFinish", Description: "Resumable UploadComplete"},
		{ApiGroup: "Chunked Upload", Method: "POST", Path: "/fileUploadAndDownload/removeChunk", Description: "UploadCompleteRemoveFile"},

		{ApiGroup: "File Upload & Download", Method: "POST", Path: "/fileUploadAndDownload/upload", Description: "FileUpload(recommended selection)"},
		{ApiGroup: "File Upload & Download", Method: "POST", Path: "/fileUploadAndDownload/deleteFile", Description: "deleteFile"},
		{ApiGroup: "File Upload & Download", Method: "POST", Path: "/fileUploadAndDownload/editFileName", Description: "file nameor remarkedit"},
		{ApiGroup: "File Upload & Download", Method: "POST", Path: "/fileUploadAndDownload/getFileList", Description: "getupload fileList"},
		{ApiGroup: "File Upload & Download", Method: "POST", Path: "/fileUploadAndDownload/importURL", Description: "importURL"},

		{ApiGroup: "System Service", Method: "POST", Path: "/system/getServerInfo", Description: "get server info"},
		{ApiGroup: "System Service", Method: "POST", Path: "/system/getSystemConfig", Description: "getconfigurationFilecontent"},
		{ApiGroup: "System Service", Method: "POST", Path: "/system/setSystemConfig", Description: "setconfigurationFilecontent"},

		{ApiGroup: "skills", Method: "GET", Path: "/skills/getTools", Description: "getSkillAbleUtilityList"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/getSkillList", Description: "getSkillAbleList"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/getSkillDetail", Description: "getSkillAbleDetails"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/saveSkill", Description: "SaveSkillAbleDefine"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/deleteSkill", Description: "deleteSkillAble"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/createScript", Description: "CreateSkillAblefootThis"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/getScript", Description: "ReadSkillAblefootThis"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/saveScript", Description: "SaveSkillAblefootThis"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/createResource", Description: "CreateSkillAbleResource"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/getResource", Description: "ReadSkillAbleResource"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/saveResource", Description: "SaveSkillAbleResource"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/createReference", Description: "CreateSkillAbleReference"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/getReference", Description: "ReadSkillAbleReference"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/saveReference", Description: "SaveSkillAbleReference"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/createTemplate", Description: "CreateSkillAbleTemplate"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/getTemplate", Description: "ReadSkillAbleTemplate"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/saveTemplate", Description: "SaveSkillAbleTemplate"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/getGlobalConstraint", Description: "ReadGlobalconstraint"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/saveGlobalConstraint", Description: "SaveGlobalconstraint"},
		{ApiGroup: "skills", Method: "POST", Path: "/skills/packageSkill", Description: "hitPackageSkillAble"},

		{ApiGroup: "Customer", Method: "PUT", Path: "/customer/customer", Description: "updateCustomer"},
		{ApiGroup: "Customer", Method: "POST", Path: "/customer/customer", Description: "CreateCustomer"},
		{ApiGroup: "Customer", Method: "DELETE", Path: "/customer/customer", Description: "deleteCustomer"},
		{ApiGroup: "Customer", Method: "GET", Path: "/customer/customer", Description: "getDocumentOneCustomer"},
		{ApiGroup: "Customer", Method: "GET", Path: "/customer/customerList", Description: "getCustomerList"},

		{ApiGroup: "Code Generator", Method: "GET", Path: "/autoCode/getDB", Description: "getAllDatabase"},
		{ApiGroup: "Code Generator", Method: "GET", Path: "/autoCode/getTables", Description: "getDatabaseTable"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/createTemp", Description: "automation code"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/preview", Description: "Previewautomation code"},
		{ApiGroup: "Code Generator", Method: "GET", Path: "/autoCode/getColumn", Description: "getofChoosetableofAllField"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/installPlugin", Description: "install plugin"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/pubPlug", Description: "Package Plugin"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/removePlugin", Description: "unloadCarryPlugin"},
		{ApiGroup: "Code Generator", Method: "GET", Path: "/autoCode/getPluginList", Description: "getAlreadyinstall plugin"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/mcp", Description: "auto-generate MCP Tool Template"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/mcpStatus", Description: "get MCP standalone servicestatus"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/mcpStart", Description: "Start MCP standalone service"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/mcpStop", Description: "Disabled MCP standalone service"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/mcpTest", Description: "MCP Tool management"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/mcpList", Description: "get MCP ToolList"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/saveAIWorkflowSession", Description: "SaveAIrequirement workflow session"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/getAIWorkflowSessionList", Description: "getAIrequirement workflow sessionList"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/getAIWorkflowSessionDetail", Description: "getAIrequirement workflow sessionDetails"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/deleteAIWorkflowSession", Description: "deleteAIrequirement workflow session"},
		{ApiGroup: "Code Generator", Method: "POST", Path: "/autoCode/dumpAIWorkflowMarkdown", Description: "AIrequirement workflowMarkdownDropdisk"},

		{ApiGroup: "Template Configuration", Method: "POST", Path: "/autoCode/createPackage", Description: "configurationTemplate"},
		{ApiGroup: "Template Configuration", Method: "GET", Path: "/autoCode/getTemplates", Description: "getTemplateFile"},
		{ApiGroup: "Template Configuration", Method: "POST", Path: "/autoCode/getPackage", Description: "getAllTemplate"},
		{ApiGroup: "Template Configuration", Method: "POST", Path: "/autoCode/delPackage", Description: "deleteTemplate"},

		{ApiGroup: "Code Generator History", Method: "POST", Path: "/autoCode/getMeta", Description: "getmetaInformation"},
		{ApiGroup: "Code Generator History", Method: "POST", Path: "/autoCode/rollback", Description: "rollbackauto-generate code"},
		{ApiGroup: "Code Generator History", Method: "POST", Path: "/autoCode/getSysHistory", Description: "QueryrollbackRecord"},
		{ApiGroup: "Code Generator History", Method: "POST", Path: "/autoCode/delSysHistory", Description: "deleterollbackRecord"},
		{ApiGroup: "Code Generator History", Method: "POST", Path: "/autoCode/addFunc", Description: "IncreaseTemplatemethod"},

		{ApiGroup: "Dictionary Detail", Method: "PUT", Path: "/sysDictionaryDetail/updateSysDictionaryDetail", Description: "update dictionarycontent"},
		{ApiGroup: "Dictionary Detail", Method: "POST", Path: "/sysDictionaryDetail/createSysDictionaryDetail", Description: "new dictionarycontent"},
		{ApiGroup: "Dictionary Detail", Method: "DELETE", Path: "/sysDictionaryDetail/deleteSysDictionaryDetail", Description: "delete dictionarycontent"},
		{ApiGroup: "Dictionary Detail", Method: "GET", Path: "/sysDictionaryDetail/findSysDictionaryDetail", Description: "get by ID dictionarycontent"},
		{ApiGroup: "Dictionary Detail", Method: "GET", Path: "/sysDictionaryDetail/getSysDictionaryDetailList", Description: "get dictionarycontentList"},

		{ApiGroup: "Dictionary Detail", Method: "GET", Path: "/sysDictionaryDetail/getDictionaryTreeList", Description: "get dictionarycount list"},
		{ApiGroup: "Dictionary Detail", Method: "GET", Path: "/sysDictionaryDetail/getDictionaryTreeListByType", Description: "According toPartClassget dictionarycount list"},
		{ApiGroup: "Dictionary Detail", Method: "GET", Path: "/sysDictionaryDetail/getDictionaryDetailsByParent", Description: "according to parentIDget dictionaryDetails"},
		{ApiGroup: "Dictionary Detail", Method: "GET", Path: "/sysDictionaryDetail/getDictionaryPath", Description: "get dictionaryfull detailpath"},

		{ApiGroup: "Dictionary", Method: "POST", Path: "/sysDictionary/createSysDictionary", Description: "new dictionary"},
		{ApiGroup: "Dictionary", Method: "DELETE", Path: "/sysDictionary/deleteSysDictionary", Description: "delete dictionary"},
		{ApiGroup: "Dictionary", Method: "PUT", Path: "/sysDictionary/updateSysDictionary", Description: "update dictionary"},
		{ApiGroup: "Dictionary", Method: "GET", Path: "/sysDictionary/findSysDictionary", Description: "get by ID dictionary(recommended selection)"},
		{ApiGroup: "Dictionary", Method: "GET", Path: "/sysDictionary/getSysDictionaryList", Description: "get dictionary list"},
		{ApiGroup: "Dictionary", Method: "POST", Path: "/sysDictionary/importSysDictionary", Description: "importDictionaryJSON"},
		{ApiGroup: "Dictionary", Method: "GET", Path: "/sysDictionary/exportSysDictionary", Description: "exportDictionaryJSON"},

		{ApiGroup: "Operation Record", Method: "POST", Path: "/sysOperationRecord/createSysOperationRecord", Description: "createOperation Record"},
		{ApiGroup: "Operation Record", Method: "GET", Path: "/sysOperationRecord/findSysOperationRecord", Description: "get by IDOperation Record"},
		{ApiGroup: "Operation Record", Method: "GET", Path: "/sysOperationRecord/getSysOperationRecordList", Description: "getOperation RecordList"},
		{ApiGroup: "Operation Record", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecord", Description: "deleteOperation Record"},
		{ApiGroup: "Operation Record", Method: "DELETE", Path: "/sysOperationRecord/deleteSysOperationRecordByIds", Description: "batch deleteOperation History"},

		{ApiGroup: "Resumable Upload(plugin edition)", Method: "POST", Path: "/simpleUploader/upload", Description: "plugin editionChunked Upload"},
		{ApiGroup: "Resumable Upload(plugin edition)", Method: "GET", Path: "/simpleUploader/checkFileMd5", Description: "FileCompleteArrangedegreeVerify"},
		{ApiGroup: "Resumable Upload(plugin edition)", Method: "GET", Path: "/simpleUploader/mergeFileMd5", Description: "UploadCompleteCombineAndFile"},

		{ApiGroup: "email", Method: "POST", Path: "/email/emailTest", Description: "send test email"},
		{ApiGroup: "email", Method: "POST", Path: "/email/sendEmail", Description: "send email"},

		{ApiGroup: "Button Permission", Method: "POST", Path: "/authorityBtn/setAuthorityBtn", Description: "set button permissions"},
		{ApiGroup: "Button Permission", Method: "POST", Path: "/authorityBtn/getAuthorityBtn", Description: "getAlreadyHaveButton Permission"},
		{ApiGroup: "Button Permission", Method: "POST", Path: "/authorityBtn/canRemoveAuthorityBtn", Description: "deleteButton"},

		{ApiGroup: "Export Template", Method: "POST", Path: "/sysExportTemplate/createSysExportTemplate", Description: "createExport Template"},
		{ApiGroup: "Export Template", Method: "DELETE", Path: "/sysExportTemplate/deleteSysExportTemplate", Description: "delete export template"},
		{ApiGroup: "Export Template", Method: "DELETE", Path: "/sysExportTemplate/deleteSysExportTemplateByIds", Description: "batch delete export templates"},
		{ApiGroup: "Export Template", Method: "PUT", Path: "/sysExportTemplate/updateSysExportTemplate", Description: "update export template"},
		{ApiGroup: "Export Template", Method: "GET", Path: "/sysExportTemplate/findSysExportTemplate", Description: "get by IDExport Template"},
		{ApiGroup: "Export Template", Method: "GET", Path: "/sysExportTemplate/getSysExportTemplateList", Description: "get export template list"},
		{ApiGroup: "Export Template", Method: "GET", Path: "/sysExportTemplate/exportExcel", Description: "export Excel"},
		{ApiGroup: "Export Template", Method: "GET", Path: "/sysExportTemplate/exportTemplate", Description: "DownloadTemplate"},
		{ApiGroup: "Export Template", Method: "GET", Path: "/sysExportTemplate/previewSQL", Description: "PreviewSQL"},
		{ApiGroup: "Export Template", Method: "POST", Path: "/sysExportTemplate/importExcel", Description: "import Excel"},

		{ApiGroup: "Error Log", Method: "POST", Path: "/sysError/createSysError", Description: "createError Log"},
		{ApiGroup: "Error Log", Method: "DELETE", Path: "/sysError/deleteSysError", Description: "deleteError Log"},
		{ApiGroup: "Error Log", Method: "DELETE", Path: "/sysError/deleteSysErrorByIds", Description: "batch delete errorsLog"},
		{ApiGroup: "Error Log", Method: "PUT", Path: "/sysError/updateSysError", Description: "updateError Log"},
		{ApiGroup: "Error Log", Method: "GET", Path: "/sysError/findSysError", Description: "get by IDError Log"},
		{ApiGroup: "Error Log", Method: "GET", Path: "/sysError/getSysErrorList", Description: "getError LogList"},
		{ApiGroup: "Error Log", Method: "GET", Path: "/sysError/getSysErrorSolution", Description: "TriggerErrorHandle(Asynchronous)"},

		{ApiGroup: "Announcement", Method: "POST", Path: "/info/createInfo", Description: "createAnnouncement"},
		{ApiGroup: "Announcement", Method: "DELETE", Path: "/info/deleteInfo", Description: "delete announcement"},
		{ApiGroup: "Announcement", Method: "DELETE", Path: "/info/deleteInfoByIds", Description: "batch delete announcements"},
		{ApiGroup: "Announcement", Method: "PUT", Path: "/info/updateInfo", Description: "update announcement"},
		{ApiGroup: "Announcement", Method: "GET", Path: "/info/findInfo", Description: "get by IDAnnouncement"},
		{ApiGroup: "Announcement", Method: "GET", Path: "/info/getInfoList", Description: "getAnnouncementList"},

		{ApiGroup: "Parameter Management", Method: "POST", Path: "/sysParams/createSysParams", Description: "createParameter"},
		{ApiGroup: "Parameter Management", Method: "DELETE", Path: "/sysParams/deleteSysParams", Description: "delete parameter"},
		{ApiGroup: "Parameter Management", Method: "DELETE", Path: "/sysParams/deleteSysParamsByIds", Description: "batch delete parameters"},
		{ApiGroup: "Parameter Management", Method: "PUT", Path: "/sysParams/updateSysParams", Description: "update parameter"},
		{ApiGroup: "Parameter Management", Method: "GET", Path: "/sysParams/findSysParams", Description: "get by IDParameter"},
		{ApiGroup: "Parameter Management", Method: "GET", Path: "/sysParams/getSysParamsList", Description: "get parameter list"},
		{ApiGroup: "Parameter Management", Method: "GET", Path: "/sysParams/getSysParam", Description: "get parameter list"},
		{ApiGroup: "Media Category", Method: "GET", Path: "/attachmentCategory/getCategoryList", Description: "category list"},
		{ApiGroup: "Media Category", Method: "POST", Path: "/attachmentCategory/addCategory", Description: "Add/editPartClass"},
		{ApiGroup: "Media Category", Method: "POST", Path: "/attachmentCategory/deleteCategory", Description: "deletePartClass"},

		{ApiGroup: "Version Control", Method: "GET", Path: "/sysVersion/findSysVersion", Description: "getDocumentOneVersion"},
		{ApiGroup: "Version Control", Method: "GET", Path: "/sysVersion/getSysVersionList", Description: "get version list"},
		{ApiGroup: "Version Control", Method: "GET", Path: "/sysVersion/downloadVersionJson", Description: "download versionjson"},
		{ApiGroup: "Version Control", Method: "POST", Path: "/sysVersion/exportVersion", Description: "create version"},
		{ApiGroup: "Version Control", Method: "POST", Path: "/sysVersion/importVersion", Description: "SynchronousVersion"},
		{ApiGroup: "Version Control", Method: "DELETE", Path: "/sysVersion/deleteSysVersion", Description: "delete version"},
		{ApiGroup: "Version Control", Method: "DELETE", Path: "/sysVersion/deleteSysVersionByIds", Description: "batch delete versions"},

		// SkyAgent BFF Proxy APIs (Epic 9)
		{ApiGroup: "SkyAgent Agent", Method: "GET", Path: "/admin-api/v1/agents", Description: "List agents"},
		{ApiGroup: "SkyAgent Agent", Method: "GET", Path: "/admin-api/v1/agents/:id", Description: "Get agent detail"},
		{ApiGroup: "SkyAgent Agent", Method: "POST", Path: "/admin-api/v1/agents", Description: "Create agent"},
		{ApiGroup: "SkyAgent Agent", Method: "PUT", Path: "/admin-api/v1/agents/:id", Description: "Update agent"},
		{ApiGroup: "SkyAgent Agent", Method: "PUT", Path: "/admin-api/v1/agents/:id/status", Description: "Update agent status"},
		{ApiGroup: "SkyAgent Order", Method: "GET", Path: "/admin-api/v1/orders", Description: "List orders"},
		{ApiGroup: "SkyAgent Order", Method: "GET", Path: "/admin-api/v1/orders/:id", Description: "Get order detail"},
		{ApiGroup: "SkyAgent Catalog", Method: "GET", Path: "/admin-api/v1/products", Description: "List products"},
		{ApiGroup: "SkyAgent Catalog", Method: "GET", Path: "/admin-api/v1/suppliers", Description: "List suppliers"},
		{ApiGroup: "SkyAgent Dashboard", Method: "GET", Path: "/admin-api/v1/dashboard/overview", Description: "Dashboard overview"},
	}
	if err := db.Create(&entities).Error; err != nil {
		return ctx, errors.Wrap(err, sysModel.SysApi{}.TableName()+"table data initialization failed!")
	}
	next := context.WithValue(ctx, i.InitializerName(), entities)
	return next, nil
}

func (i *initApi) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ? AND method = ?", "/authorityBtn/canRemoveAuthorityBtn", "POST").
		First(&sysModel.SysApi{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
