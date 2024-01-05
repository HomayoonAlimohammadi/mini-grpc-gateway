package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/HomayoonAlimohammadi/mini-grpc-gateway/pb/post"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type serverImpl struct {
	post.UnimplementedPostServiceServer
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	post.RegisterPostServiceServer(grpcServer, newServer())
	reflection.Register(grpcServer)

	fmt.Println("running server on :50051")
	log.Fatal(grpcServer.Serve(lis))
}

func newServer() *serverImpl {
	return &serverImpl{}
}

func (s *serverImpl) GetPost(ctx context.Context, _ *post.Empty) (*post.GetPostResponse, error) {
	return &post.GetPostResponse{Title: "title", Description: "description"}, nil
}