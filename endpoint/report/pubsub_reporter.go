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

// NewPubsubReporter creates a new *PubsubReporter with the provided address
// and options.
func NewPubsubReporter(options ...PubsubReporterOptions) (*PubsubReporter, error) {
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
func newPubsubReporter(options ...PubsubReporterOptions) *PubsubReporter {
	opts := PubsubReporterOptions{
		Name:    defaultReporterName,
		Topic:   defaultReporterTopic,
		Timeout: defaultReporterTimeout,
	}

	for _, o := range options {
		if len(o.Name) > 0 {
			opts.Name = o.Name
		}
		if len(o.Topic) > 0 {
			opts.Topic = o.Topic
		}
		if o.Timeout > 0 {
			opts.Timeout = o.Timeout
		}
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
