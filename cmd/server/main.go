package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"grpc-blog/internal/app/blog"
	"grpc-blog/internal/container"
	"grpc-blog/internal/infra/tracing"
	grpcTransport "grpc-blog/internal/transport/grpc"
	grpctransport "grpc-blog/internal/transport/grpc"
	"grpc-blog/proto/blogpb"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// ---- tracing ----
	shutdownTracer, err := tracing.InitTracer()
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}
	defer shutdownTracer(context.Background())

	// ---- DI container ----
	c, err := container.Build()
	if err != nil {
		log.Fatalf("failed to build container: %v", err)
	}

	// ---- start server via DI ----
	err = c.Invoke(func(
		logger *zap.Logger,
		service *blog.Service,
	) {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			logger.Fatal("failed to listen", zap.Error(err))
		}

		grpcServer := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				grpctransport.UnaryLoggingInterceptor(logger),
			),
			grpc.StatsHandler(otelgrpc.NewServerHandler()),
		)

		blogpb.RegisterBlogServiceServer(
			grpcServer,
			grpcTransport.NewBlogGRPCServer(service),
		)

		logger.Info("gRPC server started", zap.String("addr", ":50051"))

		// ---- graceful shutdown ----
		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				logger.Fatal("grpc server crashed", zap.Error(err))
			}
		}()

		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop

		logger.Info("shutting down gRPC server")
		grpcServer.GracefulStop()
	})

	if err != nil {
		log.Fatalf("failed to invoke container: %v", err)
	}
}
