package report

import (
	"encoding/json"
)

// Report represents a report with an ID and data.
type Report struct {
	ID   string
	Data []byte
}

// NewReport creates a new Report.
func NewReport(id string, data []byte) Report {
	return Report{
		ID:   id,
		Data: data,
	}
}

// JSON returns a JSON representation of a Report.
func (r Report) JSON() []byte {
	b, _ := json.Marshal(r)
	return b
}
