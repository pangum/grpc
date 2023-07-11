package grpc

type cors struct {
	// 是否开启
	Enabled bool `json:"enabled" yaml:"enabled" xml:"enabled" toml:"enabled"`
	// 允许跨域访问的源
	Allows []string `default:"['*']" json:"allows" yaml:"allows" xml:"allows" toml:"allows"`
	// 允许的请求方法
	// nolint:lll
	Methods []string `default:"['GET', 'POST', 'PUT', 'DELETE']" json:"methods" yaml:"methods" xml:"methods" toml:"methods"`
	// 允许的请求头
	// nolint:lll
	Headers []string `default:"['Content-Type', 'Authorization']" json:"headers" yaml:"headers" xml:"headers" toml:"headers"`
}
