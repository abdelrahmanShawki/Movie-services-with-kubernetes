package model

import "movie-app.com/metadata/pkg/metadatamodel"

// MovieDetails includes movie metadata its aggregated
// rating.
type MovieDetails struct {
	Rating   *float64               `json:"rating,omitEmpty"`
	Metadata metadatamodel.Metadata `json:"metadata"`
}

// Get returns the movie details including the aggregated rating and movie metadata.
