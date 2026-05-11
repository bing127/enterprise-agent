package config

// AppConfig 应用基础配置。
type AppConfig struct {
	Name  string `yaml:"name"`
	Env   string `yaml:"env"`
	Port  int    `yaml:"port"`
	Debug bool   `yaml:"debug"`
}
