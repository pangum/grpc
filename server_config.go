package grpc

import "fmt"

type server struct {
	// 绑定监听主机
	Host string `json:"host" yaml:"host" xml:"host" toml:"host"`
	// 绑定监听端口
	Port int `default:"9001" json:"port" yaml:"port" xml:"port" toml:"port" validate:"required,min=1,max=65535"`
}

func (s *server) Addr() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
