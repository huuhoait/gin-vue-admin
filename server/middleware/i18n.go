package middleware

import (
	"github.com/huuhoait/gin-vue-admin/server/global/i18n"
	"github.com/gin-gonic/gin"
)

// CtxLocaleKey is the gin.Context key under which the resolved locale tag
// (e.g. "vi-VN") is stored by I18nLocale.
const CtxLocaleKey = "i18n_locale"

// I18nLocale resolves the caller's preferred locale from the Accept-Language
// header and stashes the canonical tag on the context. Handlers and
// response helpers read it back via LocaleFrom(c).
//
// Falls back to i18n.DefaultLocale when the header is missing or unrecognised.
func I18nLocale() gin.HandlerFunc {
	return func(c *gin.Context) {
		loc := i18n.Normalize(c.GetHeader("Accept-Language"))
		c.Set(CtxLocaleKey, loc)
		c.Next()
	}
}

// LocaleFrom returns the resolved locale for the current request, or
// i18n.DefaultLocale when the middleware has not run (e.g. in unit tests).
func LocaleFrom(c *gin.Context) string {
	if v, ok := c.Get(CtxLocaleKey); ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return i18n.DefaultLocale
}
