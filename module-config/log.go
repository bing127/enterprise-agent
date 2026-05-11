package config

// LogConfig 日志配置。
type LogConfig struct {
	Level      string `yaml:"level"`  // debug | info | warn | error
	Format     string `yaml:"format"` // json | text
	Output     string `yaml:"output"` // stdout | file
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"` // MB
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"` // 天
}
