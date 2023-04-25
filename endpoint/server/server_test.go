package server

import (
	"errors"
	"net/http"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/RedeployAB/container-apps-dapr/common/report"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNew(t *testing.T) {
	var tests = []struct {
		name    string
		input   Options
		want    *server
		wantErr error
	}{
		{
			name:    "With defaults (no reporter)",
			input:   Options{},
			want:    nil,
			wantErr: errors.New("reporter is required"),
		},
		{
			name: "With defaults",
			input: Options{
				Reporter: &mockReporter{},
			},
			want: &server{
				httpServer: &http.Server{
					Addr:         ":" + strconv.Itoa(defaultPort),
					Handler:      &mockRouter{},
					ReadTimeout:  defaultReadTimeout,
					WriteTimeout: defaultWriteTimeout,
					IdleTimeout:  defaultIdleTimeout,
				},
				router:   &mockRouter{},
				log:      logr.Logger{},
				reporter: &mockReporter{},
			},
		},
		{
			name: "With options",
			input: Options{
				Host:         "localhost",
				Port:         3001,
				ReadTimeout:  time.Second * 10,
				WriteTimeout: time.Second * 10,
				IdleTimeout:  time.Second * 20,
				Logger:       mockLogger{},
				Reporter:     &mockReporter{},
			},
			want: &server{
				httpServer: &http.Server{
					Addr:         "localhost:3001",
					Handler:      &mockRouter{},
					ReadTimeout:  time.Second * 10,
					WriteTimeout: time.Second * 10,
					IdleTimeout:  time.Second * 20,
				},
				router:   &mockRouter{},
				log:      mockLogger{},
				reporter: &mockReporter{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := New(&mockRouter{}, test.input)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(server{}, mockReporter{}), cmpopts.IgnoreUnexported(http.Server{}, logr.Logger{})); diff != "" {
				t.Errorf("New(%+v) = unexpected result, (-want, +got)\n%s\n", test.input, diff)
			}

			if test.wantErr != nil && gotErr == nil {
				t.Errorf("New(%+v) = unexpected result, want error, got nil\n", test.input)
			}
		})
	}
}

func TestServer_Start(t *testing.T) {
	var tests = []struct {
		name string
		want []string
	}{
		{
			name: "Start",
			want: []string{
				"Server started.",
				"Server stopped.",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logMessages = []string{}
			srv := &server{
				httpServer: &http.Server{
					Addr: "localhost:3000",
				},
				router:   &mockRouter{},
				log:      mockLogger{},
				reporter: &mockReporter{},
			}
			go func() {
				time.Sleep(time.Millisecond * 100)
				syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			}()
			srv.Start()

			if diff := cmp.Diff(test.want, logMessages); diff != "" {
				t.Errorf("Start() = unexpected result, (-want, +got)\n%s\n", diff)
			}
			logMessages = []string{}
		})
	}
}

type mockRouter struct{}

func (ro mockRouter) Handle(pattern string, handler http.Handler) {}

func (ro mockRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {}

type mockLogger struct{}

var logMessages = []string{}

func (l mockLogger) Error(err error, msg string, keysAndValues ...any) {
	logMessages = append(logMessages, msg)
}

func (l mockLogger) Info(msg string, keysAndValues ...any) {
	logMessages = append(logMessages, msg)
}

type mockReporter struct {
	err error
}

func (r mockReporter) Create(report.Report) error {
	if r.err != nil {
		return r.err
	}
	return nil
}
