package httpmetadataGateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"movie-app.com/metadata/pkg/metadatamodel"
	"movie-app.com/movie/internal/gateway"
	"movie-app.com/pkg/discovery"
	"net/http"
)

// Gateway defines a movie metadata HTTP gateway.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new HTTP gateway for a movie metadata
// service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

// Get gets movie metadata by a movie id.
func (g *Gateway) Get(ctx context.Context, id string) (*metadatamodel.Metadata, error) {

	addrs, err := g.registry.ServiceAddresses(ctx, "metadata")
	if err != nil {
		return nil, err
	}

	url := "http://" + addrs[rand.Intn(len(addrs))] + "/metadata"
	log.Printf("Calling metadata service. Request: GET " + url)

	req, err := http.NewRequest(http.MethodGet, url, nil)

	values := req.URL.Query()
	values.Add("id", id)
	req.URL.RawQuery = values.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, gateway.ErrNotFound
	} else if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("non-2xx response: %v", resp)
	}

	var v *metadatamodel.Metadata
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return v, nil
}
