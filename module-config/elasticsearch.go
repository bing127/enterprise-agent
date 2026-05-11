package config

// ElasticsearchConfig Elasticsearch 连接配置。
type ElasticsearchConfig struct {
	Addresses   []string `yaml:"addresses" mapstructure:"addresses"`
	Username    string   `yaml:"username" mapstructure:"username"`
	Password    string   `yaml:"password" mapstructure:"password"`
	IndexPrefix string   `yaml:"index_prefix" mapstructure:"index_prefix"`
}
