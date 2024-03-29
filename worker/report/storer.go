package report

import (
	"context"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

const (
	defaultStorerName    = "reports-output"
	defaultStorerTimeout = time.Second * 10
)

// client is the interface that wraps around method InvokeBinding.
type client interface {
	InvokeBinding(ctx context.Context, in *dapr.InvokeBindingRequest) (*dapr.BindingEvent, error)
}

// BlobStorer is a storer that stores reports in a blob storage.
type BlobStorer struct {
	client
	name    string
	timeout time.Duration
}

// BlobStorerOptions contains options for BlobStorer.
type BlobStorerOptions struct {
	Name    string
	Timeout time.Duration
}

// BlobStorerOption is a function that sets *BlobStorerOptions.
type BlobStorerOption func(o *BlobStorerOptions)

// NewBlobStorer creates a BlobStorer with the provided options.
func NewBlobStorer(options ...BlobStorerOption) (*BlobStorer, error) {
	client, err := dapr.NewClient()
	if err != nil {
		return nil, err
	}

	s := newBlobStorer(options...)
	s.client = client

	return s, nil
}

// newBlobStorer creates a new *BlobStorer with the provided options.
func newBlobStorer(options ...BlobStorerOption) *BlobStorer {
	opts := BlobStorerOptions{
		Name:    defaultStorerName,
		Timeout: defaultStorerTimeout,
	}

	for _, option := range options {
		option(&opts)
	}

	return &BlobStorer{
		name:    opts.Name,
		timeout: opts.Timeout,
	}
}

// Store a report in a blob storage.
func (s BlobStorer) Store(r Report) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if _, err := s.client.InvokeBinding(ctx, &dapr.InvokeBindingRequest{
		Name:      s.name,
		Data:      r.JSON(),
		Operation: "create",
		Metadata: map[string]string{
			"key":      r.ID,
			"blobName": r.ID + ".json",
		},
	}); err != nil {
		return err
	}
	return nil
}
