package system

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/huuhoait/gin-vue-admin/server/model/system/request"
	"github.com/huuhoait/gin-vue-admin/server/utils"
)

type OIDCApi struct{}

// OIDCLogin
// @Tags      Base
// @Summary   Redirect to OIDC provider login page
// @Produce   application/json
// @Router    /base/oidc/login [get]
func (o *OIDCApi) OIDCLogin(c *gin.Context) {
	// Generate a random state to prevent CSRF on the OAuth callback
	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		response.FailWithMessage("Failed to generate state", c)
		return
	}
	state := hex.EncodeToString(stateBytes)

	// Store state in a short-lived cookie for callback verification
	c.SetCookie("oidc_state", state, 300, "/", "", true, true)

	url, err := oidcService.LoginURL(c.Request.Context(), state)
	if err != nil {
		global.GVA_LOG.Error("OIDC login URL failed", zap.Error(err))
		response.FailWithMessage("OIDC not available: "+err.Error(), c)
		return
	}
	c.Redirect(http.StatusFound, url)
}

// OIDCCallback
// @Tags      Base
// @Summary   Handle OIDC provider callback, issue JWT
// @Produce   application/json
// @Router    /base/oidc/callback [get]
func (o *OIDCApi) OIDCCallback(c *gin.Context) {
	// Verify state to prevent CSRF
	cookieState, err := c.Cookie("oidc_state")
	if err != nil || cookieState != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid state parameter"})
		return
	}
	c.SetCookie("oidc_state", "", -1, "/", "", true, true)

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "missing code parameter"})
		return
	}

	user, err := oidcService.HandleCallback(c.Request.Context(), code)
	if err != nil {
		global.GVA_LOG.Error("OIDC callback failed", zap.Error(err))
		response.FailWithMessage("OIDC authentication failed: "+err.Error(), c)
		return
	}

	// Issue internal JWT (same as password login)
	j := utils.NewJWT()
	claims := j.CreateClaims(request.BaseClaims{
		UUID:        user.UUID,
		ID:          user.ID,
		NickName:    user.NickName,
		Username:    user.Username,
		AuthorityId: user.AuthorityId,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.GVA_LOG.Error("OIDC JWT creation failed", zap.Error(err))
		response.FailWithMessage("Token creation failed", c)
		return
	}

	utils.SetToken(c, token, int(claims.RegisteredClaims.ExpiresAt.Unix()-claims.RegisteredClaims.IssuedAt.Unix()))
	response.OkWithDetailed(gin.H{
		"token":     token,
		"expiresAt": claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
		"userInfo":  user,
	}, "Login successful", c)
}
