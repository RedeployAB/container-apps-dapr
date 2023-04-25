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
					Host:  defaultHost,
					Port:  defaultPort,
					Name:  defaultName,
					Topic: defaultTopic,
				},
				Storer: Storer{
					Type:    defaultStorerType,
					Name:    defaultStorerName,
					Timeout: defaultStorerTimeout,
				},
			},
		},
		{
			name: "With environment variables",
			input: map[string]string{
				"WORKER_HOST":           "localhost",
				"WORKER_PORT":           "3001",
				"WORKER_NAME":           "reports-test",
				"WORKER_TOPIC":          "create-test",
				"WORKER_STORER_TYPE":    "blob-test",
				"WORKER_STORER_NAME":    "reports-test",
				"WORKER_STORER_TIMEOUT": "5s",
			},
			want: &Configuration{
				Server: Server{
					Host:  "localhost",
					Port:  3001,
					Name:  "reports-test",
					Topic: "create-test",
				},
				Storer: Storer{
					Type:    "blob-test",
					Name:    "reports-test",
					Timeout: time.Second * 5,
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
