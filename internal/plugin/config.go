package plugin

import (
	"github.com/pangum/grpc/internal/config"
)

type Config struct {
	// 服务器端配置
	Server *config.Server `json:"server" yaml:"server" xml:"server" toml:"server"`
	// 网关配置
	Gateway *config.Gateway `json:"gateway" yaml:"gateway" xml:"gateway" toml:"gateway"`
	// 客户端配置
	Clients []config.Client `json:"clients" yaml:"clients" toml:"clients" xml:"clients"`
	// gRPC配置
	Options config.Options `json:"options" yaml:"options" xml:"options" toml:"options"`
}
