package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDHeader is the HTTP header we propagate for correlation. Chosen to
// match the de-facto convention used by most load balancers and observability
// stacks (AWS ALB, Envoy, Istio, Datadog APM).
const RequestIDHeader = "X-Request-ID"

// RequestIDKey is the gin.Context key where the id is stored for downstream
// handlers and log fields.
const RequestIDKey = "request_id"

// RequestID attaches a correlation id to every request. If the caller
// already supplied one (trusted upstream proxies, sidecars, etc.) we keep
// it; otherwise we mint a UUIDv4 so the rest of the stack (zap fields, audit
// rows, outbound HTTP) can stitch logs together.
//
// Must run before any logger/operation/CSRF middleware so later stages can
// read the id from the context.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}
		c.Set(RequestIDKey, id)
		c.Writer.Header().Set(RequestIDHeader, id)
		c.Next()
	}
}

// GetRequestID returns the correlation id attached by RequestID(), or empty
// if the middleware wasn't run (tests / raw contexts).
func GetRequestID(c *gin.Context) string {
	if c == nil {
		return ""
	}
	if v, ok := c.Get(RequestIDKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return c.GetHeader(RequestIDHeader)
}
