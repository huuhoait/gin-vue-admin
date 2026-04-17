package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestIDGeneratesWhenMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	var got string
	r.GET("/", func(c *gin.Context) {
		got = GetRequestID(c)
		c.Status(200)
	})

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got == "" {
		t.Fatal("expected a request id to be generated")
	}
	if w.Header().Get(RequestIDHeader) != got {
		t.Fatalf("header mismatch: %q vs ctx %q", w.Header().Get(RequestIDHeader), got)
	}
}

func TestRequestIDPassesThroughUpstream(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(RequestID())
	var got string
	r.GET("/", func(c *gin.Context) {
		got = GetRequestID(c)
		c.Status(200)
	})

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set(RequestIDHeader, "ext-12345")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if got != "ext-12345" {
		t.Fatalf("expected upstream id preserved, got %q", got)
	}
}
