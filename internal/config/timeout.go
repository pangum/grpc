package config

import (
	"time"
)

type Timeout struct {
	// 读
	Read time.Duration `default:"15s" json:"read" yaml:"read" xml:"read" toml:"read"`
	// 头
	Header time.Duration `default:"15s" json:"header" yaml:"header" xml:"header" toml:"header"`
}
