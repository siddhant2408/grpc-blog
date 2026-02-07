package grpctransport

import (
	"context"
	"errors"
	"testing"

	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
)

func TestUnaryLoggingInterceptor(t *testing.T) {
	logger := zaptest.NewLogger(t)
	interceptor := UnaryLoggingInterceptor(logger)

	handlerCalled := false

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		handlerCalled = true
		return "ok", nil
	}

	_, err := interceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test"},
		handler,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !handlerCalled {
		t.Fatal("handler was not called")
	}
}

func TestUnaryLoggingInterceptor_Error(t *testing.T) {
	logger := zaptest.NewLogger(t)
	interceptor := UnaryLoggingInterceptor(logger)

	expectedErr := errors.New("boom")

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, expectedErr
	}

	_, err := interceptor(
		context.Background(),
		nil,
		&grpc.UnaryServerInfo{FullMethod: "/test"},
		handler,
	)

	if err != expectedErr {
		t.Fatal("expected handler error to propagate")
	}
}
