package proxy

import "time"

// OrderProxy wraps a Client configured for the SkyAgent Order service.
type OrderProxy struct{ *Client }

// NewOrderProxy creates an OrderProxy from the provided URL and timeout.
func NewOrderProxy(baseURL string, timeout time.Duration) *OrderProxy {
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8081"
	}
	return &OrderProxy{Client: NewClient(baseURL, timeout)}
}
