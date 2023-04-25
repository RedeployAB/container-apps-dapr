package report

import (
	"errors"
	"testing"

	"github.com/RedeployAB/container-apps-dapr/common/report"
	"github.com/google/go-cmp/cmp"
)

func TestNewService(t *testing.T) {
	var tests = []struct {
		name    string
		input   Storer
		want    *service
		wantErr error
	}{
		{
			name:    "With nil storer",
			input:   nil,
			want:    nil,
			wantErr: errors.New("error"),
		},
		{
			name:    "With storer",
			input:   &mockStorer{},
			want:    &service{s: &mockStorer{}},
			wantErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewService(test.input)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(service{}, mockStorer{})); diff != "" {
				t.Errorf("NewService(%+v) = unexpected result (-want +got):\n%s\n", test.input, diff)
			}

			if test.wantErr != nil && gotErr == nil {
				t.Errorf("NewService(%+v) = unexpected result, want error %v, got nil\n", test.input, test.wantErr)
			}
		})

	}
}

func TestService_Create(t *testing.T) {
	var tests = []struct {
		name    string
		input   Storer
		wantErr error
	}{
		{
			name:    "With nil storer",
			input:   nil,
			wantErr: errors.New("error"),
		},
		{
			name:    "With storer",
			input:   &mockStorer{},
			wantErr: nil,
		},
		{
			name:    "Error creating report",
			input:   &mockStorer{err: errors.New("error")},
			wantErr: errors.New("error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := service{s: test.input}

			gotErr := service.Create(report.Report{})

			if test.wantErr != nil && gotErr == nil {
				t.Errorf("Create(%v) = unexpected result, want error %v, got nil\n", test.input, test.wantErr)
			}
		})

	}
}

type mockStorer struct {
	err error
}

func (s mockStorer) Store(r report.Report) error {
	return s.err
}
