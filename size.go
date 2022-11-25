package grpc

type size struct {
	// 消息
	Msg msg `json:"msg" yaml:"msg" xml:"msg" toml:"msg" validate:"required"`
	// 窗口
	Window window `json:"window" yaml:"window" xml:"window" toml:"window" validate:"required"`
}
