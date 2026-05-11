package config

import "time"

// ServerConfig HTTP 服务器（Hertz）配置。
type ServerConfig struct {
	Host               string        `yaml:"host"`
	Port               int           `yaml:"port"`
	ReadTimeout        time.Duration `yaml:"read_timeout"`
	WriteTimeout       time.Duration `yaml:"write_timeout"`
	MaxRequestBodySize int           `yaml:"max_request_body_size"`
}
