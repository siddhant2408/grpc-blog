package logging

import "testing"

func TestNewLogger(t *testing.T) {
	logger, err := NewLogger()
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
}
