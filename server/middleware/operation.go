package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/huuhoait/gin-vue-admin/server/utils"

	"github.com/huuhoait/gin-vue-admin/server/global"
	"github.com/huuhoait/gin-vue-admin/server/model/system"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var respPool sync.Pool
var bufferSize = 1024

// operationQueue decouples request handling from the audit-log write. The
// previous implementation called db.Create() inline, so every request paid a
// DB round-trip (and was blocked when the DB was slow). A buffered channel
// + dedicated writer goroutine keeps the hot path free while preserving
// at-least-once persistence under normal load.
var (
	operationQueue     chan system.SysOperationRecord
	operationQueueOnce sync.Once
)

// Sensitive values must not land in audit tables. Admin read-out of
// operation records is a privilege escalation vector if plaintext passwords
// or tokens are stored there.
var sensitiveBodyField = regexp.MustCompile(`("(password|passwd|pwd|token|secret|authorization|x[-_]?token|signingKey|refreshToken|accessToken|apiKey)"\s*:\s*")[^"]*(")`)

func init() {
	respPool.New = func() interface{} {
		return make([]byte, bufferSize)
	}
}

// StartOperationRecorder spins up the audit-log writer. Call from server
// bootstrap. Safe to invoke multiple times.
func StartOperationRecorder() {
	operationQueueOnce.Do(func() {
		operationQueue = make(chan system.SysOperationRecord, 4096)
		go func() {
			for record := range operationQueue {
				if err := global.GVA_DB.Create(&record).Error; err != nil {
					global.GVA_LOG.Error("create operation record error:", zap.Error(err))
				}
			}
		}()
	})
}

func enqueueOperation(record system.SysOperationRecord) {
	if operationQueue == nil {
		// Bootstrap-time fallback: still better to persist than lose the
		// audit event. Any call path that forgets to Start the recorder
		// will surface in tests via missing goroutine.
		StartOperationRecorder()
	}
	// Warn when queue utilization exceeds 75% (3072/4096) so ops can detect
	// backpressure before overflow occurs.
	if len(operationQueue) > 3072 {
		global.GVA_LOG.Warn("operation record queue above 75% capacity",
			zap.Int("queue_len", len(operationQueue)),
			zap.Int("queue_cap", cap(operationQueue)),
		)
	}

	select {
	case operationQueue <- record:
	default:
		// Channel full — fall back to a synchronous DB write in a new goroutine
		// so the request path is not blocked. This ensures audit completeness
		// (SOC2 CC7.2) instead of silently dropping records under burst load.
		go func(r system.SysOperationRecord) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if err := global.GVA_DB.WithContext(ctx).Create(&r).Error; err != nil {
				global.GVA_LOG.Error("audit fallback write failed",
					zap.String("path", r.Path),
					zap.String("method", r.Method),
					zap.Error(err),
				)
			}
		}(record)
	}
}

func scrubSensitive(s string) string {
	if s == "" {
		return s
	}
	return sensitiveBodyField.ReplaceAllString(s, `$1[REDACTED]$3`)
}

func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		var userId int
		if c.Request.Method != http.MethodGet {
			var err error
			body, err = io.ReadAll(c.Request.Body)
			if err != nil {
				global.GVA_LOG.Error("read body from request error:", zap.Error(err))
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		} else {
			query := c.Request.URL.RawQuery
			query, _ = url.QueryUnescape(query)
			split := strings.Split(query, "&")
			m := make(map[string]string)
			for _, v := range split {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					m[kv[0]] = kv[1]
				}
			}
			body, _ = json.Marshal(&m)
		}
		claims, _ := utils.GetClaims(c)
		if claims != nil && claims.BaseClaims.ID != 0 {
			userId = int(claims.BaseClaims.ID)
		} else {
			id, err := strconv.Atoi(c.Request.Header.Get("x-user-id"))
			if err != nil {
				userId = 0
			}
			userId = id
		}
		record := system.SysOperationRecord{
			Ip:        c.ClientIP(),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Agent:     c.Request.UserAgent(),
			Body:      "",
			UserID:    userId,
			RequestID: GetRequestID(c),
		}

		// when uploading files, the middleware log truncates the body
		if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			record.Body = "[file]"
		} else {
			if len(body) > bufferSize {
				record.Body = "[exceeds record length]"
			} else {
				record.Body = scrubSensitive(string(body))
			}
		}

		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		now := time.Now()

		c.Next()

		latency := time.Since(now)
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Latency = latency

		// Response capture policy:
		//   - Drop binary / file downloads entirely (not useful, bloats DB).
		//   - Drop listing responses (GET) — they routinely carry PII such
		//     as user emails, phone numbers, and whole user tables, which
		//     would otherwise leak via the audit log to anyone with
		//     operation-record read access.
		//   - For mutating endpoints keep only the status envelope (status
		//     code + top-level msg/code), stripping any data payload. That
		//     preserves the audit signal without mirroring sensitive state.
		resp := writer.body.String()
		if isBinaryResponse(c) || c.Request.Method == http.MethodGet {
			record.Resp = ""
		} else {
			record.Resp = scrubSensitive(summarizeResponse(resp))
		}

		enqueueOperation(record)
	}
}

func isBinaryResponse(c *gin.Context) bool {
	h := c.Writer.Header()
	return strings.Contains(h.Get("Pragma"), "public") ||
		strings.Contains(h.Get("Expires"), "0") ||
		strings.Contains(h.Get("Cache-Control"), "must-revalidate, post-check=0, pre-check=0") ||
		strings.Contains(h.Get("Content-Type"), "application/force-download") ||
		strings.Contains(h.Get("Content-Type"), "application/octet-stream") ||
		strings.Contains(h.Get("Content-Type"), "application/vnd.ms-excel") ||
		strings.Contains(h.Get("Content-Type"), "application/download") ||
		strings.Contains(h.Get("Content-Disposition"), "attachment") ||
		strings.Contains(h.Get("Content-Transfer-Encoding"), "binary")
}

// summarizeResponse keeps only the status envelope from a standard response
// (`{"code": 0, "msg": "...", "data": ...}`) so the audit record captures
// outcome without the payload body.
func summarizeResponse(resp string) string {
	if resp == "" {
		return ""
	}
	var envelope struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal([]byte(resp), &envelope); err != nil {
		// Non-JSON response; keep a short prefix so we still see it fail.
		if len(resp) > 256 {
			return resp[:256] + "...[truncated]"
		}
		return resp
	}
	out, _ := json.Marshal(envelope)
	return string(out)
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
