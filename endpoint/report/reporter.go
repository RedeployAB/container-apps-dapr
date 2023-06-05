package report

import (
	"context"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

const (
	defaultReporterName    = "reports"
	defaultReporterTimeout = time.Second * 10
)

// client is the interface that wraps around method InvokeOutputBinding and PublishEvent.
type client interface {
	InvokeOutputBinding(ctx context.Context, in *dapr.InvokeBindingRequest) error
	PublishEvent(ctx context.Context, pubsubName, topic string, data any, options ...dapr.PublishEventOption) error
}
