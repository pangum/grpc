package config

import (
	"fmt"
)

type Server struct {
	// 绑定监听主机
	Host string `json:"host" yaml:"host" xml:"host" toml:"host"`
	// 绑定监听端口
	Port int `default:"9001" json:"port" yaml:"port" xml:"port" toml:"port" validate:"required,min=1,max=65535"`
	// 反射
	// 可以通过配置反射来开启服务器反射字段和方法的特性，方便客户端通过反射来调用方法
	Reflection *bool `default:"true" json:"reflection" yaml:"reflection" xml:"reflection" toml:"reflection"`
}

func (s *Server) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
