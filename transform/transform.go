package main

import (
	"context"
	"log"
	"net"

	pb "shorturl/transform/transform"

	"github.com/tal-tech/go-zero/core/hash"
	"google.golang.org/grpc"
)

type Transformer struct {
	pb.UnimplementedTransformServer
}

func (t *Transformer) Shorten(ctx context.Context, in *pb.ShortenRequest) (*pb.ShortenResponse, error) {
	key := hash.Md5Hex([]byte(in.Url))[:6]
	return &pb.ShortenResponse{
		Shorten: key,
	}, nil
}

func main() {
	listen, err := net.Listen("tcp", ":4000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		panic(err)
	}
	defer listen.Close()

	s := grpc.NewServer()
	pb.RegisterTransformServer(s, &Transformer{})

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
		panic(err)
	}
}
