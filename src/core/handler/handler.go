package handler

import (
	"balancer/src/core/config"
	"balancer/src/logger"
	"context"
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"sync/atomic"

	pb "balancer/src/proto"
)

type Handler struct {
	counter uint64
	cfg     *config.Config
	pb.UnimplementedVideoBalancerServer
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg}
}

func (h *Handler) GetRedirect(_ context.Context, req *pb.VideoRequest) (*pb.VideoResponse, error) {
	count := atomic.AddUint64(&h.counter, 1)
	originalURL := req.GetVideo()

	if count%h.cfg.Frequency == 0 {
		logger.Log.Warn("Request data", // Warn for colored log
			zap.String("Redirect url", originalURL),
			zap.Uint64("request_number", count),
		)
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

	// if we need can if else for https
	format := "http://%s/%s%s"
	cdnURL := fmt.Sprintf(format, h.cfg.CDNHost, originServer, path)

	logger.Log.Info("Request data",
		zap.String("original URL", originalURL),
		zap.String("redirect URL", cdnURL),
		zap.Uint64("request_number", count),
	)

	return &pb.VideoResponse{RedirectUrl: cdnURL}, nil
}
