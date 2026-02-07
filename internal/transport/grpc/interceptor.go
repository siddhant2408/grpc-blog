package grpctransport

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func UnaryLoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		start := time.Now()
		resp, err := handler(ctx, req)

		logger.Info("grpc request",
			zap.String("method", info.FullMethod),
			zap.Duration("latency", time.Since(start)),
			zap.Error(err),
		)

		return resp, err
	}
}
