package config

type KeepalivePolicy struct {
	// 无流许可
	Permit bool `default:"true" json:"permit" yaml:"permit" xml:"permit" toml:"permit"`
}
