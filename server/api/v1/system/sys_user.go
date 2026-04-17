package system

import (
	"strconv"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common"
	"github.com/huuhoait/gin-vue-admin/server/model/common/request"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	systemReq "github.com/huuhoait/gin-vue-admin/server/model/system/request"
	systemRes "github.com/huuhoait/gin-vue-admin/server/model/system/response"
	"github.com/huuhoait/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Login
// @Tags     Base
// @Summary  User login
// @Produce   application/json
// @Param    data  body      systemReq.Login                                             true  "Username, password, captcha"
// @Success  200   {object}  response.Response{data=systemRes.LoginResponse,msg=string}  "Returns user info, token, and expiration time"
// @Router   /base/login [post]
func (b *BaseApi) Login(c *gin.Context) {
	var l systemReq.Login
	err := c.ShouldBindJSON(&l)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(l, utils.LoginVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	key := c.ClientIP()
	// Check if captcha is enabled
	openCaptcha := global.GVA_CONFIG.Captcha.OpenCaptcha               // Whether to enable brute-force protection threshold
	openCaptchaTimeOut := global.GVA_CONFIG.Captcha.OpenCaptchaTimeOut // Cache timeout duration
	v, ok := global.BlackCache.Get(key)
	if !ok {
		global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
	}

	var oc bool = openCaptcha == 0 || openCaptcha < interfaceToInt(v)
	if oc && (l.Captcha == "" || l.CaptchaId == "" || !store.Verify(l.CaptchaId, l.Captcha, true)) {
		// Increment captcha attempt count
		global.BlackCache.Increment(key, 1)
		response.FailWithMessage("Captcha error", c)
		// Record login failure log
		loginLogService.CreateLoginLog(system.SysLoginLog{
			Username:     l.Username,
			Ip:           c.ClientIP(),
			Agent:        c.Request.UserAgent(),
			Status:       false,
			ErrorMessage: "Captcha error",
		})
		return
	}

	u := &system.SysUser{Username: l.Username, Password: l.Password}
	user, err := userService.Login(u)
	if err != nil {
		global.GVA_LOG.Error("Login failed! Username does not exist or password is incorrect!", zap.Error(err))
		// Increment captcha attempt count
		global.BlackCache.Increment(key, 1)
		response.FailWithMessage("Username does not exist or password is incorrect", c)
		// Record login failure log
		loginLogService.CreateLoginLog(system.SysLoginLog{
			Username:     l.Username,
			Ip:           c.ClientIP(),
			Agent:        c.Request.UserAgent(),
			Status:       false,
			ErrorMessage: "Username does not exist or password is incorrect",
		})
		return
	}
	if user.Enable != 1 {
		global.GVA_LOG.Error("Login failed! User is banned from logging in!")
		// Increment captcha attempt count
		global.BlackCache.Increment(key, 1)
		response.FailWithMessage("User is banned from logging in", c)
		// Record login failure log
		loginLogService.CreateLoginLog(system.SysLoginLog{
			Username:     l.Username,
			Ip:           c.ClientIP(),
			Agent:        c.Request.UserAgent(),
			Status:       false,
			ErrorMessage: "User is banned from logging in",
			UserID:       user.ID,
		})
		return
	}
	b.TokenNext(c, *user)
}

// TokenNext Issue JWT after login
func (b *BaseApi) TokenNext(c *gin.Context, user system.SysUser) {
	token, claims, err := utils.LoginToken(&user)
	if err != nil {
		global.GVA_LOG.Error("Failed to get token!", zap.Error(err))
		response.FailWithMessage("Failed to get token", c)
		return
	}
	// Record successful login log
	loginLogService.CreateLoginLog(system.SysLoginLog{
		Username: user.Username,
		Ip:       c.ClientIP(),
		Agent:    c.Request.UserAgent(),
		Status:   true,
		UserID:   user.ID,
		ErrorMessage: "Login successful",
	})
	if !global.GVA_CONFIG.System.UseMultipoint {
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "Login successful", c)
		return
	}

	if jwtStr, err := jwtService.GetRedisJWT(user.Username); err == redis.Nil {
		if err := utils.SetRedisJWT(token, user.Username); err != nil {
			global.GVA_LOG.Error("Failed to set login status!", zap.Error(err))
			response.FailWithMessage("Failed to set login status", c)
			return
		}
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "Login successful", c)
	} else if err != nil {
		global.GVA_LOG.Error("Failed to set login status!", zap.Error(err))
		response.FailWithMessage("Failed to set login status", c)
	} else {
		var blackJWT system.JwtBlacklist
		blackJWT.Jwt = jwtStr
		if err := jwtService.JsonInBlacklist(blackJWT); err != nil {
			response.FailWithMessage("Failed to invalidate JWT", c)
			return
		}
		if err := utils.SetRedisJWT(token, user.GetUsername()); err != nil {
			response.FailWithMessage("Failed to set login status", c)
			return
		}
		utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
		response.OkWithDetailed(systemRes.LoginResponse{
			User:      user,
			Token:     token,
			ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		}, "Login successful", c)
	}
}

