package config

type JobConfig struct {
	Script `yaml:"script"`
}
type Script struct {
	GlobalTimeout int `mapstructure:"global-timeout"`
	MaxConcurrent int `mapstructure:"max-concurrent"`
	TimeOut       int `yaml:"timeout"`
}
