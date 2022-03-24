package grpc

type config struct {
	// 服务器端配置
	Server serverConfig `json:"server" yaml:"server" xml:"server" toml:"server" validate:"required"`
	// 客户端配置
	Clients []clientConfig `json:"clients" yaml:"clients" toml:"clients" xml:"clients"`
	// gRPC配置
	Options options `json:"options" yaml:"options" xml:"options" toml:"options"`
}
