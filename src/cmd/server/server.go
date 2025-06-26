package main

import (
	"balancer/src/core/config"
	"balancer/src/core/handler"
	"balancer/src/core/logger"
	pb "balancer/src/proto"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

func main() {

	// Init cfg
	cfg := config.Load()

	//Init logger (RPS decrease 3-4 times)
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

		// Profiler (RPS decrease 30 %)
		go func() {
			logger.Log.Info("📈 Starting pprof on :6060")
			if err := http.ListenAndServe(":6060", nil); err != nil {
				log.Fatalf("pprof server failed: %v", err)
			}
		}()
	}

	logger.Log.Info("🚀 gRPC server started", zap.String("addr", ":50051"))

	if err := s.Serve(lis); err != nil {
		logger.Log.Fatal("failed to serve", zap.Error(err))
	}
}
