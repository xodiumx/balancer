package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	DEBUG      bool
	CDNHost    string
	ServerBind string
	Frequency  uint64
}

func Load() *Config {

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.SetDefault("cdn_host", "cdn.example.com")

	cfg := &Config{
		DEBUG:      viper.GetBool("debug"),
		ServerBind: viper.GetString("server_bind"),
		CDNHost:    viper.GetString("cdn_host"),
		Frequency:  viper.GetUint64("frequency"),
	}

	log.Printf("[config] CDN_HOST=%s DEBUG=%v FREQUENCY=%v", cfg.CDNHost, cfg.DEBUG, cfg.Frequency)
	return cfg
}
