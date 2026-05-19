package bootstrap

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSPreflightAllowsLoopbackOrigins(t *testing.T) {
	handler := NewMemoryApp(t.TempDir()).Handler()
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/workspaces", nil)
	req.Header.Set("Origin", "http://127.0.0.1:5173")
	req.Header.Set("Access-Control-Request-Headers", "content-type")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "http://127.0.0.1:5173" {
		t.Fatalf("allow origin = %q", got)
	}
	if got := rr.Header().Get("Access-Control-Allow-Headers"); got != "content-type" {
		t.Fatalf("allow headers = %q", got)
	}
}

func TestCORSHeadersOnAPIResponse(t *testing.T) {
	handler := NewMemoryApp(t.TempDir()).Handler()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("Origin", "http://localhost:5173")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("allow origin = %q", got)
	}
}
