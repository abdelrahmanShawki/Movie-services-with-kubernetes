package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"movie-app.com/movie/internal/controller/movie"
	"movie-app.com/movie/internal/gateway/metadataGateway"
	"movie-app.com/movie/internal/gateway/ratingGateway"
	"movie-app.com/movie/internal/handler/httphandler"
	"movie-app.com/pkg/discovery"
	"movie-app.com/pkg/discovery/consul"
	"net/http"
	"time"
)

const serviceName = "movie"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()
	log.Printf("Starting the movie service on port %d", port)
	registry, err := consul.NewRegistry("localhost:8500")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID,
		serviceName, fmt.Sprintf("localhost:%d", port)); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.
				ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName)

	metadataGtway := metadataGateway.New(registry)
	ratingGtway := ratingGateway.New(registry)

	ctrl := movie.New(ratingGtway, metadataGtway)
	h := httphandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
