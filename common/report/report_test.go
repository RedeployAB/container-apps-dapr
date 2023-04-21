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

// TestSerializeDeserialize tests the serialization and deserialization of a Report.
func TestSerializeDeserialize(t *testing.T) {
	var tests = []struct {
		name  string
		input Report
		want  Report
	}{
		{
			name:  "Empty",
			input: Report{},
			want:  Report{},
		},
		{
			name: "With data",
			input: Report{
				ID:   "id",
				Data: []byte("data"),
			},
			want: Report{
				ID:   "id",
				Data: []byte("data"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Report{}
			got.Deserialize(test.input.Serialize())

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Serialize/Deserialize() = unexpected, (-want +got):\n%s\n", diff)
			}
		})
	}
}
