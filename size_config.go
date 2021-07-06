package grpc

type sizeConfig struct {
	// 消息
	Msg msgConfig `json:"msg" validate:"required"`
	// 窗口
	Window windowConfig `json:"window" validate:"required"`
}
