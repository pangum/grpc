package grpc

import (
	`time`
)

type keepaliveConfig struct {
	// 保持时长
	Time time.Duration `default:"10s" json:"time"`
	// 超时
	Timeout time.Duration `default:"3s" json:"timeout"`
	// 空闲时长
	Idle time.Duration `default:"3s" json:"idle"`
	// 策略
	Policy keepalivePolicyConfig `json:"policy"`
}
