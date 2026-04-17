package middleware

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
)

const (
	csrfCookieName = "x-csrf-token"
	csrfHeaderName = "X-CSRF-Token"
	csrfTokenBytes = 32
)

// CSRF implements the double-submit-cookie pattern.
//
// When the frontend authenticates via the `x-token` HTTP header (the default
// for localStorage-stored tokens) the browser will not attach that header to
// cross-origin requests automatically, so CSRF is not reachable. We only
// enforce the check when the request carries the `x-token` cookie *without*
// the matching header — the cookie-auth flow is the CSRF-exposed one.
//
// For mutating methods we require a cookie named `x-csrf-token` whose value
// matches the `X-CSRF-Token` request header. If no cookie exists yet we mint
// one so the frontend can read it from JS and echo it back on the next write.
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		ensureCSRFCookie(c)

		if isSafeMethod(c.Request.Method) {
			c.Next()
			return
		}

		// Header-based auth flow: the attacker cannot read or forge the
		// x-token header from a cross-origin page, so the request is already
		// CSRF-safe and the CSRF pair is not required.
		if c.GetHeader("x-token") != "" {
			c.Next()
			return
		}

		// No session cookie means there is no authenticated state to protect.
		if cookie, err := c.Cookie("x-token"); err != nil || cookie == "" {
			c.Next()
			return
		}

		cookieTok, err := c.Cookie(csrfCookieName)
		headerTok := c.GetHeader(csrfHeaderName)
		if err != nil || cookieTok == "" || headerTok == "" ||
			subtle.ConstantTimeCompare([]byte(cookieTok), []byte(headerTok)) != 1 {
			response.FailWithDetailed(gin.H{"reload": false}, "csrf token mismatch", c)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	return false
}

// ensureCSRFCookie sets a readable (HttpOnly=false) CSRF cookie on the first
// request so SPAs can echo it back as an X-CSRF-Token header.
func ensureCSRFCookie(c *gin.Context) {
	if existing, err := c.Cookie(csrfCookieName); err == nil && existing != "" {
		return
	}
	buf := make([]byte, csrfTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return
	}
	tok := hex.EncodeToString(buf)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(csrfCookieName, tok, 0, "/", "", false, false)
}
