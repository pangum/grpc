package grpc

import `fmt`

type serverConfig struct {
	// 绑定监听主机
	Host string `json:"host" yaml:"host"`
	// 绑定监听端口
	Port int `default:"9001" json:"port" yaml:"port" validate:"required,min=1,max=65535"`
}

func (sc *serverConfig) Addr() string {
	return fmt.Sprintf("%s:%d", sc.Host, sc.Port)
}