// Register
// @Tags     SysUser
// @Summary  User registration
// @Produce   application/json
// @Param    data  body      systemReq.Register                                            true  "Username, nickname, password, authority ID"
// @Success  200   {object}  response.Response{data=systemRes.SysUserResponse,msg=string}  "User registration, returns user info"
// @Router   /user/admin_register [post]
func (b *BaseApi) Register(c *gin.Context) {
	var r systemReq.Register
	err := c.ShouldBindJSON(&r)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(r, utils.RegisterVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	var authorities []system.SysAuthority
	for _, v := range r.AuthorityIds {
		authorities = append(authorities, system.SysAuthority{
			AuthorityId: v,
		})
	}
	user := &system.SysUser{Username: r.Username, NickName: r.NickName, Password: r.Password, HeaderImg: r.HeaderImg, AuthorityId: r.AuthorityId, Authorities: authorities, Enable: r.Enable, Phone: r.Phone, Email: r.Email}
	userReturn, err := userService.Register(*user)
	if err != nil {
		global.GVA_LOG.Error("Failed to register!", zap.Error(err))
		response.FailWithDetailed(systemRes.SysUserResponse{User: userReturn}, "Registration failed", c)
		return
	}
	response.OkWithDetailed(systemRes.SysUserResponse{User: userReturn}, "Registration successful", c)
}

// ChangePassword
// @Tags      SysUser
// @Summary   Change user password
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      systemReq.ChangePasswordReq    true  "Username, old password, new password"
// @Success   200   {object}  response.Response{msg=string}  "Change user password"
// @Router    /user/changePassword [post]
func (b *BaseApi) ChangePassword(c *gin.Context) {
	var req systemReq.ChangePasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.ChangePasswordVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid := utils.GetUserID(c)
	u := &system.SysUser{GVA_MODEL: global.GVA_MODEL{ID: uid}, Password: req.Password}
	err = userService.ChangePassword(u, req.NewPassword)
	if err != nil {
		global.GVA_LOG.Error("Failed to change password!", zap.Error(err))
		response.FailWithMessage("Modification failed, original password does not match current account", c)
		return
	}
	response.OkWithMessage("Modified successfully", c)
}

// GetUserList
// @Tags      SysUser
// @Summary   Get user list with pagination
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.GetUserList                                        true  "Page number, page size"
// @Success   200   {object}  response.Response{data=response.PageResult,msg=string}  "Get paginated user list, returns list, total, page, page size"
// @Router    /user/getUserList [post]
func (b *BaseApi) GetUserList(c *gin.Context) {
	var pageInfo systemReq.GetUserList
	err := c.ShouldBindJSON(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := userService.GetUserInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "Retrieved successfully", c)
}

// SetUserAuthority
// @Tags      SysUser
// @Summary   Change user authority
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetUserAuth          true  "User UUID, authority ID"
// @Success   200   {object}  response.Response{msg=string}  "Set user authority"
// @Router    /user/setUserAuthority [post]
func (b *BaseApi) SetUserAuthority(c *gin.Context) {
	var sua systemReq.SetUserAuth
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if UserVerifyErr := utils.Verify(sua, utils.SetUserAuthorityVerify); UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}
	userID := utils.GetUserID(c)
	err = userService.SetUserAuthority(userID, sua.AuthorityId)
	if err != nil {
		global.GVA_LOG.Error("Failed to modify!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	claims := utils.GetUserInfo(c)
	claims.AuthorityId = sua.AuthorityId
	token, err := utils.NewJWT().CreateToken(*claims)
	if err != nil {
		global.GVA_LOG.Error("Failed to modify!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	c.Header("new-token", token)
	c.Header("new-expires-at", strconv.FormatInt(claims.ExpiresAt.Unix(), 10))
	utils.SetToken(c, token, int(claims.ExpiresAt.Unix()-time.Now().Unix()))
	response.OkWithMessage("Modified successfully", c)
}

// SetUserAuthorities
// @Tags      SysUser
// @Summary   Set user authorities
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      systemReq.SetUserAuthorities   true  "User UUID, authority IDs"
// @Success   200   {object}  response.Response{msg=string}  "Set user authorities"
// @Router    /user/setUserAuthorities [post]
func (b *BaseApi) SetUserAuthorities(c *gin.Context) {
	var sua systemReq.SetUserAuthorities
	err := c.ShouldBindJSON(&sua)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	authorityID := utils.GetUserAuthorityId(c)
	err = userService.SetUserAuthorities(authorityID, sua.ID, sua.AuthorityIds)
	if err != nil {
		global.GVA_LOG.Error("Failed to modify!", zap.Error(err))
		response.FailWithMessage("Modification failed", c)
		return
	}
	response.OkWithMessage("Modified successfully", c)
}

// DeleteUser
// @Tags      SysUser
// @Summary   Delete user
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      request.GetById                true  "User ID"
// @Success   200   {object}  response.Response{msg=string}  "Delete user"
// @Router    /user/deleteUser [delete]
func (b *BaseApi) DeleteUser(c *gin.Context) {
	var reqId request.GetById
	err := c.ShouldBindJSON(&reqId)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(reqId, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	jwtId := utils.GetUserID(c)
	if jwtId == uint(reqId.ID) {
		response.FailWithMessage("Deletion failed, cannot delete yourself.", c)
		return
	}
	err = userService.DeleteUser(reqId.ID)
	if err != nil {
		global.GVA_LOG.Error("Failed to delete!", zap.Error(err))
		response.FailWithMessage("Deletion failed", c)
		return
	}
	response.OkWithMessage("Deleted successfully", c)
}

// SetUserInfo
// @Tags      SysUser
// @Summary   Set user info
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysUser                                             true  "ID, username, nickname, avatar URL"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Set user info"
// @Router    /user/setUserInfo [put]
func (b *BaseApi) SetUserInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(user, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(user.AuthorityIds) != 0 {
		authorityID := utils.GetUserAuthorityId(c)
		err = userService.SetUserAuthorities(authorityID, user.ID, user.AuthorityIds)
		if err != nil {
			global.GVA_LOG.Error("Failed to set!", zap.Error(err))
			response.FailWithMessage("Failed to set", c)
			return
		}
	}
	err = userService.SetUserInfo(system.SysUser{
		GVA_MODEL: global.GVA_MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		Enable:    user.Enable,
	})
	if err != nil {
		global.GVA_LOG.Error("Failed to set!", zap.Error(err))
		response.FailWithMessage("Failed to set", c)
		return
	}
	response.OkWithMessage("Set successfully", c)
}

// SetSelfInfo
// @Tags      SysUser
// @Summary   Set self info
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      system.SysUser                                             true  "ID, username, nickname, avatar URL"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Set self info"
// @Router    /user/SetSelfInfo [put]
func (b *BaseApi) SetSelfInfo(c *gin.Context) {
	var user systemReq.ChangeUserInfo
	err := c.ShouldBindJSON(&user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user.ID = utils.GetUserID(c)
	err = userService.SetSelfInfo(system.SysUser{
		GVA_MODEL: global.GVA_MODEL{
			ID: user.ID,
		},
		NickName:  user.NickName,
		HeaderImg: user.HeaderImg,
		Phone:     user.Phone,
		Email:     user.Email,
		Enable:    user.Enable,
	})
	if err != nil {
		global.GVA_LOG.Error("Failed to set!", zap.Error(err))
		response.FailWithMessage("Failed to set", c)
		return
	}
	response.OkWithMessage("Set successfully", c)
}

// SetSelfSetting
// @Tags      SysUser
// @Summary   Set self settings
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Param     data  body      map[string]interface{}  true  "User settings data"
// @Success   200   {object}  response.Response{data=map[string]interface{},msg=string}  "Set user settings"
// @Router    /user/SetSelfSetting [put]
func (b *BaseApi) SetSelfSetting(c *gin.Context) {
	var req common.JSONMap
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = userService.SetSelfSetting(req, utils.GetUserID(c))
	if err != nil {
		global.GVA_LOG.Error("Failed to set!", zap.Error(err))
		response.FailWithMessage("Failed to set", c)
		return
	}
	response.OkWithMessage("Set successfully", c)
}

// GetUserInfo
// @Tags      SysUser
// @Summary   Get user info
// @Security  ApiKeyAuth
// @accept    application/json
// @Produce   application/json
// @Success   200  {object}  response.Response{data=map[string]interface{},msg=string}  "Get user info"
// @Router    /user/getUserInfo [get]
func (b *BaseApi) GetUserInfo(c *gin.Context) {
	uuid := utils.GetUserUuid(c)
	ReqUser, err := userService.GetUserInfo(uuid)
	if err != nil {
		global.GVA_LOG.Error("Failed to retrieve!", zap.Error(err))
		response.FailWithMessage("Failed to retrieve", c)
		return
	}
	response.OkWithDetailed(gin.H{"userInfo": ReqUser}, "Retrieved successfully", c)
}

// ResetPassword
// @Tags      SysUser
// @Summary   Reset user password
// @Security  ApiKeyAuth
// @Produce  application/json
// @Param     data  body      system.SysUser                 true  "ID"
// @Success   200   {object}  response.Response{msg=string}  "Reset user password"
// @Router    /user/resetPassword [post]
func (b *BaseApi) ResetPassword(c *gin.Context) {
	var rps systemReq.ResetPassword
	err := c.ShouldBindJSON(&rps)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.ResetPassword(rps.ID, rps.Password)
	if err != nil {
		global.GVA_LOG.Error("Failed to reset!", zap.Error(err))
		response.FailWithMessage("Reset failed"+err.Error(), c)
		return
	}
	response.OkWithMessage("Reset successful", c)
}
