package config

import "time"

// ServerConfig HTTP 服务器（Hertz）配置。
type ServerConfig struct {
	Host               string        `yaml:"host" mapstructure:"host"`
	Port               int           `yaml:"port" mapstructure:"port"`
	ReadTimeout        time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout       time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
	MaxRequestBodySize int           `yaml:"max_request_body_size" mapstructure:"max_request_body_size"`
}
