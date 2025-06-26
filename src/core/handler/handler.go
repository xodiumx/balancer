package handler

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync/atomic"

	pb "balancer/src/proto"
)

type Handler struct {
	CDNHost   string
	frequency uint64
	counter   uint64
	pb.UnimplementedVideoBalancerServer
}

func NewHandler(cdnHost string) *Handler {
	return &Handler{CDNHost: cdnHost, frequency: 10}
}

func (h *Handler) GetRedirect(_ context.Context, req *pb.VideoRequest) (*pb.VideoResponse, error) {
	count := atomic.AddUint64(&h.counter, 1)
	originalURL := req.GetVideo()

	if count%h.frequency == 0 {
		return &pb.VideoResponse{RedirectUrl: originalURL}, nil
	}

	parsed, err := url.Parse(originalURL)
	if err != nil {
		return nil, fmt.Errorf("invalid video URL: %w", err)
	}

	parts := strings.Split(parsed.Host, ".")
	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid origin host")
	}

	originServer := parts[0] // e.g., s1
	path := parsed.Path      // /video/123/xcg2djHckad.m3u8

	cdnURL := fmt.Sprintf("http://%s/%s%s", h.CDNHost, originServer, path)
	return &pb.VideoResponse{RedirectUrl: cdnURL}, nil
}
