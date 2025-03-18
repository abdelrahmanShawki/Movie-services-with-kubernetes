package grpcratingGateway

import (
	"context"
	"movie-app.com/movie/internal/grpcutility"
	"movie-app.com/pkg/discovery"
	"movie-app.com/src/gen"
)

// Gateway defines a gRPC gateway for a rating service.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a rating service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// GetAggregatedRating returns the aggregated rating for a
// record or ErrNotFound if there are no ratings for it.
func (g *Gateway) GetAggregatedRating(ctx context.Context, recordID int32,
	recordType int32) (float32, error) {

	conn, err := grpcutility.ServiceConnection(ctx, "rating", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := gen.NewRatingServiceClient(conn)
	resp, err := client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{
		RecordId:   string(recordID),
		RecordType: recordType,
	})
	if err != nil {
		return 0, err
	}

	return resp.RatingValue, nil
}
