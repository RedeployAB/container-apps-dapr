package server

import (
	"context"
	"errors"
	"testing"

	"github.com/dapr/go-sdk/service/common"
)

func TestPubsubReportHandler(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			reporter mockReporter
			e        *common.TopicEvent
		}
		wantRetry bool
		wantErr   error
	}{
		{
			name: "Success",
			input: struct {
				reporter mockReporter
				e        *common.TopicEvent
			}{
				reporter: mockReporter{},
				e: &common.TopicEvent{
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
				e        *common.TopicEvent
			}{
				reporter: mockReporter{},
				e: &common.TopicEvent{
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
				e        *common.TopicEvent
			}{
				reporter: mockReporter{},
				e: &common.TopicEvent{
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
				e        *common.TopicEvent
			}{
				reporter: mockReporter{
					err: errors.New("failed to create report"),
				},
				e: &common.TopicEvent{
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
			gotRetry, gotErr := s.pubsubReportHandler(context.Background(), test.input.e)

			if gotRetry != test.wantRetry {
				t.Errorf("pubsubReportHandler(%+v, %+v) = unexpected result, want retry %v, got %v\n", test.input.reporter, test.input.e, test.wantRetry, gotRetry)
			}

			if test.wantErr != nil && gotErr == nil {
				t.Errorf("pubsubReportHandler(%+v, %+v) = unexpected result, want error %v, got nil\n", test.input.reporter, test.input.e, test.wantErr)
			}
		})
	}
}
