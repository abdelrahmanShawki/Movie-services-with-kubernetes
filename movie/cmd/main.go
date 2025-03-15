package main

import (
	"log"
	"movie-app.com/movie/internal/controller/movie"
	"movie-app.com/movie/internal/gateway/metadataGateway"
	"movie-app.com/movie/internal/gateway/ratingGateway"
	"movie-app.com/movie/internal/handler/httphandler"
	"net/http"
)

func main() {
	log.Println("Starting the movie service")

	metadataGateway := metadataGateway.New("localhost:8081")
	ratingGateway := ratingGateway.New("localhost:8082")

	ctrl := movie.New(ratingGateway, metadataGateway)
	h := httphandler.New(ctrl)

	http.Handle("/movie", http.HandlerFunc(h.GetMovieDetails))

	if err := http.ListenAndServe(":8083", nil); err != nil {
		panic(err)
	}
}
