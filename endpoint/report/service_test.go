package report

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewService(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			r Reporter
		}
		want    *service
		wantErr error
	}{
		{
			name: "With reporter",
			input: struct {
				r Reporter
			}{
				r: mockReporter{},
			},
			want: &service{
				r: mockReporter{},
			},
		},
		{
			name: "With nil reporter",
			input: struct {
				r Reporter
			}{
				r: nil,
			},
			wantErr: errors.New("error creating service: reporter is nil"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := NewService(test.input.r)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(service{}, mockReporter{})); diff != "" {
				t.Errorf("NewService = unexpected, (-want +got):\n%s\n", diff)
			}

			if test.wantErr != nil && gotErr == nil {
				t.Errorf("NewService = want error, got nil\n")
			}
		})
	}
}

func TestService_Create(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			r    Reporter
			id   string
			data []byte
		}
		wantErr error
	}{
		{
			name: "Run reporter",
			input: struct {
				r    Reporter
				id   string
				data []byte
			}{
				r:    mockReporter{},
				id:   "id",
				data: []byte("data"),
			},
			wantErr: nil,
		},
		{
			name: "With nil reporter",
			input: struct {
				r    Reporter
				id   string
				data []byte
			}{
				r:    nil,
				id:   "id",
				data: []byte("data"),
			},
			wantErr: errors.New("error creating report: reporter is nil"),
		},
		{
			name: "With reporter error",
			input: struct {
				r    Reporter
				id   string
				data []byte
			}{
				r:    mockReporter{err: errors.New("error")},
				id:   "id",
				data: []byte("data"),
			},
			wantErr: errors.New("error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &service{r: test.input.r}

			gotErr := s.Create(NewReport(test.input.id, test.input.data))

			if test.wantErr != nil && gotErr == nil {
				t.Errorf("Create = want error, got nil\n")
			}
		})
	}

}

type mockReporter struct {
	err error
}

func (r mockReporter) Run(report Report) error {
	if r.err != nil {
		return r.err
	}
	return nil
}
