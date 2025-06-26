package main

import (
	"balancer/src/core/config"
	"balancer/src/core/logger"
	"go.uber.org/zap"
	"log"
	"net"
	//"runtime"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"

	"balancer/src/core/handler"
	pb "balancer/src/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// Init cfg
	cfg := config.Load()

	// Init logger
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("Init logger error: %v", err)
	}
	defer func(Log *zap.Logger) {
		err := Log.Sync()
		if err != nil {

		}
	}(logger.Log)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Log.Fatal("failed to listen", zap.Error(err))
	}

	// Init Server
	zapLogger := logger.Log.WithOptions(zap.AddCaller())

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(
				grpczap.UnaryServerInterceptor(zapLogger),
			),
		),
	)

	// Register handler
	pb.RegisterVideoBalancerServer(s, handler.NewHandler(cfg))

	// Use local requests
	if cfg.DEBUG {
		reflection.Register(s)
	}

	// gRPC Gateway HTTP server
	//go func() {
	//	ctx := context.Background()
	//	mux := runtime.NewServeMux()
	//
	//	opts := []grpc.DialOption{grpc.WithInsecure()} // WithTransportCredentials prod
	//
	//	err := pb.RegisterVideoBalancerHandlerFromEndpoint(ctx, mux, cfg.ServerBind, opts)
	//	if err != nil {
	//		logger.Log.Fatal("failed to start grpc-gateway", zap.Error(err))
	//	}
	//
	//	logger.Log.Info("🌐 gRPC-Gateway started", zap.String("addr", ":8080"))
	//	if err := http.ListenAndServe(":8080", mux); err != nil {
	//		logger.Log.Fatal("http gateway error", zap.Error(err))
	//	}
	//}()

	logger.Log.Info("🚀 gRPC server started", zap.String("addr", ":50051"))

	if err := s.Serve(lis); err != nil {
		logger.Log.Fatal("failed to serve", zap.Error(err))
	}
}
