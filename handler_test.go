package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAlertEndpoints(t *testing.T) {
	srv := NewRouter(slog.New(slog.NewTextHandler(io.Discard, nil)))
	const ok = `{"status":"ok"}`

	tests := []struct {
		name, method, path, body string
		wantCode                 int
		wantBody                 string // checked only when non-empty
	}{
		{"plug", http.MethodPost, "/alerts/plug", `{"device":"cam-1"}`, http.StatusOK, ok},
		{"unplug", http.MethodPost, "/alerts/unplug", `{"device":"cam-1"}`, http.StatusOK, ok},
		{"plug non-json", http.MethodPost, "/alerts/plug", `hello`, http.StatusOK, ok},
		{"empty body", http.MethodPost, "/alerts/unplug", ``, http.StatusOK, ok},
		{"too large", http.MethodPost, "/alerts/plug", strings.Repeat("x", maxBodyBytes+1), http.StatusRequestEntityTooLarge, ""},
		{"wrong method", http.MethodGet, "/alerts/plug", ``, http.StatusMethodNotAllowed, ""},
		{"unknown", http.MethodPost, "/alerts/other", `{}`, http.StatusNotFound, ""},
		{"health", http.MethodGet, "/healthz", ``, http.StatusOK, "."},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			rec := httptest.NewRecorder()
			srv.ServeHTTP(rec, req)
			if rec.Code != tc.wantCode {
				t.Fatalf("code: got %d, want %d", rec.Code, tc.wantCode)
			}
			if tc.wantBody != "" && rec.Body.String() != tc.wantBody {
				t.Fatalf("body: got %q, want %q", rec.Body.String(), tc.wantBody)
			}
		})
	}
}
