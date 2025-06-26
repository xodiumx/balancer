package main

import (
	"balancer/src/core/logger"
	pb "balancer/src/proto"
	"context"
	"github.com/goccy/go-json"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net/http"
	"time"
)

// TODO: refactor
func main() {

	// TODO: config

	// Init logger
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("Init logger error: %v", err)
	}
	defer func(Log *zap.Logger) {
		err := Log.Sync()
		if err != nil {

		}
	}(logger.Log)

	// Ctx for client
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Client init
	customMux := http.NewServeMux()
	customMux.HandleFunc("/watch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logger.Log.Error("failed to close body: %v", zap.Error(err))
			}
		}(r.Body)

		// Parse JSON
		var req pb.VideoRequest
		if err := json.Unmarshal(body, &req); err != nil {
			msg := "bad json format"
			logger.Log.Error(msg, zap.Error(err))
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// gRPC Client
		conn, err := grpc.DialContext(
			ctx,
			"balancer:50051",
			grpc.WithTransportCredentials(insecure.NewCredentials()), // grpc.WithTransportCredentials secure
		) // TODO: config for target url
		if err != nil {
			http.Error(w, "grpc dial error", http.StatusInternalServerError)
			return
		}
		defer func(conn *grpc.ClientConn) {
			err := conn.Close()
			if err != nil {
				log.Printf("failed to close connection: %v", err)
			}
		}(conn)

		client := pb.NewVideoBalancerClient(conn)
		ctx, cancel := context.WithTimeout(r.Context(), time.Second)
		defer cancel()

		resp, err := client.GetRedirect(ctx, &req)
		if err != nil {
			http.Error(w, "grpc error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Get Redirect
		http.Redirect(w, r, resp.RedirectUrl, http.StatusFound)
	})

	log.Println("🌐 Gateway started on :8080")
	err := http.ListenAndServe(":8080", customMux) // TODO: config for addr, loggs middleware
	if err != nil {
		return
	}

}
