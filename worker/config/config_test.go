package config

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNew(t *testing.T) {
	var tests = []struct {
		name    string
		input   map[string]string
		want    *Configuration
		wantErr error
	}{
		{
			name:  "With default values",
			input: map[string]string{},
			want: &Configuration{
				Server: Server{
					Host:         defaultHost,
					Port:         defaultPort,
					ReadTimeout:  defaultReadTimeout,
					WriteTimeout: defaultWriteTimeout,
					IdleTimeout:  defaultIdleTimeout,
				},
				Reporter: Reporter{
					Type:          defaultReporterType,
					PubsubName:    defaultReporterPubsubName,
					PubsubTopic:   defaultReporterPubsubTopic,
					PubsubTimeout: defaultReporterPubsubTimeout,
					Storer: Storer{
						Type: defaultReporterStorerType,
					},
				},
			},
		},
		{
			name: "With environment variables",
			input: map[string]string{
				"WORKER_HOST":                    "localhost",
				"WORKER_PORT":                    "3001",
				"WORKER_READ_TIMEOUT":            "10s",
				"WORKER_WRITE_TIMEOUT":           "10s",
				"WORKER_IDLE_TIMEOUT":            "10s",
				"WORKER_REPORTER_TYPE":           "storage",
				"WORKER_REPORTER_PUBSUB_NAME":    "reports-test",
				"WORKER_REPORTER_PUBSUB_TOPIC":   "create-test",
				"WORKER_REPORTER_PUBSUB_TIMEOUT": "5s",
				"WORKER_REPORTER_STORER_TYPE":    "blob-test",
			},
			want: &Configuration{
				Server: Server{
					Host:         "localhost",
					Port:         3001,
					ReadTimeout:  time.Second * 10,
					WriteTimeout: time.Second * 10,
					IdleTimeout:  time.Second * 10,
				},
				Reporter: Reporter{
					Type:          "storage",
					PubsubName:    "reports-test",
					PubsubTopic:   "create-test",
					PubsubTimeout: time.Second * 5,
					Storer: Storer{
						Type: "blob-test",
					},
				},
			},
		},
		{
			name: "With invalid environment variables",
			input: map[string]string{
				"WORKER_PORT": "invalid",
			},
			want:    nil,
			wantErr: errors.New("error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setEnvVars(test.input)
			defer unsetEnvVars(test.input)

			got, gotErr := New()

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("New() = unexpected, (-want, +got):\n%s\n", diff)
			}

			if test.wantErr != nil && gotErr == nil {
				t.Errorf("New() = unexpected, want: %v, got nil\n", test.wantErr)
			}
		})
	}

}

func setEnvVars(vars map[string]string) {
	os.Clearenv()
	for k, v := range vars {
		os.Setenv(k, v)
	}
}

func unsetEnvVars(vars map[string]string) {
	for _, v := range vars {
		os.Unsetenv(v)
	}
}
