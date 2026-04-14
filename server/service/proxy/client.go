package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"log"
)

// Client is a thin HTTP client that forwards requests to a downstream
// SkyAgent service and parses the GVA envelope response.
type Client struct {
	baseURL string
	timeout time.Duration
	http    *http.Client
}

// NewClient creates a proxy Client for the given base URL.
// timeout ≤ 0 defaults to 10 s.
func NewClient(baseURL string, timeout time.Duration) *Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		timeout: timeout,
		http:    &http.Client{Timeout: timeout},
	}
}

// RequestOpts carries optional per-request overrides.
type RequestOpts struct {
	Headers map[string]string
	Query   url.Values
}

// Do executes an HTTP request against the downstream service, reads the full
// response body, and unmarshals the GVA envelope.
func (c *Client) Do(ctx context.Context, method, path string, body any, opts *RequestOpts) (*GVAEnvelope, int, error) {
	start := time.Now()
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	reqURL, err := url.Parse(c.baseURL + path)
	if err != nil {
		return nil, 0, fmt.Errorf("proxy: build url: %w", err)
	}
	if opts != nil && len(opts.Query) > 0 {
		reqURL.RawQuery = opts.Query.Encode()
	}

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return nil, 0, fmt.Errorf("proxy: marshal body: %w", err)
		}
		reader = bytes.NewReader(payload)
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(timeoutCtx, method, reqURL.String(), reader)
	if err != nil {
		return nil, 0, fmt.Errorf("proxy: create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Apply per-request headers (e.g. X-Maker-ID, X-Checker-ID).
	var traceID string
	if opts != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
			if strings.EqualFold(k, "X-Trace-Id") {
				traceID = v
			}
		}
	}

	// Log the outbound request (avoid logging bodies/secrets).
	if traceID != "" {
		log.Printf("proxy: outbound start trace_id=%s method=%s url=%s", traceID, method, reqURL.String())
	} else {
		log.Printf("proxy: outbound start method=%s url=%s", method, reqURL.String())
	}

	resp, err := c.http.Do(req)
	if err != nil {
		log.Printf("proxy: outbound error trace_id=%s method=%s url=%s dur_ms=%d err=%v",
			traceID, method, reqURL.String(), time.Since(start).Milliseconds(), err)
		return nil, 0, fmt.Errorf("proxy: downstream unreachable")
	}
	defer resp.Body.Close()

	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("proxy: outbound read_error trace_id=%s method=%s url=%s status=%d dur_ms=%d err=%v",
			traceID, method, reqURL.String(), resp.StatusCode, time.Since(start).Milliseconds(), err)
		return nil, resp.StatusCode, fmt.Errorf("proxy: read response: %w", err)
	}

	var envelope GVAEnvelope
	if len(rawBody) > 0 {
		if err := json.Unmarshal(rawBody, &envelope); err != nil {
			log.Printf("proxy: outbound parse_error trace_id=%s method=%s url=%s status=%d dur_ms=%d err=%v body_len=%d",
				traceID, method, reqURL.String(), resp.StatusCode, time.Since(start).Milliseconds(), err, len(rawBody))
			return nil, resp.StatusCode, fmt.Errorf("proxy: parse envelope: %w", err)
		}
	}

	log.Printf("proxy: outbound done trace_id=%s method=%s url=%s status=%d dur_ms=%d code=%d msg=%q",
		traceID, method, reqURL.String(), resp.StatusCode, time.Since(start).Milliseconds(), envelope.Code, envelope.Msg)

	return &envelope, resp.StatusCode, nil
}

// Get is a convenience wrapper for GET requests.
func (c *Client) Get(ctx context.Context, path string, query url.Values, headers map[string]string) (*GVAEnvelope, int, error) {
	return c.Do(ctx, http.MethodGet, path, nil, &RequestOpts{Headers: headers, Query: query})
}

// Post is a convenience wrapper for POST requests.
func (c *Client) Post(ctx context.Context, path string, body any, headers map[string]string) (*GVAEnvelope, int, error) {
	return c.Do(ctx, http.MethodPost, path, body, &RequestOpts{Headers: headers})
}

// Put is a convenience wrapper for PUT requests.
func (c *Client) Put(ctx context.Context, path string, body any, headers map[string]string) (*GVAEnvelope, int, error) {
	return c.Do(ctx, http.MethodPut, path, body, &RequestOpts{Headers: headers})
}

// Delete is a convenience wrapper for DELETE requests.
func (c *Client) Delete(ctx context.Context, path string, body any, headers map[string]string) (*GVAEnvelope, int, error) {
	return c.Do(ctx, http.MethodDelete, path, body, &RequestOpts{Headers: headers})
}
