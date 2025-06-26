package handler

import (
	"balancer/src/core/config"
	"balancer/src/core/logger"
	"balancer/src/core/utils"
	"context"
	"fmt"
	"go.uber.org/zap"

	pb "balancer/src/proto"
)

type Handler struct {
	counterMap *utils.CounterMap
	cfg        *config.Config
	pb.UnimplementedVideoBalancerServer
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{cfg: cfg, counterMap: &utils.CounterMap{}}
}

func (h *Handler) GetRedirect(_ context.Context, req *pb.VideoRequest) (*pb.VideoResponse, error) {

	// Parse video url
	originalURL, subDomain, path, err := utils.ParseURL(req)
	if err != nil {
		return nil, err
	}

	// Counter
	count := h.counterMap.IncrementAndGet(subDomain)
	if count%h.cfg.Frequency == 0 {
		logger.Log.Warn("Request data", // Warn for colored log
			zap.String("Redirect url", originalURL),
			zap.Uint64("request_number", count),
		)
		return &pb.VideoResponse{RedirectUrl: originalURL}, nil
	}

	// if we need can if else for https
	format := "http://%s/%s%s"
	cdnURL := fmt.Sprintf(format, h.cfg.CDNHost, subDomain, path)

	logger.Log.Info("Request data",
		zap.String("original URL", originalURL),
		zap.String("redirect URL", cdnURL),
		zap.Uint64("request_number", count),
	)

	return &pb.VideoResponse{RedirectUrl: cdnURL}, nil
}
