package blog

import (
	"testing"

	"grpc-blog/proto/blogpb"

	"go.uber.org/zap/zaptest"
)

func newTestService(t *testing.T) *Service {
	t.Helper()

	logger := zaptest.NewLogger(t)
	return NewService(logger)
}

func TestNewService(t *testing.T) {
	svc := newTestService(t)
	if svc == nil {
		t.Fatal("expected service to be non-nil")
	}
	if len(svc.posts) != 0 {
		t.Fatal("expected empty post store")
	}
}

func TestCreate(t *testing.T) {
	svc := newTestService(t)

	post := &blogpb.BlogPost{
		Title:   "title",
		Content: "content",
		Author:  "author",
	}

	created, err := svc.CreatePost(post)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if created.PostId == "" {
		t.Fatal("expected PostId to be set")
	}

	if svc.posts[created.PostId] == nil {
		t.Fatal("post not stored")
	}
}

func TestReadSuccess(t *testing.T) {
	svc := newTestService(t)

	post, _ := svc.CreatePost(&blogpb.BlogPost{Title: "read-test"})

	read, err := svc.ReadPost(post.PostId)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if read.PostId != post.PostId {
		t.Fatal("returned wrong post")
	}
}

func TestReadNotFound(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.ReadPost("non-existent-id")
	if err == nil {
		t.Fatal("expected error for missing post")
	}
}

func TestUpdateSuccess(t *testing.T) {
	svc := newTestService(t)

	created, _ := svc.CreatePost(&blogpb.BlogPost{Title: "old"})

	updatedPost := &blogpb.BlogPost{
		Title:   "new",
		Content: "updated",
		Author:  "author",
	}

	updated, err := svc.UpdatePost(created.PostId, updatedPost)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if updated.PostId != created.PostId {
		t.Fatal("post id should not change")
	}

	if updated.Title != "new" {
		t.Fatal("post was not updated")
	}
}

func TestUpdateNotFound(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.UpdatePost("missing-id", &blogpb.BlogPost{})
	if err == nil {
		t.Fatal("expected error for updating missing post")
	}
}

func TestDeleteSuccess(t *testing.T) {
	svc := newTestService(t)

	created, _ := svc.CreatePost(&blogpb.BlogPost{Title: "delete-test"})

	err := svc.DeletePost(created.PostId)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, exists := svc.posts[created.PostId]; exists {
		t.Fatal("post should be deleted")
	}
}

func TestDeleteNotFound(t *testing.T) {
	svc := newTestService(t)

	err := svc.DeletePost("missing-id")
	if err == nil {
		t.Fatal("expected error for deleting missing post")
	}
}
