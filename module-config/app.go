package config

// AppConfig 应用基础配置。
type AppConfig struct {
	Name  string `yaml:"name" mapstructure:"name"`
	Env   string `yaml:"env" mapstructure:"env"`
	Port  int    `yaml:"port" mapstructure:"port"`
	Debug bool   `yaml:"debug" mapstructure:"debug"`
}
