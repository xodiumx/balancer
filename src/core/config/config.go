package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	CDNHost string
}

func Load() *Config {
	viper.SetEnvPrefix("app")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("cdn_host", "cdn.example.com")

	cfg := &Config{
		CDNHost: viper.GetString("cdn_host"),
	}

	log.Printf("[config] CDN_HOST=%s", cfg.CDNHost) // TODO : upd to zap
	return cfg
}
