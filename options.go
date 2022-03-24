package grpc

type options struct {
	// 大小配置
	Size sizeConfig `json:"size" yaml:"size" xml:"size" toml:"size"`
	// 长连接
	Keepalive keepaliveConfig `json:"keepalive" yaml:"keepalive" xml:"keepalive" toml:"keepalive"`
}
