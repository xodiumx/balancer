package main

import (
	"log"
	"net"

	"balancer/src/core/config"
	"balancer/src/core/handler"
	pb "balancer/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterVideoBalancerServer(s, handler.NewHandler(cfg.CDNHost))
	reflection.Register(s) // TODO if debug

	log.Println("Balancer service started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
