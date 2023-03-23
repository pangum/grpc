package grpc

type wrapper struct {
	// gRPC配置
	Grpc config `json:"grpc" yaml:"grpc" xml:"grpc" toml:"grpc" validate:"required"`
}
