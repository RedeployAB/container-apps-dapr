package report

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReport(t *testing.T) {
	got := Report{
		ID:   "id",
		Data: []byte("data"),
	}

	want := Report{
		ID:   "id",
		Data: []byte("data"),
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Report = unexpected, (-want +got):\n%s\n", diff)
	}
}

func TestNewReport(t *testing.T) {
	var tests = []struct {
		name  string
		input struct {
			id   string
			data []byte
		}
		want Report
	}{
		{
			name: "Empty",
			input: struct {
				id   string
				data []byte
			}{
				id:   "",
				data: nil,
			},
			want: Report{},
		},
		{
			name: "With data",
			input: struct {
				id   string
				data []byte
			}{
				id:   "id",
				data: []byte("data"),
			},
			want: Report{
				ID:   "id",
				Data: []byte("data"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := NewReport(test.input.id, test.input.data)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("NewReport() = unexpected, (-want +got):\n%s\n", diff)
			}
		})
	}
}

func TestReport_JSON(t *testing.T) {
	var tests = []struct {
		name  string
		input Report
		want  []byte
	}{
		{
			name:  "Empty",
			input: Report{},
			want:  []byte(`{"ID":"","Data":null}`),
		},
		{
			name: "With data",
			input: Report{
				ID:   "id",
				Data: []byte("data"),
			},
			want: []byte(`{"ID":"id","Data":"ZGF0YQ=="}`),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.input.JSON()

			if diff := cmp.Diff(string(test.want), string(got)); diff != "" {
				t.Errorf("JSON() = unexpected, (-want +got):\n%s\n", diff)
			}
		})
	}
}
