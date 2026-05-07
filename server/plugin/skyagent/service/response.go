package service

import "encoding/json"

// GVAEnvelope represents the standard response envelope from Core/Order services.
// Success: code == 0; Error: code != 0.
type GVAEnvelope struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
	Msg  string          `json:"msg"`
}

// IsSuccess returns true when the upstream response indicates success.
func (e *GVAEnvelope) IsSuccess() bool { return e.Code == 0 }

// ErrorPayload represents the structured error detail returned in the data
// field when a validation or business error occurs.
type ErrorPayload struct {
	ResponseCode int           `json:"responseCode"`
	TraceID      string        `json:"traceId,omitempty"`
	Details      []FieldError  `json:"details,omitempty"`
}

// FieldError represents a single field-level validation error.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ParseErrorPayload attempts to decode the data field into an ErrorPayload.
// Returns nil if decoding fails (data may be null or a different shape).
func (e *GVAEnvelope) ParseErrorPayload() *ErrorPayload {
	if e.Data == nil || string(e.Data) == "null" {
		return nil
	}
	var ep ErrorPayload
	if err := json.Unmarshal(e.Data, &ep); err != nil {
		return nil
	}
	return &ep
}
