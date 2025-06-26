package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DEBUG      bool
	CDNHost    string
	ServerBind string
	Frequency  uint64
}

func Load() *Config {

	viper.AutomaticEnv()
	viper.SetDefault("cdn_host", "cdn.default.com")

	cfg := &Config{
		DEBUG:      viper.GetBool("debug"),
		ServerBind: viper.GetString("server_bind"),
		CDNHost:    viper.GetString("cdn_host"),
		Frequency:  viper.GetUint64("frequency"),
	}

	log.Printf(
		"[config] DEBUG=%v SERVER_BINF=%s CDN_HOST=%s FREQUENCY=%v",
		cfg.DEBUG, cfg.ServerBind, cfg.CDNHost, cfg.Frequency,
	)
	return cfg
}
