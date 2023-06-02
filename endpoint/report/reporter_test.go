package report

import (
	"context"

	dapr "github.com/dapr/go-sdk/client"
)

type mockClient struct {
	err error
}

func (c *mockClient) InvokeOutputBinding(ctx context.Context, in *dapr.InvokeBindingRequest) error {
	if c.err != nil {
		return c.err
	}
	return nil
}

func (c *mockClient) PublishEvent(ctx context.Context, pubsubName, topic string, data any, options ...dapr.PublishEventOption) error {
	if c.err != nil {
		return c.err
	}
	return nil
}
