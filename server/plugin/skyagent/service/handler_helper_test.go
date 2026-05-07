package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() { gin.SetMode(gin.TestMode) }

// newTestCtx returns a gin.Context wrapping a fresh ResponseRecorder
// plus a stub request, so handlers under test can write a real response.
func newTestCtx(t *testing.T) (*gin.Context, *httptest.ResponseRecorder) {
	t.Helper()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "/admin-api/v1/agents", nil)
	c.Request = req
	return c, w
}

// decodeBody decodes the JSON body and returns code, responseCode, traceId.
// Anything missing comes back as zero/empty.
func decodeBody(t *testing.T, body []byte) (code int, respCode int, traceID string) {
	t.Helper()
	var got struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			ResponseCode int    `json:"responseCode"`
			TraceID      string `json:"traceId"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("decode body: %v\nbody=%s", err, body)
	}
	return got.Code, got.Data.ResponseCode, got.Data.TraceID
}

func TestRespondBadRequest_Shape(t *testing.T) {
	c, w := newTestCtx(t)
	c.Request.Header.Set("X-Trace-Id", "trace-abc")

	RespondBadRequest(c, errors.New("json: cannot unmarshal"))

	if w.Code != http.StatusBadRequest {
		t.Errorf("HTTP status = %d, want 400", w.Code)
	}
	code, respCode, trace := decodeBody(t, w.Body.Bytes())
	if code != codeBadRequest || respCode != codeBadRequest {
		t.Errorf("code=%d responseCode=%d, want both %d", code, respCode, codeBadRequest)
	}
	if trace != "trace-abc" {
		t.Errorf("traceId = %q, want %q", trace, "trace-abc")
	}
	if got := w.Body.String(); contains(got, "cannot unmarshal") {
		t.Errorf("response leaked binding error detail: %s", got)
	}
}

func TestRespondError_ClassifiesUpstreamErrors(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		wantCode int
		wantHTTP int
	}{
		{"timeout sentinel", fmt.Errorf("proxy: %w", ErrUpstreamTimeout), codeUpstreamTimeout, http.StatusGatewayTimeout},
		{"raw deadline", context.DeadlineExceeded, codeUpstreamTimeout, http.StatusGatewayTimeout},
		{"unreachable", fmt.Errorf("proxy: %w", ErrUpstreamUnreachable), codeUpstreamUnavailable, http.StatusServiceUnavailable},
		{"read failed", fmt.Errorf("proxy: %w", ErrUpstreamReadFailed), codeUpstreamFailure, http.StatusBadGateway},
		{"parse failed", fmt.Errorf("proxy: %w", ErrUpstreamParseFailed), codeUpstreamFailure, http.StatusBadGateway},
		{"unknown", errors.New("something else"), codeInternalError, http.StatusInternalServerError},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, w := newTestCtx(t)
			RespondError(c, tc.err)
			if w.Code != tc.wantHTTP {
				t.Errorf("HTTP status = %d, want %d", w.Code, tc.wantHTTP)
			}
			code, respCode, _ := decodeBody(t, w.Body.Bytes())
			if code != tc.wantCode || respCode != tc.wantCode {
				t.Errorf("code=%d responseCode=%d, want both %d", code, respCode, tc.wantCode)
			}
		})
	}
}

func TestRespond_NilEnvelope_Returns502_5004(t *testing.T) {
	c, w := newTestCtx(t)
	Respond(c, nil, http.StatusOK) // httpStatus is ignored when envelope is nil

	if w.Code != http.StatusBadGateway {
		t.Errorf("HTTP status = %d, want 502", w.Code)
	}
	code, respCode, _ := decodeBody(t, w.Body.Bytes())
	if code != codeUpstreamFailure || respCode != codeUpstreamFailure {
		t.Errorf("code=%d responseCode=%d, want both %d", code, respCode, codeUpstreamFailure)
	}
}

func TestRespond_PassesThroughEnvelope(t *testing.T) {
	c, w := newTestCtx(t)
	env := &GVAEnvelope{Code: 11003, Data: json.RawMessage(`{"responseCode":11003}`), Msg: "Validation failed"}
	Respond(c, env, http.StatusUnprocessableEntity)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("HTTP status = %d, want 422", w.Code)
	}
	var got map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if int(got["code"].(float64)) != 11003 {
		t.Errorf("code = %v, want 11003", got["code"])
	}
}

func TestTraceIDFrom_PrefersHeaderOverRequestID(t *testing.T) {
	c, _ := newTestCtx(t)
	c.Request.Header.Set("X-Trace-Id", "from-header")
	if got := traceIDFrom(c); got != "from-header" {
		t.Errorf("traceIDFrom = %q, want %q", got, "from-header")
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
