package config

// WeaviateConfig Weaviate 向量数据库连接配置。
type WeaviateConfig struct {
	Scheme   string `yaml:"scheme" mapstructure:"scheme"`
	Host     string `yaml:"host" mapstructure:"host"`
	GrpcPort int    `yaml:"grpc_port" mapstructure:"grpc_port"`
	APIKey   string `yaml:"api_key" mapstructure:"api_key"`
}
