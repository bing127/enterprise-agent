package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// Config 汇总所有子配置，各子配置类型定义在独立文件中。
type Config struct {
	App           AppConfig           `yaml:"app" mapstructure:"app"`
	Server        ServerConfig        `yaml:"server" mapstructure:"server"`
	Auth          AuthConfig          `yaml:"auth" mapstructure:"auth"`
	Elasticsearch ElasticsearchConfig `yaml:"elasticsearch" mapstructure:"elasticsearch"`
	Weaviate      WeaviateConfig      `yaml:"weaviate" mapstructure:"weaviate"`
	Redis         RedisConfig         `yaml:"redis" mapstructure:"redis"`
	AI            AIConfig            `yaml:"ai" mapstructure:"ai"`
	Log           LogConfig           `yaml:"log" mapstructure:"log"`
}

// Load 根据 APP_ENV 加载对应环境配置文件（仅支持 dev / prod）。
// 加载文件：config.<env>.yaml。
func Load(resourceDir string) (*Config, error) {
	env := strings.TrimSpace(os.Getenv("APP_ENV"))
	if env == "" {
		env = "dev"
	}
	if env != "dev" && env != "prod" {
		return nil, fmt.Errorf("invalid APP_ENV=%q, only dev/prod are supported", env)
	}

	envFile := filepath.Join(resourceDir, fmt.Sprintf("config.%s.yaml", env))
	if _, statErr := os.Stat(envFile); statErr != nil {
		return nil, fmt.Errorf("env config not found (%s): %w", envFile, statErr)
	}

	raw, err := os.ReadFile(envFile)
	if err != nil {
		return nil, fmt.Errorf("read env config (%s): %w", env, err)
	}
	// 兼容配置中的 ${VAR} 占位符。
	expanded := os.ExpandEnv(string(raw))

	v := viper.New()
	v.SetConfigType("yaml")
	if err := v.ReadConfig(bytes.NewBufferString(expanded)); err != nil {
		return nil, fmt.Errorf("parse env config (%s): %w", env, err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal env config (%s): %w", env, err)
	}
	cfg.App.Env = env
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate env config (%s): %w", env, err)
	}

	return &cfg, nil
}
