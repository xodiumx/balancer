package utils

import (
	"balancer/src/core/logger"
	pb "balancer/src/proto"
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
)

func ParseURL(req *pb.VideoRequest) (string, string, string, error) {

	// Parse url
	originalURL := req.GetVideo()
	if originalURL == "" {
		msg := "blank url in request"
		logger.Log.Warn(msg)
		return "", "", "", fmt.Errorf(msg)
	}

	parsed, err := url.Parse(originalURL)
	if err != nil {
		msg := "invalid video URL: %w"
		logger.Log.Warn(msg, zap.Error(err))
		return "", "", "", fmt.Errorf(msg, err)
	}

	parts := strings.Split(parsed.Host, ".")
	if len(parts) <= 1 {
		msg := "invalid origin host: %s"
		logger.Log.Warn(msg)
		return "", "", "", fmt.Errorf(msg, parts)
	}

	subDomain := parts[0] // e.g., s1 // TODO validate subDomains
	path := parsed.Path   // /video/123/xcg2djHckad.m3u8

	return originalURL, subDomain, path, nil
}

// CounterMap - Cache diffs subdomain
type CounterMap struct {
	counters sync.Map
}

func (c *CounterMap) IncrementAndGet(key string) uint64 {
	val, _ := c.counters.LoadOrStore(key, new(uint64))
	return atomic.AddUint64(val.(*uint64), 1)
}
