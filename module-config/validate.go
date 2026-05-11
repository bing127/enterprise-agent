package config

import (
	"fmt"
	"strings"
)

// Validate checks required config fields and environment-specific constraints.
func (c *Config) Validate() error {
	env := strings.TrimSpace(c.App.Env)
	if env == "" {
		env = "dev"
	}
	if env != "dev" && env != "prod" {
		return fmt.Errorf("app.env must be dev or prod")
	}

	if c.Server.Port <= 0 {
		return fmt.Errorf("server.port must be > 0")
	}
	if strings.TrimSpace(c.Auth.SecretKey) == "" {
		return fmt.Errorf("auth.secret_key is required")
	}

	output := strings.TrimSpace(c.Log.Output)
	if output == "" {
		output = "stdout"
	}
	if output != "stdout" && output != "file" && output != "both" {
		return fmt.Errorf("log.output must be stdout, file or both")
	}
	if (output == "file" || output == "both") && strings.TrimSpace(c.Log.FilePath) == "" {
		return fmt.Errorf("log.file_path is required when log.output is file/both")
	}

	if env == "prod" {
		if strings.TrimSpace(c.Redis.Addr) == "" {
			return fmt.Errorf("redis.addr is required in prod")
		}
		if strings.TrimSpace(c.Weaviate.Host) == "" {
			return fmt.Errorf("weaviate.host is required in prod")
		}
		if strings.TrimSpace(c.AI.Provider) == "" {
			return fmt.Errorf("ai.provider is required in prod")
		}
		switch strings.TrimSpace(c.AI.Provider) {
		case "openai":
			if strings.TrimSpace(c.AI.OpenAI.APIKey) == "" {
				return fmt.Errorf("ai.openai.api_key is required when ai.provider=openai")
			}
			if strings.TrimSpace(c.AI.OpenAI.Model) == "" {
				return fmt.Errorf("ai.openai.model is required when ai.provider=openai")
			}
		case "ark":
			if strings.TrimSpace(c.AI.Ark.APIKey) == "" {
				return fmt.Errorf("ai.ark.api_key is required when ai.provider=ark")
			}
			if strings.TrimSpace(c.AI.Ark.EndpointID) == "" {
				return fmt.Errorf("ai.ark.endpoint_id is required when ai.provider=ark")
			}
		case "qwen":
			if strings.TrimSpace(c.AI.Qwen.APIKey) == "" {
				return fmt.Errorf("ai.qwen.api_key is required when ai.provider=qwen")
			}
			if strings.TrimSpace(c.AI.Qwen.Model) == "" {
				return fmt.Errorf("ai.qwen.model is required when ai.provider=qwen")
			}
		default:
			return fmt.Errorf("ai.provider must be one of openai/ark/qwen")
		}
	}

	return nil
}
