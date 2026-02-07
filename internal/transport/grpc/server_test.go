package grpctransport

import (
	"context"
	"net"
	"testing"

	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"grpc-blog/internal/app/blog"
	"grpc-blog/proto/blogpb"
)

const bufSize = 1024 * 1024

func setupTestGRPCServer(t *testing.T) (*grpc.ClientConn, func()) {
	t.Helper()

	lis := bufconn.Listen(bufSize)

	logger := zaptest.NewLogger(t)
	service := blog.NewService(logger)
	server := grpc.NewServer()

	blogpb.RegisterBlogServiceServer(
		server,
		NewBlogGRPCServer(service),
	)

	errCh := make(chan error, 1)

	go func() {
		if err := server.Serve(lis); err != nil {
			errCh <- err
		}
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}

	cleanup := func() {
		server.Stop()
		conn.Close()

		select {
		case err := <-errCh:
			// gRPC returns ErrServerStopped on normal shutdown — ignore it
			if err != grpc.ErrServerStopped {
				t.Fatalf("server exited unexpectedly: %v", err)
			}
		default:
		}
	}

	return conn, cleanup
}

func TestCreatePost(t *testing.T) {
	conn, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	client := blogpb.NewBlogServiceClient(conn)

	resp, err := client.CreatePost(context.Background(), &blogpb.CreatePostRequest{
		Title:   "test",
		Content: "content",
		Author:  "author",
	})
	if err != nil {
		t.Fatalf("CreatePost failed: %v", err)
	}

	if resp.Post == nil {
		t.Fatal("expected post in response")
	}

	if resp.Post.PostId == "" {
		t.Fatal("expected post_id to be set")
	}
}

func TestReadPostNotFound(t *testing.T) {
	conn, cleanup := setupTestGRPCServer(t)
	defer cleanup()

	client := blogpb.NewBlogServiceClient(conn)

	resp, err := client.ReadPost(context.Background(), &blogpb.ReadPostRequest{
		PostId: "missing",
	})
	if err != nil {
		t.Fatalf("ReadPost failed: %v", err)
	}

	if resp.Error == "" {
		t.Fatal("expected error for missing post")
	}
}
