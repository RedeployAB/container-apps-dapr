package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAuthenticate tests the authenticate middleware with table-driven tests.
func TestAuthenticate(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			keys map[string]struct{}
			req  func() *http.Request
		}
		wantCode int
	}{
		{
			name: "valid key",
			input: struct {
				keys map[string]struct{}
				req  func() *http.Request
			}{
				keys: map[string]struct{}{
					"valid-key": {},
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					req.Header.Set(authHeader, "valid-key")
					return req
				},
			},
			wantCode: http.StatusOK,
		},
		{
			name: "missing key",
			input: struct {
				keys map[string]struct{}
				req  func() *http.Request
			}{
				keys: map[string]struct{}{
					"valid-key": {},
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					return req
				},
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "invalid key",
			input: struct {
				keys map[string]struct{}
				req  func() *http.Request
			}{
				keys: map[string]struct{}{
					"valid-key": {},
				},
				req: func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					req.Header.Set(authHeader, "invalid-key")
					return req
				},
			},
			wantCode: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a mock response recorder.
			rr := httptest.NewRecorder()

			// Create a mock handler.
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			// Create a new request.
			req := test.input.req()

			// Call the authenticate middleware with the mock handler.
			authenticate(test.input.keys, handler).ServeHTTP(rr, req)

			// Check the status code.
			if status := rr.Code; status != test.wantCode {
				t.Errorf("handler returned wrong status code: got %v want %v\n",
					status, test.wantCode)
			}
		})
	}
}
