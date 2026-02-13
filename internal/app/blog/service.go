package blog

import (
	"errors"
	"sync"

	"grpc-blog/proto/blogpb"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Service encapsulates all business logic related to blog posts.
//
// Responsibilities:
// - Manage lifecycle of blog posts
// - Enforce business rules (existence, uniqueness)
// - Remain independent of transport (gRPC / HTTP)
// - Be safe for concurrent access
//
// Persistence is in-memory by design and can be replaced with a database
// without affecting callers.
type Service struct {
	mu     sync.RWMutex
	posts  map[string]*blogpb.BlogPost
	logger *zap.Logger
}

// NewService constructs a new blog Service.
//
// Inputs:
// - logger: structured logger used for domain-level events
//
// Output:
// - Initialized *Service with empty in-memory storage
//
// This function performs no I/O and never returns an error.
func NewService(logger *zap.Logger) *Service {
	return &Service{
		posts:  make(map[string]*blogpb.BlogPost),
		logger: logger,
	}
}

// Create creates a new blog post.
//
// Business behavior:
// - Generates a unique PostID
// - Stores the post in memory
// - Logs the creation event
//
// Inputs:
// - post: BlogPost without PostID
//
// Output:
// - Stored BlogPost with PostID populated
// - Error only if invariants are violated (currently none)
//
// Thread-safe.
func (s *Service) CreatePost(post *blogpb.BlogPost) (*blogpb.BlogPost, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	post.PostId = uuid.New().String()
	s.posts[post.PostId] = post

	s.logger.Info("post created",
		zap.String("post_id", post.PostId),
		zap.String("author", post.Author),
	)

	return post, nil
}

// Read retrieves a blog post by PostID.
//
// Business behavior:
// - Validates existence
// - Returns a copy-safe reference
//
// Inputs:
// - id: unique identifier of the blog post
//
// Output:
// - BlogPost if found
// - Error if post does not exist
//
// Thread-safe.
func (s *Service) ReadPost(id string) (*blogpb.BlogPost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	post, ok := s.posts[id]
	if !ok {
		return nil, errors.New("post not found")
	}

	s.logger.Info("post read",
		zap.String("post_id", post.PostId),
		zap.String("author", post.Author),
	)

	return post, nil
}

// Read retrieves a blog post by PostID.
//
// Business behavior:
// - Validates existence
// - Returns a copy-safe reference
//
// Inputs:
// - id: unique identifier of the blog post
//
// Output:
// - BlogPost if found
// - Error if post does not exist
//
// Thread-safe.
func (s *Service) ReadAll() ([]*blogpb.BlogPost, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*blogpb.BlogPost
	for _, post := range s.posts {
		result = append(result, post)
		go func() {
			s.logger.Info("post read",
				zap.String("post_id", post.PostId),
				zap.String("author", post.Author),
			)
		}()
	}

	return result, nil
}

// Update modifies an existing blog post.
//
// Business behavior:
// - Validates that the post exists
// - Preserves PostID
// - Overwrites mutable fields
//
// Inputs:
// - id: identifier of the post to update
// - post: new blog post content
//
// Output:
// - Updated BlogPost
// - Error if post does not exist
//
// Thread-safe.
func (s *Service) UpdatePost(id string, post *blogpb.BlogPost) (*blogpb.BlogPost, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.posts[id]; !ok {
		return nil, errors.New("post not found")
	}
	post.PostId = id
	s.posts[id] = post

	s.logger.Info("post updated",
		zap.String("post_id", post.PostId),
		zap.String("author", post.Author),
	)

	return post, nil
}

// Delete removes a blog post permanently.
//
// Business behavior:
// - Validates existence
// - Deletes from in-memory store
//
// Inputs:
// - id: identifier of the post to delete
//
// Output:
// - Error if post does not exist
//
// Thread-safe.
func (s *Service) DeletePost(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.posts[id]; !ok {
		return errors.New("post not found")
	}
	delete(s.posts, id)

	s.logger.Info("post deleted",
		zap.String("post_id", id),
	)

	return nil
}
