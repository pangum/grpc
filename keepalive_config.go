package grpc

import (
	`time`
)

type keepaliveConfig struct {
	// 保持时长
	Time time.Duration `default:"10s" json:"time" yaml:"time" xml:"time" toml:"time"`
	// 超时
	Timeout time.Duration `default:"3s" json:"timeout" yaml:"timeout" xml:"timeout" toml:"timeout"`
	// 空闲时长
	Idle time.Duration `default:"3s" json:"idle" yaml:"idle" xml:"idle" toml:"idle"`
	// 策略
	Policy keepalivePolicyConfig `json:"policy" yaml:"policy" xml:"policy" toml:"policy"`
}
