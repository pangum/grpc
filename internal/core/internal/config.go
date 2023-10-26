package internal

import (
	"github.com/pangum/grpc/internal/config"
)

type Config struct {
	Server  *config.Server
	Gateway *config.Gateway
}

func NewConfig(server *config.Server, gateway *config.Gateway) *Config {
	return &Config{
		Server:  server,
		Gateway: gateway,
	}
}
