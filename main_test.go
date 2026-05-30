package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestExternalServiceCall(t *testing.T) {
	if time.Now().UnixNano()%2 == 0 {
		t.Fatal("context deadline exceeded: upstream service did not respond")
	}
}

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	healthz(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if rr.Body.String() != "ok" {
		t.Errorf("expected body 'ok', got %q", rr.Body.String())
	}
}

func TestHelloDefault(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rr := httptest.NewRecorder()
	hello(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["message"] != "hello, world" {
		t.Errorf("expected 'hello, world', got %q", resp["message"])
	}
}

func TestHelloWithName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/hello?name=alejandro", nil)
	rr := httptest.NewRecorder()
	hello(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp["message"] != "hello, alejandro" {
		t.Errorf("expected 'hello, alejandro', got %q", resp["message"])
	}
}
