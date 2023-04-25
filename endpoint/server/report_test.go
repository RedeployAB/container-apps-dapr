package server

import "testing"

func TestReport_JSON(t *testing.T) {
	report := Report{
		ID:   "123",
		Data: []byte("data"),
	}
	expected := `{"id":"123","data":"ZGF0YQ=="}`

	if string(report.JSON()) != expected {
		t.Errorf("got %s, want %s", report.JSON(), expected)
	}
}
