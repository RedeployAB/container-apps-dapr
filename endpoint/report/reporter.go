package report

import (
	"context"
	"time"

	dapr "github.com/dapr/go-sdk/client"
)

const (
	defaultPubsubReporterName  = "reports"
	defaultPubsubReporterTopic = "create"
)

// PubsubReporter is a reporter that uses Pubsub with DAPR to run reports.
type PubsubReporter struct {
	dapr.Client
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
	opts := PubsubReporterOptions{
		Name:  defaultPubsubReporterName,
		Topic: defaultPubsubReporterTopic,
	}

	for _, o := range options {
		if len(o.Name) > 0 {
			opts.Name = o.Name
		}
		if len(o.Topic) > 0 {
			opts.Topic = o.Topic
		}
	}

	client, err := dapr.NewClient()
	if err != nil {
		return nil, err
	}

	return &PubsubReporter{
		Client:  client,
		address: address,
		name:    opts.Name,
		topic:   opts.Topic,
	}, nil
}

// Run a report routine.
func (r PubsubReporter) Run(report Report) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.PublishEvent(ctx, r.name, r.topic, report.Serialize())
}
