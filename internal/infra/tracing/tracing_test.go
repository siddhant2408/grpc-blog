package tracing

import (
	"context"
	"testing"
)

func TestInitTracer(t *testing.T) {
	shutdown, err := InitTracer()
	if err != nil {
		t.Fatalf("failed to init tracer: %v", err)
	}

	if shutdown == nil {
		t.Fatal("expected shutdown func")
	}

	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown failed: %v", err)
	}
}
