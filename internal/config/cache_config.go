package config

import "time"

type CacheConfig struct {
	DefaultExpirationTime time.Duration `env:"CACHE_DEFAULT_EXPIRATION_TIME" envDefault:"24h"`
	CleanupInterval       time.Duration `env:"CACHE_CLEANUP_INTERVAL" envDefault:"24h"`
}
