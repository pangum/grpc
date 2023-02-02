package grpc

type gateway struct {
	// 是否开启
	Enabled *bool `default:"true" json:"enabled" yaml:"enabled" xml:"enabled" toml:"enabled"`
}
