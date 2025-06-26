package tests

import (
	"os"
	"testing"

	"balancer/src/core/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Default(t *testing.T) {
	err := os.Unsetenv("CDN_HOST")
	if err != nil {
		return
	}

	cfg := config.Load()
	assert.Equal(t, "cdn.example.com", cfg.CDNHost)
}
