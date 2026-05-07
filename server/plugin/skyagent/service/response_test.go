package service

import (
	"encoding/json"
	"testing"
)

func TestGVAEnvelope_IsSuccess(t *testing.T) {
	tests := []struct {
		name string
		code int
		want bool
	}{
		{"success", 0, true},
		{"error", 11001, false},
		{"generic error", 7, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &GVAEnvelope{Code: tt.code}
			if got := e.IsSuccess(); got != tt.want {
				t.Errorf("IsSuccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGVAEnvelope_ParseErrorPayload_WithDetails(t *testing.T) {
	raw := json.RawMessage(`{"responseCode":11003,"traceId":"abc","details":[{"field":"full_name","message":"max 255"}]}`)
	e := &GVAEnvelope{Code: 11003, Data: raw, Msg: "Validation failed"}

	ep := e.ParseErrorPayload()
	if ep == nil {
		t.Fatal("expected non-nil ErrorPayload")
	}
	if ep.ResponseCode != 11003 {
		t.Errorf("ResponseCode = %d, want 11003", ep.ResponseCode)
	}
	if ep.TraceID != "abc" {
		t.Errorf("TraceID = %q, want %q", ep.TraceID, "abc")
	}
	if len(ep.Details) != 1 {
		t.Fatalf("Details len = %d, want 1", len(ep.Details))
	}
	if ep.Details[0].Field != "full_name" || ep.Details[0].Message != "max 255" {
		t.Errorf("Details[0] = %+v", ep.Details[0])
	}
}

func TestGVAEnvelope_ParseErrorPayload_Null(t *testing.T) {
	e := &GVAEnvelope{Code: 11001, Data: json.RawMessage(`null`), Msg: "Agent not found"}
	if ep := e.ParseErrorPayload(); ep != nil {
		t.Errorf("expected nil ErrorPayload for null data, got %+v", ep)
	}
}

func TestGVAEnvelope_ParseErrorPayload_Nil(t *testing.T) {
	e := &GVAEnvelope{Code: 7}
	if ep := e.ParseErrorPayload(); ep != nil {
		t.Errorf("expected nil ErrorPayload for nil data, got %+v", ep)
	}
}

func TestGVAEnvelope_ParseErrorPayload_NonObject(t *testing.T) {
	e := &GVAEnvelope{Code: 7, Data: json.RawMessage(`"just a string"`)}
	if ep := e.ParseErrorPayload(); ep != nil {
		t.Errorf("expected nil ErrorPayload for string data, got %+v", ep)
	}
}
