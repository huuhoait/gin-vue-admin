package middleware

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/huuhoait/gin-vue-admin/server/config"
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/gin-gonic/gin"
)

var corsWarnOnce sync.Once

// Cors allows all cross-origin requests — development only.
//
// Reflecting the request Origin together with Access-Control-Allow-Credentials
// means any page the user visits can fire authenticated requests at the API,
// which is equivalent to disabling the Same-Origin Policy for credentialed
// traffic. Release mode therefore refuses to return credentials in this mode;
// operators must configure cors.mode = "strict-whitelist" and list the real
// frontend origins.
func Cors() gin.HandlerFunc {
	release := gin.Mode() == gin.ReleaseMode
	if release {
		corsWarnOnce.Do(func() {
			fmt.Println("[WARN] CORS allow-all mode is enabled in release; switch cors.mode to 'strict-whitelist' in production")
		})
	}
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			// Only send credentials when we are explicitly *not* in release:
			// reflecting Origin + credentials = trust any site that loads in
			// the victim's browser.
			if !release {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token,X-Token,X-User-Id")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type, New-Token, New-Expires-At")

		// allow all OPTIONS methods
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// process request
		c.Next()
	}
}

// CorsByRules handles cross-origin requests according to configuration
func CorsByRules() gin.HandlerFunc {
	// allow all
	if global.GVA_CONFIG.Cors.Mode == "allow-all" {
		return Cors()
	}
	return func(c *gin.Context) {
		whitelist := checkCors(c.GetHeader("origin"))

		// passed check, add request headers
		if whitelist != nil {
			c.Header("Access-Control-Allow-Origin", whitelist.AllowOrigin)
			c.Header("Access-Control-Allow-Headers", whitelist.AllowHeaders)
			c.Header("Access-Control-Allow-Methods", whitelist.AllowMethods)
			c.Header("Access-Control-Expose-Headers", whitelist.ExposeHeaders)
			if whitelist.AllowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		// strict whitelist mode and check failed, reject request directly
		if whitelist == nil && global.GVA_CONFIG.Cors.Mode == "strict-whitelist" && !(c.Request.Method == "GET" && c.Request.URL.Path == "/health") {
			c.AbortWithStatus(http.StatusForbidden)
		} else {
			// non-strict whitelist mode, allow all OPTIONS methods regardless of check result
			if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusNoContent)
			}
		}

		// process request
		c.Next()
	}
}

func checkCors(currentOrigin string) *config.CORSWhitelist {
	for _, whitelist := range global.GVA_CONFIG.Cors.Whitelist {
		// iterate through configured CORS headers to find a match
		if currentOrigin == whitelist.AllowOrigin {
			return &whitelist
		}
	}
	return nil
}
