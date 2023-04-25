package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/RedeployAB/container-apps-dapr/common/report"
	"github.com/dapr/go-sdk/service/common"
)

// reportHandler is the handler for the report topic.
func (s server) reportHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	s.log.Info("Event received.", "id", e.ID, "pubsub", e.PubsubName, "topic", e.Topic)

	data, ok := e.Data.(string)
	if !ok {
		s.log.Error(fmt.Errorf("Failed to convert data to string."), "Failed to cast data to string.", "id", e.ID, "pubsub", e.PubsubName, "topic", e.Topic)
		return false, err
	}

	var r report.Report
	if err := json.Unmarshal([]byte(data), &r); err != nil {
		s.log.Error(err, "Failed to deserialize report.", "id", e.ID, "pubsub", e.PubsubName, "topic", e.Topic)
		return false, err
	}

	if err := s.reporter.Create(r); err != nil {
		s.log.Error(err, "Failed to create report.", "id", e.ID, "pubsub", e.PubsubName, "topic", e.Topic)
		return false, err
	}
	s.log.Info("Report created.", "id", e.ID, "pubsub", e.PubsubName, "topic", e.Topic)

	return false, nil
}
