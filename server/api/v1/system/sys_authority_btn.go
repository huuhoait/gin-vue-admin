package system

import (
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthorityBtnApi struct{}

// GetAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   Get authority buttons
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SysAuthorityBtnReq                                      true  "Menu ID, role ID, selected button IDs"
// @Success   200   {object}  response.Response{data=response.SysAuthorityBtnRes,msg=string}  "Return list successfully"
// @Router    /authorityBtn/getAuthorityBtn [post]
func (a *AuthorityBtnApi) GetAuthorityBtn(c *gin.Context) {
	var req request.SysAuthorityBtnReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := authorityBtnService.GetAuthorityBtn(req)
	if err != nil {
		global.GVA_LOG.Error("Failed to query!", zap.Error(err))
		response.FailWithMessage("Failed to query", c)
		return
	}
	response.OkWithDetailed(res, "Query successful", c)
}

// SetAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   Set authority buttons
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.SysAuthorityBtnReq     true  "Menu ID, role ID, selected button IDs"
// @Success   200   {object}  response.Response{msg=string}  "Return list successfully"
// @Router    /authorityBtn/setAuthorityBtn [post]
func (a *AuthorityBtnApi) SetAuthorityBtn(c *gin.Context) {
	var req request.SysAuthorityBtnReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = authorityBtnService.SetAuthorityBtn(req)
	if err != nil {
		global.GVA_LOG.Error("Failed to assign!", zap.Error(err))
		response.FailWithMessage("Failed to assign", c)
		return
	}
	response.OkWithMessage("Assigned successfully", c)
}

// CanRemoveAuthorityBtn
// @Tags      AuthorityBtn
// @Summary   Check if authority button can be removed
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{msg=string}  "Deleted successfully"
// @Router    /authorityBtn/canRemoveAuthorityBtn [post]
func (a *AuthorityBtnApi) CanRemoveAuthorityBtn(c *gin.Context) {
	id := c.Query("id")
	err := authorityBtnService.CanRemoveAuthorityBtn(id)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}
