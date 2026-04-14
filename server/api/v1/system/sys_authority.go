package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"
	systemRes "github.com/flipped-aurora/gin-vue-admin/server/model/system/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthorityApi struct{}

// CreateAuthority
// @Tags      Authority
// @Summary   Create role
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysAuthority                                                true  "Authority ID, authority name, parent role ID"
// @Success   200   {object}  response.Response{data=systemRes.SysAuthorityResponse,msg=string}  "Create role, returns system role details"
// @Router    /authority/createAuthority [post]
func (a *AuthorityApi) CreateAuthority(c *gin.Context) {
	var authority, authBack system.SysAuthority
	var err error

	if err = c.ShouldBindJSON(&authority); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err = utils.Verify(authority, utils.AuthorityVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if *authority.ParentId == 0 && global.GVA_CONFIG.System.UseStrictAuth {
		authority.ParentId = utils.Pointer(utils.GetUserAuthorityId(c))
	}

	if authBack, err = authorityService.CreateAuthority(authority); err != nil {
		global.GVA_LOG.Error("Failed to create!", zap.Error(err))
		response.FailWithMessage("Creation failed: "+err.Error(), c)
		return
	}
	err = casbinService.FreshCasbin()
	if err != nil {
		global.GVA_LOG.Error("Created successfully, but failed to refresh permissions.", zap.Error(err))
		response.FailWithMessage("Created successfully, but failed to refresh permissions: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(systemRes.SysAuthorityResponse{Authority: authBack}, "Created successfully", c)
}

// CopyAuthority
// @Tags      Authority
// @Summary   Copy role
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      response.SysAuthorityCopyResponse                                  true  "Old role ID, new authority ID, new authority name, new parent role ID"
// @Success   200   {object}  response.Response{data=systemRes.SysAuthorityResponse,msg=string}  "Copy role, returns system role details"
// @Router    /authority/copyAuthority [post]
func (a *AuthorityApi) CopyAuthority(c *gin.Context) {
	var copyInfo systemRes.SysAuthorityCopyResponse
	err := c.ShouldBindJSON(&copyInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(copyInfo, utils.OldAuthorityVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(copyInfo.Authority, utils.AuthorityVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	adminAuthorityID := utils.GetUserAuthorityId(c)
	authBack, err := authorityService.CopyAuthority(adminAuthorityID, copyInfo)
	if err != nil {
		global.GVA_LOG.Error("Failed to copy!", zap.Error(err))
		response.FailWithMessage("Copy failed: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(systemRes.SysAuthorityResponse{Authority: authBack}, "Copied successfully", c)
}

// DeleteAuthority
// @Tags      Authority
// @Summary   Delete role
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysAuthority            true  "Delete role"
// @Success   200   {object}  response.Response{msg=string}  "Delete role"
// @Router    /authority/deleteAuthority [post]
func (a *AuthorityApi) DeleteAuthority(c *gin.Context) {
	var authority system.SysAuthority
	var err error
	if err = c.ShouldBindJSON(&authority); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err = utils.Verify(authority, utils.AuthorityIdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// Check if any user is using this role before deleting
	if err = authorityService.DeleteAuthority(&authority); err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Deletion failed: "+err.Error(), c)
		return
	}
	_ = casbinService.FreshCasbin()
	response.OkWithMessage("Deleted successfully", c)
}

// UpdateAuthority
// @Tags      Authority
// @Summary   Update role information
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysAuthority                                                true  "Authority ID, authority name, parent role ID"
// @Success   200   {object}  response.Response{data=systemRes.SysAuthorityResponse,msg=string}  "Update role information, returns system role details"
// @Router    /authority/updateAuthority [put]
func (a *AuthorityApi) UpdateAuthority(c *gin.Context) {
	var auth system.SysAuthority
	err := c.ShouldBindJSON(&auth)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(auth, utils.AuthorityVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	authority, err := authorityService.UpdateAuthority(auth)
	if err != nil {
		global.GVA_LOG.Error("Failed to update!", zap.Error(err))
		response.FailWithMessage("Update failed: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(systemRes.SysAuthorityResponse{Authority: authority}, "Updated successfully", c)
}

// GetAuthorityList
// @Tags      Authority
// @Summary   Get role list with pagination
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.PageInfo                                        true  "Page number, page size"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "Get paginated role list, returns list, total, page, page size"
// @Router    /authority/getAuthorityList [post]
func (a *AuthorityApi) GetAuthorityList(c *gin.Context) {
	authorityID := utils.GetUserAuthorityId(c)
	list, err := authorityService.GetAuthorityInfoList(authorityID)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Retrieval failed: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(list, "Retrieved successfully", c)
}

// SetDataAuthority
// @Tags      Authority
// @Summary   Set role data authority
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysAuthority            true  "Set role data authority"
// @Success   200   {object}  response.Response{msg=string}  "Set role data authority"
// @Router    /authority/setDataAuthority [post]
func (a *AuthorityApi) SetDataAuthority(c *gin.Context) {
	var auth system.SysAuthority
	err := c.ShouldBindJSON(&auth)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(auth, utils.AuthorityIdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	adminAuthorityID := utils.GetUserAuthorityId(c)
	err = authorityService.SetDataAuthority(adminAuthorityID, auth)
	if err != nil {
		global.GVA_LOG.Error("Failed to set!", zap.Error(err))
		response.FailWithMessage("Setting failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Set successfully", c)
}

// GetUsersByAuthority
// @Tags      Authority
// @Summary   Get user ID list by authority
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     authorityId  query     uint                                                        true  "Role ID"
// @Success   200          {object}  response.Response{data=[]uint,msg=string}                   "Retrieved successfully"
// @Router    /authority/getUsersByAuthority [get]
func (a *AuthorityApi) GetUsersByAuthority(c *gin.Context) {
	var req systemReq.SetRoleUsers
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userIds, err := authorityService.GetUserIdsByAuthorityId(req.AuthorityId)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Retrieval failed: "+err.Error(), c)
		return
	}
	if userIds == nil {
		userIds = []uint{}
	}
	response.OkWithDetailed(userIds, "Retrieved successfully", c)
}

// SetRoleUsers
// @Tags      Authority
// @Summary   Fully replace the user list associated with a role
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetRoleUsers         true  "Role ID and user ID list"
// @Success   200   {object}  response.Response{msg=string}  "Set successfully"
// @Router    /authority/setRoleUsers [post]
func (a *AuthorityApi) SetRoleUsers(c *gin.Context) {
	var req systemReq.SetRoleUsers
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.AuthorityId == 0 {
		response.FailWithMessage("Role ID cannot be empty", c)
		return
	}
	if err := authorityService.SetRoleUsers(req.AuthorityId, req.UserIds); err != nil {
		global.GVA_LOG.Error("Failed to set!", zap.Error(err))
		response.FailWithMessage("Setting failed: "+err.Error(), c)
		return
	}
	response.OkWithMessage("Set successfully", c)
}
