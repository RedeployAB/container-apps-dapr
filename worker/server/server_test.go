package server

import (
	"errors"
	"syscall"
	"testing"
	"time"

	"github.com/RedeployAB/container-apps-dapr/worker/report"
	"github.com/dapr/go-sdk/service/common"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestServerNew(t *testing.T) {
	tests := []struct {
		name    string
		input   Options
		want    *server
		wantErr error
	}{
		{
			name:    "With empty options",
			input:   Options{},
			want:    nil,
			wantErr: errors.New("error"),
		},
		{
			name: "With reporter in options",
			input: Options{
				Reporter: &mockReporter{},
			},
			want: &server{
				reporter: &mockReporter{},
				log:      logr.Logger{},
				address:  defaultAddress,
				name:     defaultPubsubName,
				topic:    defaultPubsubTopic,
			},
			wantErr: nil,
		},
		{
			name: "With options",
			input: Options{
				Reporter: &mockReporter{},
				Logger:   &mockLogger{},
				Address:  "localhost:3002",
				Name:     "reports-test",
				Topic:    "create-test",
			},
			want: &server{
				reporter: &mockReporter{},
				log:      &mockLogger{},
				address:  "localhost:3002",
				name:     "reports-test",
				topic:    "create-test",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := new(test.input)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(server{}, mockService{}, mockReporter{}), cmpopts.IgnoreUnexported(logr.Logger{})); diff != "" {
				t.Errorf("New() = unexpected result, (-want +got):\n%s\n", diff)
			}

			if test.wantErr != nil && gotErr == nil {
				t.Errorf("New() = unexpected result, want error %v, got nil\n", test.wantErr)
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
				service: &mockService{},
				log:     &mockLogger{},
			}
			go func() {
				time.Sleep(time.Millisecond * 100)
				syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			}()
			srv.Start()

			if diff := cmp.Diff(test.want, logMessages); diff != "" {
				t.Errorf("Start() = unexpected result, (-want +got):\n%s\n", diff)
			}
			logMessages = []string{}
		})
	}
}

type mockService struct {
	err      error
	startErr error
	stopErr  error
}

func (s mockService) Start() error {
	return s.startErr
}

func (s mockService) Stop() error {
	return s.stopErr
}

func (s mockService) AddTopicEventHandler(sub *common.Subscription, fn common.TopicEventHandler) error {
	return s.err
}

type mockReporter struct {
	err error
}

func (r mockReporter) Create(report.Report) error {
	return r.err
}

type mockLogger struct{}

var logMessages = []string{}

func (l mockLogger) Error(err error, msg string, keysAndValues ...any) {
	logMessages = append(logMessages, msg)
}

func (l mockLogger) Info(msg string, keysAndValues ...any) {
	logMessages = append(logMessages, msg)
}
