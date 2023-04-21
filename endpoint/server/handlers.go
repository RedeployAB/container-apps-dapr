package server

import (
	"encoding/json"
	"net/http"

	"github.com/RedeployAB/container-apps-dapr/common/report"
)

// reportHandler returns a handler for incoming reports.
func (s server) reportHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var re Report
		if err := json.NewDecoder(r.Body).Decode(&re); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		s.log.Info("Received report.", "handler", "report", "id", re.ID, "data", string(re.Data))

		if err := s.reporter.Create(report.NewReport(re.ID, []byte(re.Data))); err != nil {
			s.log.Error(err, "Error creating report.")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		s.log.Info("Report sent for creation.", "handler", "report", "id", re.ID, "data", string(re.Data))

		w.WriteHeader(http.StatusOK)
		w.Write(re.JSON())
	})
}