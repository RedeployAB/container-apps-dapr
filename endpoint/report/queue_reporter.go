package report

import (
	"context"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

const (
	defaultReporterQueue = "create"
)

// QueueReporter is a reporter that uses queue binding with DAPR to run reports.
type QueueReporter struct {
	client
	name    string
	queue   string
	timeout time.Duration
}

// QueueReporterOptions contains settings for a QueueReporter.
type QueueReporterOptions struct {
	Name    string
	Queue   string
	Timeout time.Duration
}

// QueueReporterOption is a function that sets *QueueReporterOptions.
type QueueReporterOption func(o *QueueReporterOptions)

// NewQueueReporter creates a new *QueueReporter with the provided address
// and options.
func NewQueueReporter(options ...QueueReporterOption) (*QueueReporter, error) {
	client, err := dapr.NewClient()
	if err != nil {
		return nil, err
	}

	r := newQueueReporter(options...)
	r.client = client

	return r, nil
}

// newQueueReporter creates a new *QueueReporter with the provided address and
// options.
func newQueueReporter(options ...QueueReporterOption) *QueueReporter {
	opts := QueueReporterOptions{
		Name:    defaultReporterName,
		Queue:   defaultReporterQueue,
		Timeout: defaultReporterTimeout,
	}

	for _, option := range options {
		option(&opts)
	}

	return &QueueReporter{
		name:    opts.Name,
		queue:   opts.Queue,
		timeout: opts.Timeout,
	}
}

// Run a report routine.
func (r QueueReporter) Run(report Report) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.InvokeOutputBinding(ctx, &dapr.InvokeBindingRequest{
		Name:      r.name,
		Operation: "create",
		Data:      report.JSON(),
		Metadata: map[string]string{
			"queueName": r.queue,
		},
	})
}
