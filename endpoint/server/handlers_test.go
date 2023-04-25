package server

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestReportHandler tests the reportHandler method with table driven tests,
// with the help of the httptest package and a mock reporter and logger.
func TestReportHandler(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			method string
			body   string
			err    error
		}
		wantCode int
		wantBody string
	}{
		{
			name: "With valid request",
			input: struct {
				method string
				body   string
				err    error
			}{
				method: http.MethodPost,
				body:   `{"id":"123","data":"data"}`,
				err:    nil,
			},
			wantCode: http.StatusOK,
			wantBody: `{"id":"123","data":"data"}`,
		},
		{
			name: "With invalid method",
			input: struct {
				method string
				body   string
				err    error
			}{
				method: http.MethodGet,
				body:   `{"id":"123","data":"data"}`,
				err:    nil,
			},
			wantCode: http.StatusMethodNotAllowed,
			wantBody: "Method not allowed\n",
		},
		{
			name: "With invalid body",
			input: struct {
				method string
				body   string
				err    error
			}{
				method: http.MethodPost,
				body:   `{"id":"123","data":"data`,
				err:    nil,
			},
			wantCode: http.StatusBadRequest,
			wantBody: "Invalid request body\n",
		},
		{
			name: "With reporter error",
			input: struct {
				method string
				body   string
				err    error
			}{
				method: http.MethodPost,
				body:   `{"id":"123","data":"data"}`,
				err:    errors.New("error"),
			},
			wantCode: http.StatusInternalServerError,
			wantBody: "Internal server error\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &server{
				reporter: &mockReporter{
					err: test.input.err,
				},
				log: &mockLogger{},
			}

			req := httptest.NewRequest(test.input.method, "/reports", strings.NewReader(test.input.body))
			w := httptest.NewRecorder()

			s.reportHandler().ServeHTTP(w, req)

			resp := w.Result()

			if resp.StatusCode != test.wantCode {
				t.Errorf("reportHandler() = unexpected result, want %d, got: %d\n", resp.StatusCode, test.wantCode)
			}

			body, _ := io.ReadAll(resp.Body)
			if string(body) != test.wantBody {
				t.Errorf("reportHandler() = unexpected result, want %s, got: %s\n", test.wantBody, string(body))
			}
		})
	}
}
