package server

import (
	"context"
	"errors"
	"testing"

	"github.com/dapr/go-sdk/service/common"
)

func TestQueueReportHandler(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			reporter mockReporter
			in       *common.BindingEvent
		}
		want    []byte
		wantErr error
	}{
		{
			name: "Success",
			input: struct {
				reporter mockReporter
				in       *common.BindingEvent
			}{
				reporter: mockReporter{},
				in: &common.BindingEvent{
					Data:     []byte(`{"id":"123","data":"testdata"}`),
					Metadata: map[string]string{},
				},
			},
			want:    []byte(`Message processed.`),
			wantErr: nil,
		},
		{
			name: "Failed to deserialize report",
			input: struct {
				reporter mockReporter
				in       *common.BindingEvent
			}{
				reporter: mockReporter{},
				in: &common.BindingEvent{
					Data:     []byte(`{"id":"123","data":"testdata`),
					Metadata: map[string]string{},
				},
			},
		},
		{
			name: "Failed to create report",
			input: struct {
				reporter mockReporter
				in       *common.BindingEvent
			}{
				reporter: mockReporter{
					err: errors.New("failed to create report"),
				},
				in: &common.BindingEvent{
					Data:     []byte(`{"id":"123","data":"testdata"}`),
					Metadata: map[string]string{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &server{
				log:      mockLogger{},
				reporter: &test.input.reporter,
			}
			got, gotErr := s.queueReportHandler(context.Background(), test.input.in)

			if string(got) != string(test.want) {
				t.Errorf("queueReportHandler(%+v, %+v) = unexpected result, want retry %v, got %v\n", test.input.reporter, test.input.in, test.want, got)
			}

			if test.wantErr != nil && gotErr == nil {
				t.Errorf("queueReportHandler(%+v, %+v) = unexpected result, want error %v, got nil\n", test.input.reporter, test.input.in, test.wantErr)
			}
		})
	}
}
