package api

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/plugin/oauth2server/model/request"
	"github.com/huuhoait/gin-vue-admin/server/plugin/oauth2server/service"
	"github.com/huuhoait/gin-vue-admin/server/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type oauth2Api struct{}

// Authorize handles GET /oauth2/authorize. Requires the admin's JWT (mounted
// under the private group). Auto-approves and redirects to redirect_uri with
// the code. Errors are reported via redirect when redirect_uri is valid, per
// RFC 6749 §4.1.2.1; otherwise returned as JSON to avoid open-redirect.
//
// @Tags     OAuth2
// @Summary  authorization-code endpoint
// @Security ApiKeyAuth
// @Produce  text/html
// @Param    response_type query string true "must be 'code'"
// @Param    client_id query string true "registered client_id"
// @Param    redirect_uri query string true "must match a registered URI exactly"
// @Param    scope query string false "space-separated scopes"
// @Param    state query string false "opaque value echoed back"
// @Success  302
// @Router   /oauth2/authorize [get]
func (a *oauth2Api) Authorize(c *gin.Context) {
	var req request.AuthorizeReq
	if err := c.ShouldBindQuery(&req); err != nil {
		writeOAuthJSONError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	client, err := serviceClient.FindByClientID(req.ClientID)
	if err != nil {
		writeOAuthJSONError(c, http.StatusBadRequest, "invalid_client", "unknown client_id")
		return
	}
	if !client.Enabled {
		writeOAuthJSONError(c, http.StatusBadRequest, "unauthorized_client", "client disabled")
		return
	}
	// Validate redirect_uri BEFORE we trust it for error redirects.
	if !serviceClient.ValidateRedirect(client, req.RedirectURI) {
		writeOAuthJSONError(c, http.StatusBadRequest, "invalid_request", "redirect_uri not whitelisted")
		return
	}
	userID := utils.GetUserID(c)
	if userID == 0 {
		redirectErr(c, req.RedirectURI, req.State, "access_denied", "user not authenticated")
		return
	}
	code, err := serviceOAuth2.IssueCode(client, userID, req.RedirectURI, req.Scope)
	if err != nil {
		redirectErr(c, req.RedirectURI, req.State, oauthErrorCode(err), err.Error())
		return
	}
	target := req.RedirectURI
	q := url.Values{"code": {code}}
	if req.State != "" {
		q.Set("state", req.State)
	}
	if strings.Contains(target, "?") {
		target += "&" + q.Encode()
	} else {
		target += "?" + q.Encode()
	}
	c.Redirect(http.StatusFound, target)
}

// Token handles POST /oauth2/token. Public endpoint (no JWT). Client
// authentication via Basic auth or body (client_id/client_secret).
//
// @Tags     OAuth2
// @Summary  token endpoint
// @accept   application/x-www-form-urlencoded
// @Produce  application/json
// @Router   /oauth2/token [post]
func (a *oauth2Api) Token(c *gin.Context) {
	var req request.TokenReq
	if err := c.ShouldBind(&req); err != nil {
		writeOAuthJSONError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	clientID, secret := extractClientCreds(c, req.ClientID, req.ClientSecret)
	if clientID == "" || secret == "" {
		writeOAuthJSONError(c, http.StatusUnauthorized, "invalid_client", "missing client credentials")
		return
	}
	client, err := serviceClient.AuthenticateClient(clientID, secret)
	if err != nil {
		global.GVA_LOG.Warn("oauth2 client auth failed", zap.String("client_id", clientID))
		writeOAuthJSONError(c, http.StatusUnauthorized, "invalid_client", "")
		return
	}
	switch req.GrantType {
	case "authorization_code":
		resp, err := serviceOAuth2.ExchangeCode(client, req.Code, req.RedirectURI)
		if err != nil {
			writeOAuthJSONError(c, http.StatusBadRequest, oauthErrorCode(err), "")
			return
		}
		c.JSON(http.StatusOK, resp)
	case "refresh_token":
		resp, err := serviceOAuth2.RefreshAccessToken(client, req.RefreshToken)
		if err != nil {
			writeOAuthJSONError(c, http.StatusBadRequest, oauthErrorCode(err), "")
			return
		}
		c.JSON(http.StatusOK, resp)
	case "client_credentials":
		resp, err := serviceOAuth2.ClientCredentials(client, req.Scope)
		if err != nil {
			writeOAuthJSONError(c, http.StatusBadRequest, oauthErrorCode(err), "")
			return
		}
		c.JSON(http.StatusOK, resp)
	default:
		writeOAuthJSONError(c, http.StatusBadRequest, "unsupported_grant_type", "")
	}
}

// Introspect handles POST /oauth2/introspect (RFC 7662). Requires client
// authentication. Returns {active:false} for any unknown/expired/revoked
// token to avoid leaking which tokens existed.
//
// @Tags     OAuth2
// @Summary  token introspection
// @Router   /oauth2/introspect [post]
func (a *oauth2Api) Introspect(c *gin.Context) {
	var req request.IntrospectReq
	if err := c.ShouldBind(&req); err != nil {
		writeOAuthJSONError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	clientID, secret := extractClientCreds(c, "", "")
	if clientID == "" || secret == "" {
		writeOAuthJSONError(c, http.StatusUnauthorized, "invalid_client", "")
		return
	}
	if _, err := serviceClient.AuthenticateClient(clientID, secret); err != nil {
		writeOAuthJSONError(c, http.StatusUnauthorized, "invalid_client", "")
		return
	}
	c.JSON(http.StatusOK, serviceOAuth2.Introspect(c.Request.Context(), req.Token))
}

// Revoke handles POST /oauth2/revoke (RFC 7009).
// @Tags     OAuth2
// @Summary  token revocation
// @Router   /oauth2/revoke [post]
func (a *oauth2Api) Revoke(c *gin.Context) {
	var req request.RevokeReq
	if err := c.ShouldBind(&req); err != nil {
		writeOAuthJSONError(c, http.StatusBadRequest, "invalid_request", err.Error())
		return
	}
	clientID, secret := extractClientCreds(c, "", "")
	if clientID == "" || secret == "" {
		writeOAuthJSONError(c, http.StatusUnauthorized, "invalid_client", "")
		return
	}
	if _, err := serviceClient.AuthenticateClient(clientID, secret); err != nil {
		writeOAuthJSONError(c, http.StatusUnauthorized, "invalid_client", "")
		return
	}
	_ = serviceOAuth2.Revoke(req.Token)
	c.JSON(http.StatusOK, gin.H{"revoked": true})
}

// extractClientCreds returns clientID and secret from either Authorization:
// Basic header or body fallback. RFC 6749 §2.3.1 requires Basic to be
// supported and lets servers choose whether to also accept body params.
func extractClientCreds(c *gin.Context, bodyID, bodySecret string) (string, string) {
	if id, sec, ok := c.Request.BasicAuth(); ok {
		return id, sec
	}
	return bodyID, bodySecret
}

func writeOAuthJSONError(c *gin.Context, status int, code, desc string) {
	body := gin.H{"error": code}
	if desc != "" {
		body["error_description"] = desc
	}
	c.JSON(status, body)
}

func redirectErr(c *gin.Context, redirectURI, state, code, desc string) {
	q := url.Values{"error": {code}}
	if desc != "" {
		q.Set("error_description", desc)
	}
	if state != "" {
		q.Set("state", state)
	}
	target := redirectURI
	if strings.Contains(target, "?") {
		target += "&" + q.Encode()
	} else {
		target += "?" + q.Encode()
	}
	c.Redirect(http.StatusFound, target)
}

func oauthErrorCode(err error) string {
	switch {
	case errors.Is(err, service.ErrInvalidGrant):
		return "invalid_grant"
	case errors.Is(err, service.ErrUnsupportedGrant):
		return "unsupported_grant_type"
	case errors.Is(err, service.ErrInvalidScope):
		return "invalid_scope"
	case errors.Is(err, service.ErrInvalidRedirectURI):
		return "invalid_request"
	case errors.Is(err, service.ErrClientDisabled):
		return "unauthorized_client"
	default:
		return "server_error"
	}
}
