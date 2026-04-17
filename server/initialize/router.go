package initialize

import (
	"net/http"
	"os"

	"github.com/huuhoait/gin-vue-admin/server/docs"
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"github.com/huuhoait/gin-vue-admin/server/router"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err == nil && stat.IsDir() {
		return nil, os.ErrPermission
	}

	return f, nil
}

// initializeTotalRoute

func Routers() *gin.Engine {
	Router := gin.New()
	// Correlation id must be the very first middleware so every log line,
	// audit row, and outbound call in this request can reference it.
	Router.Use(otelgin.Middleware("gin-vue-admin"))
	Router.Use(middleware.RequestID())
	// Story 8.3: resolve caller locale from Accept-Language for the i18n
	// response helpers. Must run before any handler that calls Ok/Fail.
	Router.Use(middleware.I18nLocale())
	// UseCustomof Recovery InIntervalpiece, Record panic AndInbound
	Router.Use(middleware.GinRecovery(true))
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}
	// CORS + CSRF must run before auth so preflights and the CSRF cookie are
	// issued even for unauthenticated requests. Pick the configured CORS
	// strategy: strict-whitelist is enforced, otherwise allow-all is used
	// (allow-all refuses to send credentials in release mode — see Cors()).
	if global.GVA_CONFIG.Cors.Mode == "strict-whitelist" {
		Router.Use(middleware.CorsByRules())
	} else {
		Router.Use(middleware.Cors())
	}
	Router.Use(middleware.CSRF())

	systemRouter := router.RouterGroupApp.System
	// IfThinkRequireNotUsenginxagentFrontendWeb Page, CanByUpdate web/.env.production Downof
	// VUE_APP_BASE_API = /
	// VUE_APP_BASE_PATH = http://localhost
	// ThenAfterExecutehitPackageCommand npm run build. AtopenDownPage3RowNoteExplain
	// Router.StaticFile("/favicon.ico", "./dist/favicon.ico")
	// Router.Static("/assets", "./dist/assets")   // distInsidePageofStaticResource
	// Router.StaticFile("/", "./dist/index.html") // FrontendWeb PageInPortPage

	Router.StaticFS(global.GVA_CONFIG.Local.StorePath, justFilesFilesystem{http.Dir(global.GVA_CONFIG.Local.StorePath)})
	// Router.Use(middleware.LoadTls())  // if neededUsehttps PleaseopenThisInIntervalpiece ThenAfterBeforeToward core/server.go WillStartMode UpdateChangeFor Router.RunTLS("port","Youofcre/pemFile","YouofkeyFile")
	// CORS, Such AsNeedCORSCanByopenDownPageofNoteExplain
	// Router.Use(middleware.Cors()) // DirectPlaceRowAllCORSRequest
	// Router.Use(middleware.CorsByRules()) // According toconfigurationofRulePlaceRowCORSRequest
	// global.GVA_LOG.Info("use middleware cors")
	docs.SwaggerInfo.BasePath = global.GVA_CONFIG.System.RouterPrefix
	Router.GET(global.GVA_CONFIG.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.GVA_LOG.Info("register swagger handler")
	// ConvenientUnifiedAddroute groupPrefix MultipleServerUpperLineUse

	PublicGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)
	PrivateGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)

	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())

	// Detailed liveness + readiness health probes
	RegisterHealthRoutes(Router, global.GVA_CONFIG.System.RouterPrefix)

	{
		// HealthMonitorTest — backward-compat alias kept
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
	{
		systemRouter.InitBaseRouter(PublicGroup)               // RegisterBasicFunctionRoute NotDoAuthorization
		systemRouter.InitInitRouter(PublicGroup)               // AutomaticinitializeRelated
	}

	{
		systemRouter.InitApiRouter(PrivateGroup, PublicGroup)                 // RegisterFunctionapiRoute
		systemRouter.InitJwtRouter(PrivateGroup)                              // jwtrelated routes
		systemRouter.InitUserRouter(PrivateGroup)                             // RegisterUserRoute
		systemRouter.InitMenuRouter(PrivateGroup)                             // RegistermenuRoute
		systemRouter.InitSystemRouter(PrivateGroup)                           // systemrelated routes
		systemRouter.InitSysVersionRouter(PrivateGroup)                       // releaserelated routes
		systemRouter.InitCasbinRouter(PrivateGroup)                           // Permissionrelated routes
		systemRouter.InitAutoCodeRouter(PrivateGroup, PublicGroup)             // Createautomation code
		systemRouter.InitAuthorityRouter(PrivateGroup)                        // RegisterRoleRoute
		systemRouter.InitSysDictionaryRouter(PrivateGroup)                    // Dictionary Management
		systemRouter.InitAutoCodeHistoryRouter(PrivateGroup)                  // automation codeHistory
		systemRouter.InitSysOperationRecordRouter(PrivateGroup)               // Operation Record
		systemRouter.InitSysDictionaryDetailRouter(PrivateGroup)              // DictionaryDetailsmanagement
		systemRouter.InitAuthorityBtnRouterRouter(PrivateGroup)               // Button Permissionmanagement
		systemRouter.InitSysExportTemplateRouter(PrivateGroup, PublicGroup)    // Export Template
		systemRouter.InitSysParamsRouter(PrivateGroup, PublicGroup)            // Parameter Management
		systemRouter.InitSysErrorRouter(PrivateGroup, PublicGroup)             // Error Log
		systemRouter.InitLoginLogRouter(PrivateGroup)                         // Login Log
		systemRouter.InitApiTokenRouter(PrivateGroup)                         // apiTokenSignsend
		systemRouter.InitSkillsRouter(PrivateGroup, PublicGroup)              // Skills DefineDevice
		systemRouter.InitAuditRouter(PrivateGroup)                            // Audit chain verification
	}

	// SkyAgent domain placeholder endpoints (Epic 8)
	{
		PrivateGroup.GET("/kyc/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
		PrivateGroup.GET("/commission/ping", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })
	}

	// SkyAgent BFF proxy routes (Epic 9)
	{
		proxyRouter := router.RouterGroupApp.Proxy
		skyAgentGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix + "/admin-api/v1")
		skyAgentGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
		proxyRouter.InitSkyAgentRouter(skyAgentGroup)
	}

	//PluginRoutesafeInstall
	InstallPlugin(PrivateGroup, PublicGroup, Router)

	// RegisterBusinessRoute
	initBizRouter(PrivateGroup, PublicGroup)

	global.GVA_ROUTERS = Router.Routes()

	global.GVA_LOG.Info("router register success")
	return Router
}
