package report

import (
	"context"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

const (
	defaultReporterTopic = "create"
)

// PubsubReporter is a reporter that uses Pubsub with DAPR to run reports.
type PubsubReporter struct {
	client
	name    string
	topic   string
	timeout time.Duration
}

// PubsubReporterOptions contains settings for a PubsubReporter.
type PubsubReporterOptions struct {
	Name    string
	Topic   string
	Timeout time.Duration
}

// PubsubReporterOption is a function that sets *PubsubReporterOptions.
type PubsubReporterOption func(o *PubsubReporterOptions)

// NewPubsubReporter creates a new *PubsubReporter with the provided address
// and options.
func NewPubsubReporter(options ...PubsubReporterOption) (*PubsubReporter, error) {
	client, err := dapr.NewClient()
	if err != nil {
		return nil, err
	}

	r := newPubsubReporter(options...)
	r.client = client

	return r, nil
}

// newPubsubReporter creates a new *PubsubReporter with the provided address and
// options.
func newPubsubReporter(options ...PubsubReporterOption) *PubsubReporter {
	opts := PubsubReporterOptions{
		Name:    defaultReporterName,
		Topic:   defaultReporterTopic,
		Timeout: defaultReporterTimeout,
	}

	for _, option := range options {
		option(&opts)
	}

	return &PubsubReporter{
		name:    opts.Name,
		topic:   opts.Topic,
		timeout: opts.Timeout,
	}
}

// Run a report routine.
func (r PubsubReporter) Run(report Report) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.PublishEvent(ctx, r.name, r.topic, report.JSON())
}
