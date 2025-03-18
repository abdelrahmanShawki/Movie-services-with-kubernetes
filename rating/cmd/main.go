package main

import (
	"google.golang.org/grpc"
	"log"
	"movie-app.com/rating/internal/controller/rating"
	"movie-app.com/rating/internal/handler/grpchandler"
	"movie-app.com/rating/internal/repository/memory"
	"movie-app.com/src/gen"
	"net"
)

func main() {
	log.Println("Starting the rating service")
	repo := memory.New()
	svc := rating.New(repo)
	h := grpchandler.New(svc)
	lis, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterRatingServiceServer(srv, h)
	srv.Serve(lis)
}
