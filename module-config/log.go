package config

// LogConfig 日志配置。
type LogConfig struct {
	Level      string `yaml:"level" mapstructure:"level"`   // debug | info | warn | error
	Format     string `yaml:"format" mapstructure:"format"` // json | text
	Output     string `yaml:"output" mapstructure:"output"` // stdout | file
	FilePath   string `yaml:"file_path" mapstructure:"file_path"`
	MaxSize    int    `yaml:"max_size" mapstructure:"max_size"` // MB
	MaxBackups int    `yaml:"max_backups" mapstructure:"max_backups"`
	MaxAge     int    `yaml:"max_age" mapstructure:"max_age"` // 天
}
