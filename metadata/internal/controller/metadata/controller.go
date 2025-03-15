package metadata

import (
	"context"
	"errors"
	"movie-app.com/metadata/internal/repository"
	"movie-app.com/metadata/pkg/metadatamodel"
)

// ErrNotFound is returned when a requested record is not found.
var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	Get(ctx context.Context, id string) (*metadatamodel.Metadata, error)
}

// Controller defines a metadata service controller.
type Controller struct {
	repo metadataRepository
}

// New creates a metadata service controller.
func New(repo metadataRepository) *Controller {
	return &Controller{repo}
}

// Get returns movie metadata by id.
func (c *Controller) Get(ctx context.Context, id string) (*metadatamodel.Metadata, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}
