package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	common "github.com/huuhoait/gin-vue-admin/server/model/common/request"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/huuhoait/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

type AutoCodePackageApi struct{}

// Create
// @Tags      AutoCodePackage
// @Summary   Create package
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SysAutoCodePackageCreate                                         true  "Create package"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Successfully created package"
// @Router    /autoCode/createPackage [post]
func (a *AutoCodePackageApi) Create(c *gin.Context) {
	var info request.SysAutoCodePackageCreate
	_ = c.ShouldBindJSON(&info)
	if err := utils.Verify(info, utils.AutoPackageVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if strings.Contains(info.PackageName, "\\") || strings.Contains(info.PackageName, "/") || strings.Contains(info.PackageName, "..") {
		response.FailWithMessage("Invalid package name", c)
		return
	} // PackageName may cause path traversal issues, prevent both / and \
	err := autoCodePackageService.Create(c.Request.Context(), &info)
	if err != nil {
		global.GVA_LOG.Error("Failed to create!", zap.Error(err))
		response.FailWithMessage("Failed to create", c)
		return
	}
	response.OkWithMessage("Created successfully", c)
}

// Delete
// @Tags      AutoCode
// @Summary   Delete package
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      common.GetById                                         true  "Delete package"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Successfully deleted package"
// @Router    /autoCode/delPackage [post]
func (a *AutoCodePackageApi) Delete(c *gin.Context) {
	var info common.GetById
	_ = c.ShouldBindJSON(&info)
	err := autoCodePackageService.Delete(c.Request.Context(), info)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Failed to delete", c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

// All
// @Tags      AutoCodePackage
// @Summary   Get all packages
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "Successfully retrieved packages"
// @Router    /autoCode/getPackage [post]
func (a *AutoCodePackageApi) All(c *gin.Context) {
	data, err := autoCodePackageService.All(c.Request.Context())
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(gin.H{"pkgs": data}, "Retrieved successfully", c)
}

// Templates
// @Tags      AutoCodePackage
// @Summary   Get templates
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "Successfully retrieved templates"
// @Router    /autoCode/getTemplates [get]
func (a *AutoCodePackageApi) Templates(c *gin.Context) {
	data, err := autoCodePackageService.Templates(c.Request.Context())
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(data, "Retrieved successfully", c)
}
