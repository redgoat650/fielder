package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/redgoat650/fielder/server/grpc"
)

var (
	port = flag.Uint("port", 8080, "Port to launch the server on")
)

type testServer struct {
}

func (s *testServer) TestRPC(ctx context.Context, in *testgrpc.Ping) (*testgrpc.Pong, error) {
	fmt.Println(in)
	return &testgrpc.Pong{Message: "Fuck you"}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	testgrpc.RegisterTestServiceServer(grpcServer, &testServer{})
	fmt.Println("http://localhost:8080")
	grpcServer.Serve(lis)
}
