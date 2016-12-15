package config

import (
	"time"
)

type Config struct {
	Period time.Duration `config:"period"`
	URLs   []string      `config:"urls" validate:"nonzero,required"`
}

var DefaultConfig = Config{
	Period: 1 * time.Second,
	URLs:   []string{"http://127.0.0.1/status"},
}
