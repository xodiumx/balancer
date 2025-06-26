package tests

import (
	"testing"

	"balancer/src/core/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Default(t *testing.T) {
	cfg := config.Load()
	assert.Equal(t, "cdn.example.com", cfg.CDNHost)
	assert.Equal(t, uint64(0), cfg.Frequency)
	assert.Equal(t, "", cfg.ServerBind)
	assert.Equal(t, false, cfg.DEBUG)
}
