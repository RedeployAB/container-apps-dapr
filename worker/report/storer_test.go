package report

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/RedeployAB/container-apps-dapr/common/report"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/google/go-cmp/cmp"
)

func TestNewBlobStorer(t *testing.T) {
	var tests = []struct {
		name    string
		input   BlobStorerOptions
		want    *BlobStorer
		wantErr error
	}{
		{
			name:  "With empty options",
			input: BlobStorerOptions{},
			want: &BlobStorer{
				client:  nil,
				name:    defaultStorerName,
				timeout: defaultStorerTimeout,
			},
		},
		{
			name: "With options",
			input: BlobStorerOptions{
				Name:    "test",
				Timeout: time.Second * 30,
			},
			want: &BlobStorer{
				client:  nil,
				name:    "test",
				timeout: time.Second * 30,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := newBlobStorer(test.input)

			if diff := cmp.Diff(test.want, got, cmp.AllowUnexported(BlobStorer{})); diff != "" {
				t.Errorf("newBlobStorer(%+v) = unexpected result (-want +got):\n%s\n", test.input, diff)
			}
		})
	}
}

func TestBlogStorer_Store(t *testing.T) {
	var tests = []struct {
		name    string
		input   mockClient
		wantErr error
	}{
		{
			name:    "With successful store",
			input:   mockClient{},
			wantErr: nil,
		},
		{
			name:    "With unsuccessful store",
			input:   mockClient{err: errors.New("error")},
			wantErr: errors.New("error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storer := &BlobStorer{
				client:  &test.input,
				name:    "test",
				timeout: time.Second * 30,
			}

			gotErr := storer.Store(report.NewReport("123", []byte("test")))

			if test.wantErr != nil && gotErr == nil {
				t.Errorf("Store() = unexpected result, want error %v, got nil\n", test.wantErr)
			}
		})
	}
}

type mockClient struct {
	err error
}

func (c *mockClient) InvokeBinding(ctx context.Context, in *dapr.InvokeBindingRequest) (*dapr.BindingEvent, error) {
	return nil, c.err
}
