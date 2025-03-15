package memory

import (
	"context"
	"movie-app.com/rating/internal/repository"
	"movie-app.com/rating/pkg/ratingmodel"
)

// Repository defines a rating repository.
type Repository struct {
	data map[ratingmodel.RecordType]map[ratingmodel.RecordID][]ratingmodel.Rating
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{map[ratingmodel.RecordType]map[ratingmodel.RecordID][]ratingmodel.Rating{}}
}

func (r *Repository) Get(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType) ([]ratingmodel.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}
	if ratings, ok := r.data[recordType][recordID]; !ok ||
		len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}
	return r.data[recordType][recordID], nil
}

// Put adds a rating for a given record.
func (r *Repository) Put(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error {
	// Check if there's an entry for this record type. If not, initialize it.
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[ratingmodel.RecordID][]ratingmodel.Rating{}
	}
	// Append the new rating to the slice for the specific record ID.
	r.data[recordType][recordID] = append(r.data[recordType][recordID], *rating)
	return nil
}
