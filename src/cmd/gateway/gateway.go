package main

import (
	pb "balancer/src/proto"
	"context"
	"encoding/json"
	"google.golang.org/grpc"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {

	// TODO: logs, config

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
				log.Printf("failed to close body: %v", err)
			}
		}(r.Body)

		// Parse JSON
		var req pb.VideoRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		// gRPC Client
		conn, err := grpc.Dial("balancer:50051", grpc.WithInsecure()) // TODO: config for target url
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
	err := http.ListenAndServe(":8080", customMux) // TODO: config for addr
	if err != nil {
		return
	}

}
