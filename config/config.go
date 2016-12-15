package config

import (
	"time"
	"net/url"
)

type Config struct {
	Period time.Duration `config:"period"`
	URLs   []*url.URL    `config:"urls"`
}

var defaultUrl, _ = url.Parse("http://127.0.0.1/status")

var DefaultConfig = Config{
	Period: 1 * time.Second,
	URLs: []*url.URL{defaultUrl},
}
