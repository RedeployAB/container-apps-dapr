package server

import (
	"net/http"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestNew(t *testing.T) {
	var tests = []struct {
		name  string
		input Options
		want  *server
	}{
		{
			name:  "defaults",
			input: Options{},
			want: &server{
				httpServer: &http.Server{
					Addr:         ":" + strconv.Itoa(defaultPort),
					Handler:      &mockRouter{},
					ReadTimeout:  defaultReadTimeout,
					WriteTimeout: defaultWriteTimeout,
					IdleTimeout:  defaultIdleTimeout,
				},
				router: &mockRouter{},
				log:    logr.Logger{},
			},
		},
		{
			name: "with options",
			input: Options{
				Host:         "localhost",
				Port:         3001,
				ReadTimeout:  time.Second * 10,
				WriteTimeout: time.Second * 10,
				IdleTimeout:  time.Second * 20,
				Logger:       mockLogger{},
			},
			want: &server{
				httpServer: &http.Server{
					Addr:         "localhost:3001",
					Handler:      &mockRouter{},
					ReadTimeout:  time.Second * 10,
					WriteTimeout: time.Second * 10,
					IdleTimeout:  time.Second * 20,
				},
				router: &mockRouter{},
				log:    mockLogger{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := New(&mockRouter{}, test.input)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(server{}, mockLogger{}), cmpopts.IgnoreUnexported(http.Server{}, logr.Logger{})); diff != "" {
				t.Errorf("New(%+v) = unexpected result, (-want, +got)\n%s\n", test.input, diff)
			}
		})
	}
}

func TestServerStart(t *testing.T) {
	var tests = []struct {
		name string
		want []string
	}{
		{
			name: "start",
			want: []string{
				"Server started",
				"Server stopped",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			srv := &server{
				httpServer: &http.Server{
					Addr: "localhost:3000",
				},
				log: mockLogger{},
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

type mockLogger struct {
	logs []string
}

var logMessages = []string{}

func (l mockLogger) Error(err error, msg string, keysAndValues ...any) {
	logMessages = append(logMessages, msg)
}

func (l mockLogger) Info(msg string, keysAndValues ...any) {
	logMessages = append(logMessages, msg)
}
