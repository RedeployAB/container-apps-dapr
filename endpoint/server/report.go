package server

import "encoding/json"

// Report is a incoming report request.
type Report struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

// JSON returns the JSON representation of the message.
func (r Report) JSON() []byte {
	b, _ := json.Marshal(&r)
	return b
}
