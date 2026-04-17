package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type DBApi struct{}

// InitDB
// @Tags     InitDB
// @Summary  Initialize user database
// @Produce  application/json
// @Param    data  body      request.InitDB                  true  "Database initialization parameters"
// @Success  200   {object}  response.Response{data=string}  "Initialize user database"
// @Router   /init/initdb [post]
func (i *DBApi) InitDB(c *gin.Context) {
	if global.GVA_DB != nil {
		global.GVA_LOG.Error("Database configuration already exists!")
		response.FailWithMessage("Database configuration already exists", c)
		return
	}
	var dbInfo request.InitDB
	if err := c.ShouldBindJSON(&dbInfo); err != nil {
		global.GVA_LOG.Error("Parameter validation failed!", zap.Error(err))
		response.FailWithMessage("Parameter validation failed", c)
		return
	}
	if err := initDBService.InitDB(dbInfo); err != nil {
		global.GVA_LOG.Error("Failed to auto-create database!", zap.Error(err))
		response.FailWithMessage("Failed to auto-create database, please check the server logs and try again", c)
		return
	}
	response.OkWithMessage("Database created successfully", c)
}

// CheckDB
// @Tags     CheckDB
// @Summary  Check database initialization status
// @Produce  application/json
// @Success  200  {object}  response.Response{data=map[string]interface{},msg=string}  "Check database initialization status"
// @Router   /init/checkdb [post]
func (i *DBApi) CheckDB(c *gin.Context) {
	var (
		message  = "Please initialize the database"
		needInit = true
	)

	if global.GVA_DB != nil {
		message = "Database is already initialized"
		needInit = false
	}
	global.GVA_LOG.Info(message)
	response.OkWithDetailed(gin.H{"needInit": needInit}, message, c)
}
