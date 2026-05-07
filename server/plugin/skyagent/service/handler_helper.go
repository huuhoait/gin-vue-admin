package service

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Respond writes the GVA envelope back to the client.
// On success (code == 0) it passes through data + msg.
// On error it maps the upstream envelope to the gin-vue-admin error format,
// including field-level validation details when present.
func Respond(c *gin.Context, envelope *GVAEnvelope, httpStatus int) {
	if envelope == nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"code": 7,
			"data": nil,
			"msg":  "Upstream service unavailable",
		})
		return
	}

	// Pass-through: the upstream envelope is already in GVA format.
	c.JSON(httpStatus, envelope)
}

// RespondError writes a proxy-level error (timeout, unreachable, parse failure).
func RespondError(c *gin.Context, err error) {
	c.JSON(http.StatusOK, gin.H{
		"code": 7,
		"data": nil,
		"msg":  "Service request failed",
	})
}

// RespondRaw writes an arbitrary data payload wrapped in a success envelope.
func RespondRaw(c *gin.Context, data any) {
	raw, _ := json.Marshal(data)
	c.JSON(http.StatusOK, &GVAEnvelope{
		Code: 0,
		Data: raw,
		Msg:  "Success",
	})
}
