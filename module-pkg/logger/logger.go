package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/lumberjack.v2"
)

var (
	global *zap.Logger
	once   sync.Once
)

// Config 日志配置
type Config struct {
	Level      string // "debug" | "info" | "warn" | "error"
	Env        string // "dev" | "prod"
	Output     string // "stdout" | "file" | "both"
	FilePath   string // 日志文件路径
	MaxSizeMB  int    // 单文件最大 MB
	MaxBackups int    // 保留备份数
	MaxAgeDays int    // 最大保留天数（天）
}

// Init 初始化全局 logger（仅 stdout）。建议在 main() 最开始调用一次。
// level: "debug" | "info" | "warn" | "error"
// env:   "dev" | "prod"
func Init(level, env string) {
	once.Do(func() {
		global = newLogger(Config{Level: level, Env: env, Output: "stdout"})
	})
}

// InitWithConfig 初始化全局 logger，支持文件输出。
func InitWithConfig(cfg Config) {
	once.Do(func() {
		global = newLogger(cfg)
	})
}

// L 返回全局 logger，未初始化时返回默认 info 级别 logger。
func L() *zap.Logger {
	if global == nil {
		global = newLogger(Config{Level: "info", Env: "prod", Output: "stdout"})
	}
	return global
}

// Sync 刷新缓冲日志，程序退出前调用。
func Sync() {
	if global != nil {
		_ = global.Sync()
	}
}

func newLogger(cfg Config) *zap.Logger {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(cfg.Level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "time"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	var encoder zapcore.Encoder
	if cfg.Env == "dev" {
		devCfg := zap.NewDevelopmentEncoderConfig()
		devCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		devCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(devCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	stdoutWriter := zapcore.AddSync(os.Stdout)
	output := strings.TrimSpace(cfg.Output)
	if output == "" {
		output = "stdout"
	}

	fileWriter, fileErr := buildFileWriter(cfg)

	var core zapcore.Core
	switch output {
	case "file":
		if fileErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "logger file output disabled: %v\n", fileErr)
			core = zapcore.NewCore(encoder, stdoutWriter, zapLevel)
			break
		}
		core = zapcore.NewCore(encoder, fileWriter, zapLevel)
	case "both":
		if fileErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "logger file output disabled: %v\n", fileErr)
			core = zapcore.NewCore(encoder, stdoutWriter, zapLevel)
			break
		}
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, stdoutWriter, zapLevel),
			zapcore.NewCore(encoder, fileWriter, zapLevel),
		)
	default: // "stdout"
		core = zapcore.NewCore(encoder, stdoutWriter, zapLevel)
	}

	opts := []zap.Option{zap.AddCaller(), zap.AddCallerSkip(0)}
	if cfg.Env == "dev" {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return zap.New(core, opts...)
}

func buildFileWriter(cfg Config) (zapcore.WriteSyncer, error) {
	path := strings.TrimSpace(cfg.FilePath)
	if path == "" {
		return nil, fmt.Errorf("empty log file path")
	}
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   path,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   true,
	}), nil
}
