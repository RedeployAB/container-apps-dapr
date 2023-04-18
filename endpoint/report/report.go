package report

import (
	"bytes"
	"encoding/gob"
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

// Serialize Report as gob as []byte.
func (r Report) Serialize() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(r)
	return buf.Bytes()
}

// Desierialize a gob from []byte to Report.
func (r *Report) Deserialize(b []byte) error {
	enc := gob.NewDecoder(bytes.NewBuffer(b))
	return enc.Decode(&r)
}
