package grpctransport

import (
	"context"

	"grpc-blog/internal/app/blog"
	"grpc-blog/proto/blogpb"
)

type BlogGRPCServer struct {
	blogpb.UnimplementedBlogServiceServer
	service *blog.Service
}

func NewBlogGRPCServer(service *blog.Service) *BlogGRPCServer {
	return &BlogGRPCServer{service: service}
}

func (s *BlogGRPCServer) CreatePost(
	ctx context.Context,
	req *blogpb.CreatePostRequest,
) (*blogpb.PostResponse, error) {

	post := &blogpb.BlogPost{
		Title:           req.Title,
		Content:         req.Content,
		Author:          req.Author,
		PublicationDate: req.PublicationDate,
		Tags:            req.Tags,
	}

	created, err := s.service.CreatePost(post)
	if err != nil {
		return &blogpb.PostResponse{
			Error: err.Error(),
		}, nil
	}

	return &blogpb.PostResponse{
		Post: created,
	}, nil
}

func (s *BlogGRPCServer) ReadPost(
	ctx context.Context,
	req *blogpb.ReadPostRequest,
) (*blogpb.PostResponse, error) {

	post, err := s.service.ReadPost(req.PostId)
	if err != nil {
		return &blogpb.PostResponse{
			Error: err.Error(),
		}, nil
	}

	return &blogpb.PostResponse{
		Post: post,
	}, nil
}

func (s *BlogGRPCServer) UpdatePost(
	ctx context.Context,
	req *blogpb.UpdatePostRequest,
) (*blogpb.PostResponse, error) {

	post := &blogpb.BlogPost{
		Title:   req.Title,
		Content: req.Content,
		Author:  req.Author,
		Tags:    req.Tags,
	}

	updated, err := s.service.UpdatePost(req.PostId, post)
	if err != nil {
		return &blogpb.PostResponse{
			Error: err.Error(),
		}, nil
	}

	return &blogpb.PostResponse{
		Post: updated,
	}, nil
}

func (s *BlogGRPCServer) DeletePost(
	ctx context.Context,
	req *blogpb.DeletePostRequest,
) (*blogpb.DeletePostResponse, error) {

	err := s.service.DeletePost(req.PostId)
	if err != nil {
		return &blogpb.DeletePostResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &blogpb.DeletePostResponse{
		Success: true,
	}, nil
}
