package server

import (
	"context"
	"encoding/json"

	"github.com/RedeployAB/container-apps-dapr/worker/report"
	"github.com/dapr/go-sdk/service/common"
)

// queueReportHandler is the handler for the report queue.
func (s server) queueReportHandler(ctx context.Context, in *common.BindingEvent) (out []byte, err error) {
	s.log.Info("Message received.", "metadata", in.Metadata)

	var r report.Report
	if err := json.Unmarshal(in.Data, &r); err != nil {
		s.log.Error(err, "Failed to deserialize report.", "metadata", in.Metadata)
		return nil, err
	}

	if err := s.reporter.Create(r); err != nil {
		s.log.Error(err, "Failed to create report.", "metadata", in.Metadata)
		return nil, err
	}

	return []byte(`Message processed.`), nil
}
