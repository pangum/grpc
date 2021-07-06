package grpc

type clientConfig struct {
	// 名称
	Name string `json:"name" yaml:"name" xml:"name" validate:"required"`
	// 连接地址
	Addr string `json:"addr" yaml:"addr" xml:"addr" validate:"required,url"`
}
