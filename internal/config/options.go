package config

type Options struct {
	// 大小配置
	Size Size `json:"size" yaml:"size" xml:"size" toml:"size"`
	// 长连接
	Keepalive Keepalive `json:"keepalive" yaml:"keepalive" xml:"keepalive" toml:"keepalive"`
}
