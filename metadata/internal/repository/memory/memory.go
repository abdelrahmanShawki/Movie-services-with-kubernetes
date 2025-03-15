package memory

import (
	"context"
	"movie-app.com/metadata/internal/repository"
	"movie-app.com/metadata/pkg/metadatamodel"
	"sync"
)

// Repository defines a memory movie metadata repository.
type Repository struct {
	sync.RWMutex
	data map[string]*metadatamodel.Metadata
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{data: map[string]*metadatamodel.Metadata{}}
}

// Get retrieves movie metadata for by movie id.
func (r *Repository) Get(_ context.Context, id string) (*metadatamodel.Metadata, error) {
	r.RLock()
	defer r.RUnlock()
	m, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(_ context.Context, id string, metadata *metadatamodel.Metadata) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata
	return nil
}
