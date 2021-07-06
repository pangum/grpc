package grpc

type panguConfig struct {
	// gRPC配置
	Grpc config `json:"grpc" yaml:"grpc" validate:"required"`
}
