package main

import (
	"google.golang.org/grpc"
	"log"
	"movie-app.com/metadata/internal/controller/metadata"
	"movie-app.com/metadata/internal/handler/grpchandler"
	"movie-app.com/metadata/internal/repository/memory"
	"movie-app.com/src/gen"
	"net"
)

func main() {
	log.Println("Starting the movie metadata service")
	repo := memory.New()
	svc := metadata.New(repo)
	h := grpchandler.New(svc)
	lis, err := net.Listen("tcp", "localhost:8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	gen.RegisterMetadataServiceServer(srv, h)
	srv.Serve(lis)
}
