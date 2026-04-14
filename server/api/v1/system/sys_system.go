package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemApi struct{}

// GetSystemConfig
// @Tags      System
// @Summary   Get configuration file content
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  response.Response{data=systemRes.SysConfigResponse,msg=string}  "Get configuration file content, returns system configuration"
// @Router    /system/getSystemConfig [post]
func (s *SystemApi) GetSystemConfig(c *gin.Context) {
	config, err := systemConfigService.GetSystemConfig()
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(systemRes.SysConfigResponse{Config: config}, "Retrieved successfully", c)
}

// SetSystemConfig
// @Tags      System
// @Summary   Set configuration file content
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     data  body      system.System                   true  "Set configuration file content"
// @Success   200   {object}  response.Response{data=string}  "Set configuration file content"
// @Router    /system/setSystemConfig [post]
func (s *SystemApi) SetSystemConfig(c *gin.Context) {
	var sys system.System
	err := c.ShouldBindJSON(&sys)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = systemConfigService.SetSystemConfig(sys)
	if err != nil {
		global.GVA_LOG.Error("Failed to set!", zap.Error(err))
		response.FailWithMessage("Failed to set", c)
		return
	}
	response.OkWithMessage("Set successfully", c)
}

// ReloadSystem
// @Tags      System
// @Summary   Reload system
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  response.Response{msg=string}  "Reload system"
// @Router    /system/reloadSystem [post]
func (s *SystemApi) ReloadSystem(c *gin.Context) {
	// Trigger system reload event
	err := utils.GlobalSystemEvents.TriggerReload()
	if err != nil {
		global.GVA_LOG.Error("Failed to reload system!", zap.Error(err))
		response.FailWithMessage("Failed to reload system: "+err.Error(), c)
		return
	}
	response.OkWithMessage("System reloaded successfully", c)
}

// GetServerInfo
// @Tags      System
// @Summary   Get server information
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "Get server information"
// @Router    /system/getServerInfo [post]
func (s *SystemApi) GetServerInfo(c *gin.Context) {
	server, err := systemConfigService.GetServerInfo()
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(gin.H{"server": server}, "Retrieved successfully", c)
}
