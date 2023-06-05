package report

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNewPubsubReporter(t *testing.T) {
	var tests = []struct {
		name  string
		input []PubsubReporterOptions
		want  *PubsubReporter
	}{
		{
			name:  "Empty",
			input: nil,
			want: &PubsubReporter{
				name:    defaultReporterName,
				topic:   defaultReporterTopic,
				timeout: defaultReporterTimeout,
			},
		},
		{
			name: "With options",
			input: []PubsubReporterOptions{
				{
					Name:    "name",
					Topic:   "topic",
					Timeout: time.Second * 5,
				},
			},
			want: &PubsubReporter{
				name:    "name",
				topic:   "topic",
				timeout: time.Second * 5,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newPubsubReporter(test.input...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(PubsubReporter{})); diff != "" {
				t.Errorf("newPubsubReporter() = unexpected, (-want +got):\n%s\n", diff)
			}
		})
	}
}

func TestPubsubReporter_Run(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			reporter *PubsubReporter
			report   Report
		}
		wantErr error
	}{
		{
			name: "Empty",
			input: struct {
				reporter *PubsubReporter
				report   Report
			}{
				reporter: &PubsubReporter{
					client: &mockClient{
						err: nil,
					},
					name:    defaultReporterName,
					topic:   defaultReporterTopic,
					timeout: defaultReporterTimeout,
				},
			},
		},
		{
			name: "With data",
			input: struct {
				reporter *PubsubReporter
				report   Report
			}{
				reporter: &PubsubReporter{
					client: &mockClient{
						err: nil,
					},
					name:    defaultReporterName,
					topic:   defaultReporterTopic,
					timeout: defaultReporterTimeout,
				},
				report: Report{
					ID:   "id",
					Data: []byte("data"),
				},
			},
		},
		{
			name: "With error",
			input: struct {
				reporter *PubsubReporter
				report   Report
			}{
				reporter: &PubsubReporter{
					client: &mockClient{
						err: errors.New("error"),
					},
					name:    defaultReporterName,
					topic:   defaultReporterTopic,
					timeout: defaultReporterTimeout,
				},
			},
			wantErr: errors.New("error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotErr := test.input.reporter.Run(test.input.report)

			if test.wantErr != nil && gotErr == nil {
				t.Errorf("PubsubReporter.Run() = unexpected, want error, got nil\n")
			}
		})
	}
}
