package grpc

type msgConfig struct {
	// 发送大小
	// 4GB
	Send int `default:"4294967296" json:"send"`
	// 接收大小
	// 4GB
	Receive int `default:"4294967296" json:"receive"`
}
