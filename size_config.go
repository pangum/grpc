package grpc

type sizeConfig struct {
	// 消息
	Msg msgConfig `json:"msg" yaml:"msg" xml:"msg" toml:"msg" validate:"required"`
	// 窗口
	Window windowConfig `json:"window" yaml:"window" xml:"window" toml:"window" validate:"required"`
}
