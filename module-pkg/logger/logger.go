package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	global *zap.Logger
	once   sync.Once
)

// Init 初始化全局 logger。建议在 main() 最开始调用一次。
// level: "debug" | "info" | "warn" | "error"
// env:   "dev" | "prod"
func Init(level, env string) {
	once.Do(func() {
		global = newLogger(level, env)
	})
}

// L 返回全局 logger，未初始化时返回默认 info 级别 logger。
func L() *zap.Logger {
	if global == nil {
		global = newLogger("info", "prod")
	}
	return global
}

// Sync 刷新缓冲日志，程序退出前调用。
func Sync() {
	if global != nil {
		_ = global.Sync()
	}
}

func newLogger(level, env string) *zap.Logger {
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zapLevel = zapcore.InfoLevel
	}

	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "time"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	var encoder zapcore.Encoder
	if env == "dev" {
		// 开发环境：彩色、易读的控制台格式
		devCfg := zap.NewDevelopmentEncoderConfig()
		devCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		devCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(devCfg)
	} else {
		// 生产环境：JSON 结构化日志
		encoder = zapcore.NewJSONEncoder(encCfg)
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout),
		zapLevel,
	)

	opts := []zap.Option{zap.AddCaller(), zap.AddCallerSkip(0)}
	if env == "dev" {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	return zap.New(core, opts...)
}
