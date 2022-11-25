package grpc

import (
	"github.com/goexl/gox"
)

type window struct {
	// 初始
	// 1GB
	Initial gox.Size `default:"1GB" json:"initial" yaml:"initial" xml:"initial" toml:"initial"`
	// 连接
	// 1GB
	Connection gox.Size `default:"1GB" json:"connection" yaml:"connection" xml:"connection" toml:"connection"`
}
