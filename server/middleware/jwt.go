package middleware

import (
	"errors"
	"strconv"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/utils"
	"github.com/golang-jwt/jwt/v5"

	"github.com/huuhoait/gin-vue-admin/server/model/common/response"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWT authentication retrieves the x-token header; the frontend should store the token in cookies or localStorage and coordinate expiration time with the backend
		token := utils.GetToken(c)
		if token == "" {
			response.NoAuth("Not logged in or unauthorized access, please log in", c)
			c.Abort()
			return
		}
		if isBlacklist(token) {
			response.NoAuth("Your account has been logged in from another location or the token is invalid", c)
			utils.ClearToken(c)
			c.Abort()
			return
		}
		j := utils.NewJWT()
		// parseToken parses the information contained in the token
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				response.NoAuth("Login has expired, please log in again", c)
				utils.ClearToken(c)
				c.Abort()
				return
			}
			response.NoAuth(err.Error(), c)
			utils.ClearToken(c)
			c.Abort()
			return
		}

		// logged-in user disabled by admin - need to invalidate their JWT; this is performance-intensive, enable if needed
		// user deletion logic needs optimization; this is performance-intensive, enable if needed

		//if user, err := userService.FindUserByUuid(claims.UUID.String()); err != nil || user.Enable == 2 {
		//	_ = jwtService.JsonInBlacklist(system.JwtBlacklist{Jwt: token})
		//	response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
		//	c.Abort()
		//}
		c.Set("claims", claims)
		if claims.ExpiresAt.Unix()-time.Now().Unix() < claims.BufferTime {
			// Cooldown: once a token has been refreshed in the last 5 minutes
			// skip re-issuance. Prevents a storm of Redis writes on every
			// request once we enter the BufferTime window (commonly 1 day).
			refreshKey := "jwt-refresh:" + token
			if _, hit := global.BlackCache.Get(refreshKey); !hit {
				global.BlackCache.Set(refreshKey, struct{}{}, 5*time.Minute)
				dr, _ := utils.ParseDuration(global.GVA_CONFIG.JWT.ExpiresTime)
				claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(dr))
				newToken, _ := j.CreateTokenByOldToken(token, *claims)
				newClaims, _ := j.ParseToken(newToken)
				c.Header("new-token", newToken)
				c.Header("new-expires-at", strconv.FormatInt(newClaims.ExpiresAt.Unix(), 10))
				utils.SetToken(c, newToken, int(dr.Seconds()/60))
				if global.GVA_CONFIG.System.UseMultipoint {
					// record new active JWT
					_ = utils.SetRedisJWT(newToken, newClaims.Username)
				}
			}
		}
		c.Next()

		if newToken, exists := c.Get("new-token"); exists {
			c.Header("new-token", newToken.(string))
		}
		if newExpiresAt, exists := c.Get("new-expires-at"); exists {
			c.Header("new-expires-at", newExpiresAt.(string))
		}
	}
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: IsBlacklist
//@description: check if JWT is in the blacklist
//@param: jwt string
//@return: bool

func isBlacklist(jwt string) bool {
	_, ok := global.BlackCache.Get(jwt)
	return ok
}
