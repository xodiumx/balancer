package tests

import (
	"balancer/src/core/config"
	"balancer/src/core/logger"
	"context"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"

	"balancer/src/core/handler"
	"balancer/src/proto"
	"github.com/stretchr/testify/assert"
)

func TestGetRedirect_Every10thRequestGoesToOrigin(t *testing.T) {

	err := logger.InitLogger()
	require.NoError(t, err)
	defer func(Log *zap.Logger) {
		err := Log.Sync()
		if err != nil {
		}
	}(logger.Log)

	cfg := config.Load()
	cfg.Frequency = 10

	h := handler.NewHandler(cfg)

	originalURL := "http://s1.origin-cluster/video/123/file.m3u8"
	cdnURL := "http://cdn.example.com/s1/video/123/file.m3u8"

	// отправим 9 "CDN"-запросов
	for i := 0; i < 9; i++ {
		resp, err := h.GetRedirect(context.Background(), &proto.VideoRequest{
			Video: originalURL,
		})
		assert.NoError(t, err)
		assert.Equal(t, cdnURL, resp.RedirectUrl)
	}

	// 10-й должен идти на origin
	resp, err := h.GetRedirect(context.Background(), &proto.VideoRequest{
		Video: originalURL,
	})
	assert.NoError(t, err)
	assert.Equal(t, originalURL, resp.RedirectUrl)
}
