package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 汇总所有子配置，各子配置类型定义在独立文件中。
type Config struct {
	App           AppConfig           `yaml:"app"`
	Server        ServerConfig        `yaml:"server"`
	Auth          AuthConfig          `yaml:"auth"`
	Elasticsearch ElasticsearchConfig `yaml:"elasticsearch"`
	Weaviate      WeaviateConfig      `yaml:"weaviate"`
	Redis         RedisConfig         `yaml:"redis"`
	AI            AIConfig            `yaml:"ai"`
	Log           LogConfig           `yaml:"log"`
}

// Load 根据 APP_ENV 环境变量加载 YAML 配置文件。
// 加载顺序：config.yaml → config.<env>.yaml（覆盖）。
func Load(resourceDir string) (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}

	base, err := loadFile(filepath.Join(resourceDir, "config.yaml"))
	if err != nil {
		return nil, fmt.Errorf("load base config: %w", err)
	}

	envFile := filepath.Join(resourceDir, fmt.Sprintf("config.%s.yaml", env))
	if _, statErr := os.Stat(envFile); statErr == nil {
		override, err := loadFile(envFile)
		if err != nil {
			return nil, fmt.Errorf("load env config (%s): %w", env, err)
		}
		mergeConfig(base, override)
	}

	return base, nil
}

func loadFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// 支持 ${VAR} 环境变量展开
	data = []byte(os.ExpandEnv(string(data)))

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// mergeConfig 将 override 中非零值字段覆盖到 base。
func mergeConfig(base, override *Config) {
	if override.App.Env != "" {
		base.App = override.App
	}
	if override.Server.Port != 0 {
		base.Server = override.Server
	}
	if override.Auth.SecretKey != "" {
		base.Auth = override.Auth
	}
	if len(override.Elasticsearch.Addresses) > 0 {
		base.Elasticsearch = override.Elasticsearch
	}
	if override.Weaviate.Host != "" {
		base.Weaviate = override.Weaviate
	}
	if override.Redis.Addr != "" {
		base.Redis = override.Redis
	}
	if override.AI.Provider != "" {
		base.AI = override.AI
	}
	if override.Log.Level != "" {
		base.Log = override.Log
	}
}
