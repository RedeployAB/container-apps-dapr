package server

import (
	"context"
	"errors"
	"testing"

	"github.com/dapr/go-sdk/service/common"
)

func TestReportHandler(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			reporter mockReporter
			req      *common.TopicEvent
		}
		wantRetry bool
		wantErr   error
	}{
		{
			name: "Success",
			input: struct {
				reporter mockReporter
				req      *common.TopicEvent
			}{
				reporter: mockReporter{},
				req: &common.TopicEvent{
					ID:         "test",
					PubsubName: "test",
					Topic:      "test",
					Data:       `{"id":"123","data":"testdata"}`,
				},
			},
			wantRetry: false,
			wantErr:   nil,
		},
		{
			name: "Failed to convert data to string",
			input: struct {
				reporter mockReporter
				req      *common.TopicEvent
			}{
				reporter: mockReporter{},
				req: &common.TopicEvent{
					ID:         "test",
					PubsubName: "test",
					Topic:      "test",
					Data:       123,
				},
			},
		},
		{
			name: "Failed to deserialize report",
			input: struct {
				reporter mockReporter
				req      *common.TopicEvent
			}{
				reporter: mockReporter{},
				req: &common.TopicEvent{
					ID:         "test",
					PubsubName: "test",
					Topic:      "test",
					Data:       `{"id":"123","data":"testdata`,
				},
			},
		},
		{
			name: "Failed to create report",
			input: struct {
				reporter mockReporter
				req      *common.TopicEvent
			}{
				reporter: mockReporter{
					err: errors.New("failed to create report"),
				},
				req: &common.TopicEvent{
					ID:         "test",
					PubsubName: "test",
					Topic:      "test",
					Data:       `{"id":"123","data":"testdata"}`,
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
			gotRetry, gotErr := s.reportHandler(context.Background(), test.input.req)

			if gotRetry != test.wantRetry {
				t.Errorf("reportHandler(%+v, %+v) = unexpected result, want retry %v, got %v\n", test.input.reporter, test.input.req, test.wantRetry, gotRetry)
			}

			if test.wantErr != nil && gotErr == nil {
				t.Errorf("reportHandler(%+v, %+v) = unexpected result, want error %v, got nil\n", test.input.reporter, test.input.req, test.wantErr)
			}
		})
	}
}
