package grpc

type config struct {
	// 服务器端配置
	Server serverConfig `json:"server" yaml:"server" xml:"server" validate:"required"`
	// 客户端配置
	Clients []clientConfig `json:"clients" yaml:"clients" xml:"clients"`
	// gRPC配置
	Options options `json:"options"`
}
