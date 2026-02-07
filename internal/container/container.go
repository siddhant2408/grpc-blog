package container

import (
	"go.uber.org/dig"

	"grpc-blog/internal/app/blog"
	"grpc-blog/internal/infra/logging"
)

func Build() (*dig.Container, error) {
	c := dig.New()

	if err := c.Provide(logging.NewLogger); err != nil {
		return nil, err
	}
	if err := c.Provide(blog.NewService); err != nil {
		return nil, err
	}

	return c, nil
}
