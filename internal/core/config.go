package core

import (
	"github.com/pangum/grpc/internal/config"
)

type Config struct {
	// 服务器端配置
	Server config.Server `json:"Server" yaml:"Server" xml:"Server" toml:"Server" validate:"required"`
	// 客户端配置
	Clients []config.Client `json:"clients" yaml:"clients" toml:"clients" xml:"clients"`
	// gRPC配置
	Options config.Options `json:"Options" yaml:"Options" xml:"Options" toml:"Options"`
}
