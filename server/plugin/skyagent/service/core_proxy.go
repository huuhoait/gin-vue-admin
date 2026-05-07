package service

import "time"

// CoreProxy wraps a Client configured for the SkyAgent Core service.
type CoreProxy struct{ *Client }

// NewCoreProxy creates a CoreProxy from the provided URL and timeout.
func NewCoreProxy(baseURL string, timeout time.Duration) *CoreProxy {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8081"
	}
	return &CoreProxy{Client: NewClient(baseURL, timeout)}
}
