package config

import "time"

type Config struct {
	Redis  RedisConfig
	Server ServerConfig
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type ServerConfig struct {
	Port                      int
	HttpServerShutdownTimeout time.Duration
}
