package config

// AIConfig AI 模型（Eino ChatModel）配置。
type AIConfig struct {
	Provider string       `yaml:"provider" mapstructure:"provider"` // openai | ark | qwen
	OpenAI   OpenAIConfig `yaml:"openai" mapstructure:"openai"`
	Ark      ArkConfig    `yaml:"ark" mapstructure:"ark"`
	Qwen     QwenConfig   `yaml:"qwen" mapstructure:"qwen"`
}

// OpenAIConfig OpenAI / 兼容接口配置。
type OpenAIConfig struct {
	APIKey      string  `yaml:"api_key" mapstructure:"api_key"`
	BaseURL     string  `yaml:"base_url" mapstructure:"base_url"`
	Model       string  `yaml:"model" mapstructure:"model"`
	Temperature float32 `yaml:"temperature" mapstructure:"temperature"`
	MaxTokens   int     `yaml:"max_tokens" mapstructure:"max_tokens"`
}

// ArkConfig 火山引擎 Ark 大模型配置。
type ArkConfig struct {
	APIKey     string `yaml:"api_key" mapstructure:"api_key"`
	EndpointID string `yaml:"endpoint_id" mapstructure:"endpoint_id"`
	Region     string `yaml:"region" mapstructure:"region"`
}

// QwenConfig 阿里云通义千问配置。
type QwenConfig struct {
	APIKey string `yaml:"api_key" mapstructure:"api_key"`
	Model  string `yaml:"model" mapstructure:"model"`
}
