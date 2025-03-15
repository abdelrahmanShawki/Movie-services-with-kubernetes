package main

import (
	"log"
	"movie-app.com/rating/internal/controller/rating"
	"movie-app.com/rating/internal/handler/httphandler"
	"movie-app.com/rating/internal/repository/memory"
	"net/http"
)

func main() {
	log.Println("Starting the rating service")
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/rating", http.HandlerFunc(h.Handle))
	if err := http.ListenAndServe(":8082", nil); err != nil {
		panic(err)
	}
}
