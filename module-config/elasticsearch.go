package config

// ElasticsearchConfig Elasticsearch 连接配置。
type ElasticsearchConfig struct {
	Addresses   []string `yaml:"addresses"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	IndexPrefix string   `yaml:"index_prefix"`
}
