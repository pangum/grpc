package grpc

type metric struct {
	// 是否开启
	Enabled *bool `default:"true" json:"enabled" yaml:"enabled" xml:"enabled" toml:"enabled"`
	// 路径
	Pattern string `default:"/metrics" json:"pattern" yaml:"pattern" xml:"pattern" toml:"pattern"`
}
