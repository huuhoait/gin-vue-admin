package request

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

func HttpRequest(
	urlStr string,
	method string,
	headers map[string]string,
	params map[string]string,
	data any) (*http.Response, error) {
	return doJSONRequest(context.Background(), 0, urlStr, method, headers, params, data)
}

// HttpRequestWithTimeout sendHTTPRequest, SupportCustomtimeout
// timeout ParameterCanChoose, UnitForTime.Duration, defaultValueFor 10 PartClock
func HttpRequestWithTimeout(
	urlStr string,
	method string,
	headers map[string]string,
	params map[string]string,
	data any,
	timeout ...time.Duration) (*http.Response, error) {
	t := 10 * time.Minute
	if len(timeout) > 0 && timeout[0] > 0 {
		t = timeout[0]
	}
	return doJSONRequest(context.Background(), t, urlStr, method, headers, params, data)
}

// HttpRequestWithContextAndTimeout sendHTTPRequest, SupportCustomtimeoutAndContext
func HttpRequestWithContextAndTimeout(
	ctx context.Context,
	urlStr string,
	method string,
	headers map[string]string,
	params map[string]string,
	data any,
	timeout ...time.Duration) (*http.Response, error) {
	t := 10 * time.Minute
	if len(timeout) > 0 {
		if timeout[0] < 0 {
			t = 0 // negativeValueTableShowNotsetTimeout(Used for SSE etc.FlowstyleScenario)
		} else if timeout[0] > 0 {
			t = timeout[0]
		}
	}
	return doJSONRequest(ctx, t, urlStr, method, headers, params, data)
}

func doJSONRequest(
	ctx context.Context,
	timeout time.Duration,
	urlStr string,
	method string,
	headers map[string]string,
	params map[string]string,
	data any) (*http.Response, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// URL
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	// AddQuery Parameter
	query := u.Query()
	for k, v := range params {
		query.Set(k, v)
	}
	u.RawQuery = query.Encode()

	// WillDataCodeForJSON
	buf := new(bytes.Buffer)
	if data != nil {
		b, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}

	// CreateRequest
	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if data != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	client := &http.Client{}
	if timeout > 0 {
		client.Timeout = timeout
	}

	// WhenRequest SSE FlowWhen, Disable Transport LayerofAutomatic gzip compress
	// Avoid gzip unzipRelieverushCause SSE EventUnablerealWhenToReach
	if req.Header.Get("Accept") == "text/event-stream" {
		client.Transport = &http.Transport{
			DisableCompression: true,
		}
	}

	// sendRequest
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// ReturnResponse, InvokeSideResponsibleDisable
	return resp, nil
}
