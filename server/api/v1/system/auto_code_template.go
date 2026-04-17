package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/huuhoait/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AutoCodeTemplateApi struct{}

// Preview
// @Tags      AutoCodeTemplate
// @Summary   Preview generated code
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.AutoCode                                      true  "Preview code generation"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Preview of generated code"
// @Router    /autoCode/preview [post]
func (a *AutoCodeTemplateApi) Preview(c *gin.Context) {
	var info request.AutoCode
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(info, utils.AutoCodeVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = info.Pretreatment()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	info.PackageT = utils.FirstUpper(info.Package)
	autoCode, err := autoCodeTemplateService.Preview(c.Request.Context(), info)
	if err != nil {
		global.GVA_LOG.Error(err.Error(), zap.Error(err))
		response.FailWithMessage("Preview failed: "+err.Error(), c)
	} else {
		response.OkWithDetailed(gin.H{"autoCode": autoCode}, "Preview successful", c)
	}
}

// Create
// @Tags      AutoCodeTemplate
// @Summary   Auto code template
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.AutoCode  true  "Create auto code"
// @Success   200   {string}  string                 "{"success":true,"data":{},"msg":"Created successfully"}"
// @Router    /autoCode/createTemp [post]
func (a *AutoCodeTemplateApi) Create(c *gin.Context) {
	var info request.AutoCode
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(info, utils.AutoCodeVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = info.Pretreatment()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = autoCodeTemplateService.Create(c.Request.Context(), info)
	if err != nil {
		global.GVA_LOG.Error("Failed to create!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("Created successfully", c)
	}
}

// AddFunc
// @Tags      AddFunc
// @Summary   Add function
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.AutoCode  true  "Add function"
// @Success   200   {string}  string                 "{"success":true,"data":{},"msg":"Created successfully"}"
// @Router    /autoCode/addFunc [post]
func (a *AutoCodeTemplateApi) AddFunc(c *gin.Context) {
	var info request.AutoFunc
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var tempMap map[string]string
	if info.IsPreview {
		info.Router = "placeholder_router"
		info.FuncName = "placeholder_funcName"
		info.Method = "placeholder_method"
		info.Description = "placeholder_description"
		tempMap, err = autoCodeTemplateService.GetApiAndServer(info)
	} else {
		err = autoCodeTemplateService.AddFunc(info)
	}
	if err != nil {
		global.GVA_LOG.Error("Failed to inject!", zap.Error(err))
		response.FailWithMessage("Injection failed", c)
	} else {
		if info.IsPreview {
			response.OkWithDetailed(tempMap, "Injected successfully", c)
			return
		}
		response.OkWithMessage("Injected successfully", c)
	}
}
