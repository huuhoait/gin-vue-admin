package response

import (
	"net/http"

	"github.com/huuhoait/gin-vue-admin/server/global/i18n"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

const (
	ERROR   = 7
	SUCCESS = 0
)

// localeKey mirrors middleware.CtxLocaleKey.
// Duplicated as a string literal to avoid an import cycle
// (middleware imports this package).
const localeKey = "i18n_locale"

func localeFrom(c *gin.Context) string {
	if c == nil {
		return i18n.DefaultLocale
	}
	if v, ok := c.Get(localeKey); ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return i18n.DefaultLocale
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

// ---------------------------------------------------------------------------
// Legacy helpers (message passed in directly).
// Story 8.3 replaces the original Chinese defaults with i18n-resolved
// fallbacks sourced from the request's Accept-Language header. Callers that
// still pass a message literal are unchanged; only the zero-arg defaults are
// affected.
// ---------------------------------------------------------------------------

func Ok(c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, i18n.T(localeFrom(c), "admin.common.ok"), c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(SUCCESS, data, i18n.T(localeFrom(c), "admin.common.ok"), c)
}

func OkWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(SUCCESS, data, message, c)
}

func Fail(c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, i18n.T(localeFrom(c), "admin.common.fail"), c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(ERROR, map[string]interface{}{}, message, c)
}

func NoAuth(message string, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		7,
		nil,
		message,
	})
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(ERROR, data, message, c)
}

// ---------------------------------------------------------------------------
// Story 8.3 — i18n key-based helpers.
// Prefer these in new code. The `key` argument is an i18n bundle key such as
// "admin.common.create_failed"; the optional kv varargs are key/value pairs
// substituted into {{.name}} placeholders.
//
// Handlers MUST NOT pass raw localised strings to these helpers — pass the
// bundle key and keep translations in resource/i18n/*.toml.
// ---------------------------------------------------------------------------

// OkWithCode responds 200/SUCCESS with a translated message.
func OkWithCode(c *gin.Context, key string, kv ...string) {
	Result(SUCCESS, map[string]interface{}{}, i18n.T(localeFrom(c), key, kv...), c)
}

// OkWithDataCode responds 200/SUCCESS with data and a translated message.
func OkWithDataCode(c *gin.Context, data interface{}, key string, kv ...string) {
	Result(SUCCESS, data, i18n.T(localeFrom(c), key, kv...), c)
}

// FailWithCode responds 200/ERROR with a translated message.
// (HTTP 200 is retained for parity with existing gin-vue-admin clients that
// branch on the JSON `code` field rather than HTTP status.)
func FailWithCode(c *gin.Context, key string, kv ...string) {
	Result(ERROR, map[string]interface{}{}, i18n.T(localeFrom(c), key, kv...), c)
}

// FailWithDataCode responds 200/ERROR with data and a translated message.
func FailWithDataCode(c *gin.Context, data interface{}, key string, kv ...string) {
	Result(ERROR, data, i18n.T(localeFrom(c), key, kv...), c)
}
