package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/middleware"
	"go.uber.org/zap"
)

// SkyAgent BFF responseCode catalog (subset).
// Source of truth: external-frontend-integration.md §5.
const (
	codeBadRequest          = 4001
	codeInternalError       = 5001
	codeUpstreamUnavailable = 5002
	codeUpstreamTimeout     = 5003
	codeUpstreamFailure     = 5004
)

// Respond writes the GVA envelope back to the client.
// On success (code == 0) it passes through data + msg.
// When the upstream envelope is missing it surfaces a contract-compliant
// 502/5004 instead of a 200/7 sentinel — see external-frontend-integration.md §3.
func Respond(c *gin.Context, envelope *GVAEnvelope, httpStatus int) {
	if envelope == nil {
		c.JSON(http.StatusBadGateway, errorEnvelope(
			codeUpstreamFailure,
			"Upstream service unavailable",
			traceIDFrom(c),
		))
		return
	}
	c.JSON(httpStatus, envelope)
}

// RespondError writes a proxy-level error (timeout, unreachable, parse failure).
// The error is classified into a SkyAgent responseCode + matching HTTP status
// so the FE's translateError(`errors.<code>`) pipeline lights up.
func RespondError(c *gin.Context, err error) {
	code, status, msg := classifyUpstreamError(err)
	c.JSON(status, errorEnvelope(code, msg, traceIDFrom(c)))
}

// RespondBadRequest writes a 400/4001 response for malformed input
// (typically a c.ShouldBindJSON failure). The raw err is logged server-side
// (correlatable via traceId) but NOT echoed to the client — gin's binding
// error strings can leak struct internals.
func RespondBadRequest(c *gin.Context, err error) {
	if err != nil && global.GVA_LOG != nil {
		global.GVA_LOG.Warn("skyagent bad request",
			zap.String("trace_id", traceIDFrom(c)),
			zap.String("path", c.FullPath()),
			zap.Error(err),
		)
	}
	c.JSON(http.StatusBadRequest, errorEnvelope(
		codeBadRequest,
		"Invalid request body",
		traceIDFrom(c),
	))
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

func errorEnvelope(code int, msg, traceID string) gin.H {
	data := gin.H{"responseCode": code}
	if traceID != "" {
		data["traceId"] = traceID
	}
	return gin.H{
		"code": code,
		"data": data,
		"msg":  msg,
	}
}

func traceIDFrom(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if id := c.GetHeader("X-Trace-Id"); id != "" {
		return id
	}
	return middleware.GetRequestID(c)
}

func classifyUpstreamError(err error) (code, status int, msg string) {
	switch {
	case errors.Is(err, ErrUpstreamTimeout), errors.Is(err, context.DeadlineExceeded):
		return codeUpstreamTimeout, http.StatusGatewayTimeout, "Request timed out"
	case errors.Is(err, ErrUpstreamUnreachable):
		return codeUpstreamUnavailable, http.StatusServiceUnavailable, "Upstream service unavailable"
	case errors.Is(err, ErrUpstreamReadFailed), errors.Is(err, ErrUpstreamParseFailed):
		return codeUpstreamFailure, http.StatusBadGateway, "Upstream dependency failure"
	default:
		return codeInternalError, http.StatusInternalServerError, "Internal error"
	}
}
