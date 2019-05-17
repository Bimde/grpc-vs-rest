package main

import (
	context "context"
	"log"

	"github.com/Bimde/grpc-vs-rest/pb"
	"google.golang.org/grpc"
)

func random(c context.Context, input *pb.Random) (*pb.Random, error) {
	conn, err := grpc.Dial("sever_address:port")
	if err != nil {
		log.Fatalf("Dial failed: %v", err)
	}

	client := pb.NewRandomServiceClient(conn)
	return client.DoSomething(c, &pb.Random{})
}
