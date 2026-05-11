package config

// WeaviateConfig Weaviate 向量数据库连接配置。
type WeaviateConfig struct {
	Scheme   string `yaml:"scheme"`
	Host     string `yaml:"host"`
	GrpcPort int    `yaml:"grpc_port"`
	APIKey   string `yaml:"api_key"`
}
