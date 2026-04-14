package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AutoCodePluginApi struct{}

// Install
// @Tags      AutoCodePlugin
// @Summary   Install plugin
// @Security  ApiKeyAuth
// @accept    multipart/form-data
// @Produce   application/json
// @Param     plug  formData  file                                              true  "this is a test file"
// @Success   200   {object}  response.Response{data=[]interface{},msg=string}  "Successfully installed plugin"
// @Router    /autoCode/installPlugin [post]
func (a *AutoCodePluginApi) Install(c *gin.Context) {
	header, err := c.FormFile("plug")
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	web, server, err := autoCodePluginService.Install(header)
	webStr := "Web plugin installed successfully"
	serverStr := "Server plugin installed successfully"
	if web == -1 {
		webStr = "Web plugin was not installed successfully. Please extract and install manually per the documentation. Ignore this if it is a server-only plugin."
	}
	if server == -1 {
		serverStr = "Server plugin was not installed successfully. Please extract and install manually per the documentation. Ignore this if it is a web-only plugin."
	}
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData([]interface{}{
		gin.H{
			"code": web,
			"msg":  webStr,
		},
		gin.H{
			"code": server,
			"msg":  serverStr,
		}}, c)
}

// Packaged
// @Tags      AutoCodePlugin
// @Summary   Package plugin
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     plugName  query    string  true  "Plugin name"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Successfully packaged plugin"
// @Router    /autoCode/pubPlug [post]
func (a *AutoCodePluginApi) Packaged(c *gin.Context) {
	plugName := c.Query("plugName")
	zipPath, err := autoCodePluginService.PubPlug(plugName)
	if err != nil {
		global.GVA_LOG.Error("Failed to package!", zap.Error(err))
		response.FailWithMessage("Packaging failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage(fmt.Sprintf("Packaged successfully, file path: %s", zipPath), c)
}

// InitMenu
// @Tags      AutoCodePlugin
// @Summary   Initialize plugin menu
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Successfully initialized plugin menu"
// @Router    /autoCode/initMenu [post]
func (a *AutoCodePluginApi) InitMenu(c *gin.Context) {
	var menuInfo request.InitMenu
	err := c.ShouldBindJSON(&menuInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = autoCodePluginService.InitMenu(menuInfo)
	if err != nil {
		global.GVA_LOG.Error("Failed to initialize Menu!", zap.Error(err))
		response.FailWithMessage("Failed to initialize Menu: "+err.Error(), c)
		return
	}
	response.OkWithMessage("File changes applied successfully", c)
}

// InitAPI
// @Tags      AutoCodePlugin
// @Summary   Initialize plugin API
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Successfully initialized plugin API"
// @Router    /autoCode/initAPI [post]
func (a *AutoCodePluginApi) InitAPI(c *gin.Context) {
	var apiInfo request.InitApi
	err := c.ShouldBindJSON(&apiInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = autoCodePluginService.InitAPI(apiInfo)
	if err != nil {
		global.GVA_LOG.Error("Failed to initialize API!", zap.Error(err))
		response.FailWithMessage("Failed to initialize API: "+err.Error(), c)
		return
	}
	response.OkWithMessage("File changes applied successfully", c)
}

// InitDictionary
// @Tags      AutoCodePlugin
// @Summary   Initialize plugin dictionary
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Successfully initialized plugin dictionary"
// @Router    /autoCode/initDictionary [post]
func (a *AutoCodePluginApi) InitDictionary(c *gin.Context) {
	var dictInfo request.InitDictionary
	err := c.ShouldBindJSON(&dictInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = autoCodePluginService.InitDictionary(dictInfo)
	if err != nil {
		global.GVA_LOG.Error("Failed to initialize Dictionary!", zap.Error(err))
		response.FailWithMessage("Failed to initialize Dictionary: "+err.Error(), c)
		return
	}
	response.OkWithMessage("File changes applied successfully", c)
}

// GetPluginList
// @Tags      AutoCodePlugin
// @Summary   Get plugin list
// @Security  ApiKeyAuth
// @Produce   application/json
// @Success   200   {object}  response.Response{data=[]systemRes.PluginInfo}  "Successfully retrieved plugin list"
// @Router    /autoCode/getPluginList [get]
func (a *AutoCodePluginApi) GetPluginList(c *gin.Context) {
	serverDir := filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Server, "plugin")
	webDir := filepath.Join(global.GVA_CONFIG.AutoCode.Root, global.GVA_CONFIG.AutoCode.Web, "plugin")

	serverEntries, _ := os.ReadDir(serverDir)
	webEntries, _ := os.ReadDir(webDir)

	configMap := make(map[string]string)

	for _, entry := range serverEntries {
		if entry.IsDir() {
			configMap[entry.Name()] = "server"
		}
	}

	for _, entry := range webEntries {
		if entry.IsDir() {
			if val, ok := configMap[entry.Name()]; ok {
				if val == "server" {
					configMap[entry.Name()] = "full"
				}
			} else {
				configMap[entry.Name()] = "web"
			}
		}
	}

	var list []systemRes.PluginInfo
	for k, v := range configMap {
		apis, menus, dicts := utils.GetPluginData(k)
		list = append(list, systemRes.PluginInfo{
			PluginName:   k,
			PluginType:   v,
			Apis:         apis,
			Menus:        menus,
			Dictionaries: dicts,
		})
	}

	response.OkWithDetailed(list, "Retrieved successfully", c)
}

// Remove
// @Tags      AutoCodePlugin
// @Summary   Remove plugin
// @Security  ApiKeyAuth
// @Produce   application/json
// @Param     pluginName  query    string  true  "Plugin name"
// @Param     pluginType  query    string  true  "Plugin type"
// @Success   200   {object}  response.Response{msg=string}  "Successfully removed plugin"
// @Router    /autoCode/removePlugin [post]
func (a *AutoCodePluginApi) Remove(c *gin.Context) {
	pluginName := c.Query("pluginName")
	pluginType := c.Query("pluginType")
	err := autoCodePluginService.Remove(pluginName, pluginType)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Deletion failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}
