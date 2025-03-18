package grpcmetadataGateway

import (
	"context"
	"movie-app.com/metadata/pkg/metadatamodel"
	"movie-app.com/movie/internal/grpcutility"
	"movie-app.com/pkg/discovery"
	"movie-app.com/src/gen"
)

// Gateway defines a movie metadata gRPC gateway.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a movie metadata service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

// Get returns movie metadata by a movie id.
func (g *Gateway) Get(ctx context.Context, id string) (*metadatamodel.Metadata, error) {
	conn, err := grpcutility.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &gen.GetMetadataRequest{MovieId: id})
	if err != nil {
		return nil, err
	}

	return metadatamodel.MetadataFromProto(resp.Metadata), nil
}
