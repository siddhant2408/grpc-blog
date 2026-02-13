package main

import (
	"context"
	"flag"
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

	// Define the flags. Each function takes the flag name, default value, and a help message.
	opType := flag.String("type", "fetch", "the type of operation")
	postID := flag.String("id", "", "the id to fetch")

	// Parse the command line arguments
	flag.Parse()

	switch *opType {
	case "create":
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
			zap.String("post_id", resp.Post[0].PostId),
			zap.String("title", resp.Post[0].Title),
		)

	case "fetch":
		// ---- call API --i--
		resp, err := client.ReadPost(ctx, &blogpb.ReadPostRequest{
			PostId: *postID,
		})
		if err != nil {
			logger.Fatal("FetchPost failed", zap.Error(err))
		}

		logger.Info("post fetched",
			zap.String("post_id", resp.Post[0].PostId),
			zap.String("title", resp.Post[0].Title),
		)

	case "fetchall":
		// ---- call API --i--
		resp, err := client.ReadAll(ctx, &blogpb.ReadAllRequest{})
		if err != nil {
			logger.Fatal("FetchPost failed", zap.Error(err))
		}

		for _, post := range resp.Post {
			logger.Info("post fetched",
				zap.String("post_id", post.PostId),
				zap.String("title", post.Title),
			)
		}

	case "update":

		// ---- call API ----
		resp, err := client.UpdatePost(ctx, &blogpb.UpdatePostRequest{
			PostId: *postID,
			Title:  "updated title",
		})
		if err != nil {
			logger.Fatal("FetchPost failed", zap.Error(err))
		}

		logger.Info("post updated",
			zap.String("post_id", resp.Post[0].PostId),
			zap.String("title", resp.Post[0].Title),
		)

	case "delete":
		_, err = client.DeletePost(ctx, &blogpb.DeletePostRequest{
			PostId: *postID,
		})
		if err != nil {
			logger.Fatal("FetchPost failed", zap.Error(err))
		}

		logger.Info("post deleted",
			zap.String("post_id", *postID),
		)
	}
}
