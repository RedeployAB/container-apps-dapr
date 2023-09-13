package report

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestNewQueueReporter(t *testing.T) {
	var tests = []struct {
		name  string
		input []QueueReporterOption
		want  *QueueReporter
	}{
		{
			name:  "Empty",
			input: nil,
			want: &QueueReporter{
				name:    defaultReporterName,
				queue:   defaultReporterQueue,
				timeout: defaultReporterTimeout,
			},
		},
		{
			name: "With options",
			input: []QueueReporterOption{
				func(o *QueueReporterOptions) {
					o.Name = "name"
					o.Queue = "queue"
					o.Timeout = time.Second * 5
				},
			},
			want: &QueueReporter{
				name:    "name",
				queue:   "queue",
				timeout: time.Second * 5,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newQueueReporter(test.input...)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(QueueReporter{})); diff != "" {
				t.Errorf("newQueueReporter() = unexpected, (-want +got):\n%s\n", diff)
			}
		})
	}
}

func TestQueueReporter_Run(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			reporter *QueueReporter
			report   Report
		}
		wantErr error
	}{
		{
			name: "Empty",
			input: struct {
				reporter *QueueReporter
				report   Report
			}{
				reporter: &QueueReporter{
					client: &mockClient{
						err: nil,
					},
					name:    defaultReporterName,
					queue:   defaultReporterQueue,
					timeout: defaultReporterTimeout,
				},
			},
		},
		{
			name: "With data",
			input: struct {
				reporter *QueueReporter
				report   Report
			}{
				reporter: &QueueReporter{
					client: &mockClient{
						err: nil,
					},
					name:    defaultReporterName,
					queue:   defaultReporterQueue,
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
				reporter *QueueReporter
				report   Report
			}{
				reporter: &QueueReporter{
					client: &mockClient{
						err: errors.New("error"),
					},
					name:    defaultReporterName,
					queue:   defaultReporterQueue,
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
				t.Errorf("QueueReporter.Run() = unexpected, want error, got nil\n")
			}
		})
	}
}
