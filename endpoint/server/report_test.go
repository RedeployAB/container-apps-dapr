package server

import "testing"

// TestReport_JSON tests the JSON method.
func TestReport_JSON(t *testing.T) {
	report := Report{
		ID:   "123",
		Data: "data",
	}
	expected := `{"id":"123","data":"data"}`

	if string(report.JSON()) != expected {
		t.Errorf("got %s, want %s", report.JSON(), expected)
	}
}
