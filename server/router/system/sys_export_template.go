package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type SysExportTemplateRouter struct {
}

// InitSysExportTemplateRouter initialize Export Template route information
func (s *SysExportTemplateRouter) InitSysExportTemplateRouter(Router *gin.RouterGroup, pubRouter *gin.RouterGroup) {
	sysExportTemplateRouter := Router.Group("sysExportTemplate").Use(middleware.OperationRecord())
	sysExportTemplateRouterWithoutRecord := Router.Group("sysExportTemplate")
	sysExportTemplateRouterWithoutAuth := pubRouter.Group("sysExportTemplate")

	{
		sysExportTemplateRouter.POST("createSysExportTemplate", exportTemplateApi.CreateSysExportTemplate)             // createExport Template
		sysExportTemplateRouter.DELETE("deleteSysExportTemplate", exportTemplateApi.DeleteSysExportTemplate)           // delete export template
		sysExportTemplateRouter.DELETE("deleteSysExportTemplateByIds", exportTemplateApi.DeleteSysExportTemplateByIds) // batch delete export templates
		sysExportTemplateRouter.PUT("updateSysExportTemplate", exportTemplateApi.UpdateSysExportTemplate)              // update export template
		sysExportTemplateRouter.POST("importExcel", exportTemplateApi.ImportExcel)                                     // importexcelTemplateData
	}
	{
		sysExportTemplateRouterWithoutRecord.GET("findSysExportTemplate", exportTemplateApi.FindSysExportTemplate)       // get by IDExport Template
		sysExportTemplateRouterWithoutRecord.GET("getSysExportTemplateList", exportTemplateApi.GetSysExportTemplateList) // get export template list
		sysExportTemplateRouterWithoutRecord.GET("exportExcel", exportTemplateApi.ExportExcel)                           // getexporttoken
		sysExportTemplateRouterWithoutRecord.GET("exportTemplate", exportTemplateApi.ExportTemplate)                     // exportTableformatTemplate
        sysExportTemplateRouterWithoutRecord.GET("previewSQL", exportTemplateApi.PreviewSQL)                         // PreviewSQL
	}
	{
		sysExportTemplateRouterWithoutAuth.GET("exportExcelByToken", exportTemplateApi.ExportExcelByToken)       // ApprovedtokenexportTableformat
		sysExportTemplateRouterWithoutAuth.GET("exportTemplateByToken", exportTemplateApi.ExportTemplateByToken) // ApprovedtokenExport Template
	}
}
