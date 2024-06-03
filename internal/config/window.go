package config

import (
	"github.com/goexl/gox"
)

type Window struct {
	// 初始
	// 1GB
	Initial gox.Bytes `default:"1GB" json:"initial" yaml:"initial" xml:"initial" toml:"initial"`
	// 连接
	// 1GB
	Connection gox.Bytes `default:"1GB" json:"connection" yaml:"connection" xml:"connection" toml:"connection"`
}
