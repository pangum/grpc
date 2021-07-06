package grpc

type windowConfig struct {
	// 初始
	// 1GB
	Initial int32 `default:"1073741824" json:"initial"`
	// 连接
	// 1GB
	Connection int32 `default:"1073741824" json:"connection"`
}
