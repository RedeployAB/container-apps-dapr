package report

import (
	"context"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

const (
	defaultPubsubReporterName    = "reports"
	defaultPubsubReporterTopic   = "create"
	defaultPubsubReporterTimeout = time.Second * 10
)

// client is the interface that wraps around method PublishEvent.
type client interface {
	PublishEvent(ctx context.Context, pubsubName, topic string, data any, options ...dapr.PublishEventOption) error
}

// PubsubReporter is a reporter that uses Pubsub with DAPR to run reports.
type PubsubReporter struct {
	client
	address string
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
func NewPubsubReporter(address string, options ...PubsubReporterOptions) (*PubsubReporter, error) {
	client, err := dapr.NewClient()
	if err != nil {
		return nil, err
	}

	r := newPubsubReporter(address, options...)
	r.client = client

	return r, nil
}

// newPubsubReporter creates a new *PubsubReporter with the provided address and
// options.
func newPubsubReporter(address string, options ...PubsubReporterOptions) *PubsubReporter {
	opts := PubsubReporterOptions{
		Name:    defaultPubsubReporterName,
		Topic:   defaultPubsubReporterTopic,
		Timeout: defaultPubsubReporterTimeout,
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
		address: address,
		name:    opts.Name,
		topic:   opts.Topic,
		timeout: opts.Timeout,
	}
}

// Run a report routine.
func (r PubsubReporter) Run(report Report) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.PublishEvent(ctx, r.name, r.topic, report.Serialize())
}