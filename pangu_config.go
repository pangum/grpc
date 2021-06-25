package grpc

type panguConfig struct {
	// GRPC配置
	GRPC struct {
		// 服务器端配置
		Server config `json:"server" yaml:"server" validate:"required"`
	} `json:"grpc" yaml:"grpc" validate:"required"`
}
