package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Get_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/v1/agents" {
			t.Errorf("expected /v1/agents, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("expected page=1, got %s", r.URL.Query().Get("page"))
		}
		if r.Header.Get("X-Maker-ID") != "42" {
			t.Errorf("expected X-Maker-ID=42, got %s", r.Header.Get("X-Maker-ID"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0,
			"data": map[string]any{"list": []any{}, "total": 0},
			"msg":  "Success",
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, 5*time.Second)
	query := make(map[string][]string)
	query["page"] = []string{"1"}

	env, status, err := client.Get(context.Background(), "/v1/agents", query, map[string]string{"X-Maker-ID": "42"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status != 200 {
		t.Errorf("expected status 200, got %d", status)
	}
	if !env.IsSuccess() {
		t.Errorf("expected success envelope, got code=%d", env.Code)
	}
}

func TestClient_Post_ErrorEnvelope(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"code": 11003,
			"data": map[string]any{
				"responseCode": 11003,
				"details":      []map[string]string{{"field": "full_name", "message": "max 255"}},
			},
			"msg": "Validation failed",
		})
	}))
	defer srv.Close()

	client := NewClient(srv.URL, 5*time.Second)
	env, _, err := client.Post(context.Background(), "/v1/agents", map[string]string{"full_name": "x"}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.IsSuccess() {
		t.Error("expected error envelope")
	}
	if env.Code != 11003 {
		t.Errorf("expected code 11003, got %d", env.Code)
	}
	ep := env.ParseErrorPayload()
	if ep == nil {
		t.Fatal("expected ErrorPayload")
	}
	if len(ep.Details) != 1 || ep.Details[0].Field != "full_name" {
		t.Errorf("unexpected details: %+v", ep.Details)
	}
}

func TestClient_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	client := NewClient(srv.URL, 50*time.Millisecond)
	_, _, err := client.Get(context.Background(), "/v1/agents", nil, nil)
	if err == nil {
		t.Error("expected timeout error")
	}
}

func TestClient_Unreachable(t *testing.T) {
	client := NewClient("http://127.0.0.1:1", 1*time.Second)
	_, _, err := client.Get(context.Background(), "/v1/agents", nil, nil)
	if err == nil {
		t.Error("expected connection error")
	}
}
