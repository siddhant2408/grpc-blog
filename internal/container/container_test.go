package container

import (
	"testing"

	"go.uber.org/zap"

	"grpc-blog/internal/app/blog"
)

func TestBuildContainer(t *testing.T) {
	c, err := Build()
	if err != nil {
		t.Fatalf("failed to build container: %v", err)
	}

	err = c.Invoke(func(
		logger *zap.Logger,
		service *blog.Service,
	) {
		if logger == nil || service == nil {
			t.Fatal("dependencies not resolved")
		}
	})

	if err != nil {
		t.Fatalf("failed to invoke container: %v", err)
	}
}
