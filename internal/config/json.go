package config

type Json struct {
	// 是否允许多行
	Multiline bool `json:"multiline" yaml:"multiline" xml:"multiline" toml:"multiline"`
	// 前缀
	Indent string `json:"indent" yaml:"indent" xml:"indent" toml:"indent"`
	// 允许部分
	Partial bool `json:"partial" yaml:"partial" xml:"partial" toml:"partial"`
	// 选项列表
	// nolint: lll
	Options []string `default:"['enum_as_numbers', 'name_as_proto']" json:"options" yaml:"options" xml:"options" toml:"options"`
	// 是否允许不填充
	Unpopulated *bool `default:"true" json:"unpopulated" yaml:"unpopulated" xml:"unpopulated" toml:"unpopulated"`
	// 是否允许丢弃
	Discard *bool `default:"true" json:"discard" yaml:"discard" xml:"discard" toml:"discard"`
}
