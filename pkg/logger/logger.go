package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Log 全局日志实例
	Log *zap.Logger
	// Sugar 全局 SugaredLogger 实例（支持类似 printf 的格式化输出）
	Sugar *zap.SugaredLogger
)

// Init 初始化日志系统
// env: "dev" (控制台彩色输出) 或 "prod" (JSON格式 + 文件写入)
func Init(env string) {
	var core zapcore.Core

	if env == "dev" {
		// 开发环境：控制台输出，彩色，人类可读
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ := config.Build()
		Log = logger
	} else {
		// 生产环境：JSON格式，写入文件，带日志轮转
		writeSyncer := getLogWriter()
		encoder := getEncoder()
		core = zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)

		// 添加调用者信息（文件名+行号），但在高并发下会有性能损耗，按需开启
		Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	Sugar = Log.Sugar()
}

// getEncoder 配置 JSON 编码器
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // ISO8601 时间格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// getLogWriter 配置日志文件轮转（Lumberjack）
func getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   "./logs/app.log", // 日志文件路径
		MaxSize:    10,               // 每个文件最大 10MB
		MaxBackups: 5,                // 保留 5 个旧文件
		MaxAge:     30,               // 保留 30 天
		Compress:   true,             // 压缩旧日志
	}
	return zapcore.AddSync(lumberJackLogger)
}

// --- 常用便捷方法封装 ---

func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Debug 仅在开发模式有效
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Fatal 记录致命错误并退出程序
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}
