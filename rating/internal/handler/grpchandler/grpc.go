package grpchandler

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movie-app.com/rating/internal/controller/rating"
	"movie-app.com/rating/pkg/ratingmodel"
	"movie-app.com/src/gen"
)

// Handler defines a gRPC rating API handler.”
type Handler struct {
	gen.UnimplementedRatingServiceServer
	svc *rating.Controller
}

// New creates a new movie metadata gRPC handler
func New(ctrl *rating.Controller) *Handler {
	return &Handler{svc: ctrl}
}

// GetAggregatedRating returns the aggregated rating for a
// record.
func (h *Handler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.RecordType == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	v, err := h.svc.GetAggregatedRating(ctx, ratingmodel.RecordID(req.RecordId), ratingmodel.RecordType(req.RecordType))
	if err != nil && errors.Is(err, rating.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &gen.GetAggregatedRatingResponse{RatingValue: v}, nil
}

// PutRating writes a rating for a given record.
func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty user id or record id")
	}
	if err := h.svc.PutRating(
		ctx,
		ratingmodel.RecordID(req.RecordId),
		ratingmodel.RecordType(req.RecordType),
		&ratingmodel.Rating{
			UserID: ratingmodel.UserID(req.UserId),
			Value:  ratingmodel.RatingValue(req.RatingValue),
		},
	); err != nil {
		return nil, err
	}
	return &gen.PutRatingResponse{}, nil
}
