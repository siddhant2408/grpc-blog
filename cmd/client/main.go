package main

import (
	"context"
	"log"
	"time"

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	"grpc-blog/internal/infra/tracing"
	"grpc-blog/proto/blogpb"
)

func main() {
	// ---- tracing ----
	shutdown, err := tracing.InitTracer()
	if err != nil {
		log.Fatalf("failed to init tracer: %v", err)
	}
	defer shutdown(context.Background())

	// ---- logger ----
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// ---- context ----
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// ---- grpc client ----
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Fatal("failed to create grpc client", zap.Error(err))
	}
	defer conn.Close()

	client := blogpb.NewBlogServiceClient(conn)

	// ---- call API ----
	resp, err := client.CreatePost(ctx, &blogpb.CreatePostRequest{
		Title:           "Create Post",
		Content:         "First Post",
		Author:          "Siddhant",
		PublicationDate: timestamppb.New(time.Now()),
		Tags:            []string{"create_post"},
	})
	if err != nil {
		logger.Fatal("CreatePost failed", zap.Error(err))
	}

	logger.Info("post created",
		zap.String("post_id", resp.Post.PostId),
		zap.String("title", resp.Post.Title),
	)
}
