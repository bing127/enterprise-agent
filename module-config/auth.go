package config

// AuthConfig JWT 鉴权配置。
type AuthConfig struct {
	SecretKey  string `yaml:"secret_key"`
	TokenTTL   int    `yaml:"token_ttl"`   // 单位：秒
	RefreshTTL int    `yaml:"refresh_ttl"` // 单位：秒
}
