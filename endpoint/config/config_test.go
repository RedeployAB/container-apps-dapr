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
					Type:    defaultReporterType,
					Name:    defaultReporterName,
					Timeout: defaultReporterTimeout,
					Queue:   defaultReporterQueue,
					Topic:   defaultReporterTopic,
				},
			},
		},
		{
			name: "With environment variables",
			input: map[string]string{
				"ENDPOINT_HOST":             "localhost",
				"ENDPOINT_PORT":             "3001",
				"ENDPOINT_READ_TIMEOUT":     "10s",
				"ENDPOINT_WRITE_TIMEOUT":    "10s",
				"ENDPOINT_IDLE_TIMEOUT":     "10s",
				"ENDPOINT_REPORTER_TYPE":    "pubsub-test",
				"ENDPOINT_REPORTER_NAME":    "reports-test",
				"ENDPOINT_REPORTER_TIMEOUT": "5s",
				"ENDPOINT_REPORTER_QUEUE":   "create-test",
				"ENDPOINT_REPORTER_TOPIC":   "create-test",
				"ENDPOINT_SECURITY_KEYS":    "key1,key2",
			},
			want: &Configuration{
				Server: Server{
					Host:         "localhost",
					Port:         3001,
					ReadTimeout:  time.Second * 10,
					WriteTimeout: time.Second * 10,
					IdleTimeout:  time.Second * 10,
					Security: Security{
						Keys: map[string]struct{}{
							"key1": {},
							"key2": {},
						},
					},
				},
				Reporter: Reporter{
					Type:    "pubsub-test",
					Name:    "reports-test",
					Timeout: time.Second * 5,
					Queue:   "create-test",
					Topic:   "create-test",
				},
			},
		},
		{
			name: "With invalid environment variables",
			input: map[string]string{
				"ENDPOINT_PORT": "invalid",
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
