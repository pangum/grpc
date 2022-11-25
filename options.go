package grpc

type options struct {
	// 大小配置
	Size size `json:"size" yaml:"size" xml:"size" toml:"size"`
	// 长连接
	Keepalive _keepalive `json:"_keepalive" yaml:"_keepalive" xml:"_keepalive" toml:"_keepalive"`
}
