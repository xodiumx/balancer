package tests

import (
	"balancer/src/core/config"
	"balancer/src/core/logger"
	"context"
	"go.uber.org/zap"
	"log"
	"os"
	"testing"

	"balancer/src/core/handler"
	"balancer/src/proto"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := logger.InitLogger(); err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer func(Log *zap.Logger) {
		err := Log.Sync()
		if err != nil {
			log.Fatalf("failed to sync logger: %v", err)
		}
	}(logger.Log)

	os.Exit(m.Run())
}

type handlerFixture struct {
	Handler           *handler.Handler
	firstOriginalURL  string
	secondOriginalURL string
	firstCDNURL       string
	secondCDNURL      string
}

func newHandlerFixture() *handlerFixture {
	cfg := config.Load()
	cfg.Frequency = 10
	return &handlerFixture{
		Handler:           handler.NewHandler(cfg),
		firstOriginalURL:  "http://s1.origin-cluster/video/123/file.m3u8",
		secondOriginalURL: "http://s2.origin-cluster/video/123/file.m3u8",
		firstCDNURL:       "http://cdn.default.com/s1/video/123/file.m3u8",
		secondCDNURL:      "http://cdn.default.com/s2/video/123/file.m3u8",
	}
}

func TestGetRedirect_Every10thRequestGoesToOrigin(t *testing.T) {

	h := newHandlerFixture()

	// send 9 "CDN"-request
	for i := 0; i < 9; i++ {
		resp, err := h.Handler.GetRedirect(context.Background(), &proto.VideoRequest{
			Video: h.firstOriginalURL,
		})
		assert.NoError(t, err)
		assert.Equal(t, h.firstCDNURL, resp.RedirectUrl)
	}

	// 10 must redirect on origin
	resp, err := h.Handler.GetRedirect(context.Background(), &proto.VideoRequest{
		Video: h.firstOriginalURL,
	})
	assert.NoError(t, err)
	assert.Equal(t, h.firstOriginalURL, resp.RedirectUrl)
}

func TestGetRedirect_SubDomainCached(t *testing.T) {

	h := newHandlerFixture()

	// send 9 "CDN"-request
	for i := 0; i < 9; i++ {
		resp, err := h.Handler.GetRedirect(context.Background(), &proto.VideoRequest{
			Video: h.firstOriginalURL,
		})
		assert.NoError(t, err)
		assert.Equal(t, h.firstCDNURL, resp.RedirectUrl)
	}

	// 10 A request with a different domain should not return the original url
	{
		resp, err := h.Handler.GetRedirect(context.Background(), &proto.VideoRequest{
			Video: h.secondOriginalURL,
		})
		assert.NoError(t, err)
		assert.Equal(t, h.secondCDNURL, resp.RedirectUrl)
	}

	// 11 The request with the first domain should return the original URL
	{
		resp, err := h.Handler.GetRedirect(context.Background(), &proto.VideoRequest{
			Video: h.firstOriginalURL,
		})
		assert.NoError(t, err)
		assert.Equal(t, h.firstOriginalURL, resp.RedirectUrl)
	}
}
