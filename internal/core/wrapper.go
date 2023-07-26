package core

type Wrapper struct {
	// gRPC配置
	Grpc Config `json:"grpc" yaml:"grpc" xml:"grpc" toml:"grpc" validate:"required"`
}
