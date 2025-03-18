package rating

import (
	"context"
	"errors"
	"movie-app.com/rating/internal/repository"
	"movie-app.com/rating/pkg/ratingmodel"
)

var ErrNotFound = errors.New("ratings not found for a record")

type ratingRepository interface {
	Get(ctx context.Context, recordID ratingmodel.RecordID,
		recordType ratingmodel.RecordType) ([]ratingmodel.Rating, error)
	Put(ctx context.Context, recordID ratingmodel.RecordID,
		recordType ratingmodel.RecordType, rating *ratingmodel.Rating) error
}

// Controller defines a rating service controller.
type Controller struct {
	repo ratingRepository
}

// New creates a rating service controller.
func New(repo ratingRepository) *Controller {
	return &Controller{repo}
}

// GetAggregatedRating returns the aggregated rating for a record
// or ErrNotFound if there are no ratings for it.
func (c *Controller) GetAggregatedRating(
	ctx context.Context,
	recordID ratingmodel.RecordID,
	recordType ratingmodel.RecordType,
) (float32, error) {

	// Fetch ratings from the repository
	ratings, err := c.repo.Get(ctx, recordID, recordType)

	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return 0, ErrNotFound
		}
		return 0, err
	}

	// Compute the sum of ratings
	var sum float64
	for _, r := range ratings {
		sum += float64(r.Value)
	}

	// Return the average rating
	return sum / float64(len(ratings)), nil
}

// PutRating writes a rating for a given record.
func (c *Controller) PutRating(ctx context.Context, recordID ratingmodel.RecordID, recordType ratingmodel.RecordType,
	rating *ratingmodel.Rating) error {
	return c.repo.Put(ctx, recordID, recordType, rating)
}
