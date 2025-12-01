package config

type JobConfig struct {
	Script struct {
		GlobalTimeout int `mapstructure:"global-timeout"`
		TimeOut       int `yaml:"timeout"`
	} `yaml:"script"`
}
