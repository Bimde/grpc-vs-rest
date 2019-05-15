package main

import (
	"log"
	"net"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	
	"github.com/Bimde/grpc-vs-rest/pb"
)

type server struct{}

func main() {
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRandomServiceServer(s, &server{})
	log.Println("Starting gRPC server")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) DoSomething(_ context.Context, random *pb.Random) (*pb.Random, error) {
	random.RandomString = "[Updated] " + random.RandomString;
	return random, nil
}