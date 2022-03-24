package grpc

type clientConfig struct {
	// 名称
	Name string `json:"name" yaml:"name" xml:"name" toml:"name" validate:"required_without=Names"`
	// 名称列表
	Names []string `json:"names" yaml:"names" xml:"names" toml:"names" validate:"required_without=Name"`
	// 连接地址
	Addr string `json:"addr" yaml:"addr" xml:"addr" toml:"addr" validate:"required,url"`
}
