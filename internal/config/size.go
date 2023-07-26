package config

type Size struct {
	// 消息
	Msg Msg `json:"msg" yaml:"msg" xml:"msg" toml:"msg" validate:"required"`
	// 窗口
	Window Window `json:"window" yaml:"window" xml:"window" toml:"window" validate:"required"`
}
