package logx

import (
	"os"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// logger 是当前实际对外提供服务的 sugared logger。
	// 选择 SugaredLogger 是为了兼容项目现有的 `msg + key/value` 与 `Infof` 两种调用方式。
	logger = newDefaultLogger()
)

// newDefaultLogger 创建一个仅输出到标准输出的默认 logger。
// 即使配置尚未初始化，日志也不会完全丢失。
func newDefaultLogger() *zap.SugaredLogger {
	baseLogger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	), zap.AddCaller(), zap.AddCallerSkip(1))
	return baseLogger.Sugar()
}

// InitFromViper 使用现有 viper 配置重建日志系统。
//
// 当前实现采用 zap + lumberjack：
// 1. zap 负责结构化日志与高性能编码
// 2. lumberjack 负责文件滚动
//
// 这样我们就不需要自己维护复杂的文件句柄和滚动细节。
func InitFromViper() error {
	encoder := buildEncoder(viper.GetBool("log.log_format_text"))
	syncer := buildWriteSyncer()
	level := parseLevel(viper.GetString("log.logger_level"))

	baseLogger := zap.New(
		zapcore.NewCore(encoder, syncer, level),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)

	oldLogger := logger
	logger = baseLogger.Sugar()

	if oldLogger != nil {
		_ = oldLogger.Sync()
	}
	_ = baseLogger.Sync()
	return nil
}

func Info(msg string, args ...any) {
	current := getLogger()
	if len(args) == 0 {
		current.Info(msg)
		return
	}
	current.Infow(msg, sanitizeKeyValues(args)...)
}

func Infof(format string, args ...any) {
	getLogger().Infof(format, args...)
}

func Warn(msg string, args ...any) {
	current := getLogger()
	if len(args) == 0 {
		current.Warn(msg)
		return
	}
	current.Warnw(msg, sanitizeKeyValues(args)...)
}

func Error(msg string, args ...any) {
	current := getLogger()
	if len(args) == 0 {
		current.Error(msg)
		return
	}
	current.Errorw(msg, sanitizeKeyValues(args)...)
}

func Errorf(format string, args ...any) {
	getLogger().Errorf(format, args...)
}

func Fatal(msg string, args ...any) {
	current := getLogger()
	_ = current.Sync()
	if len(args) == 0 {
		current.Fatal(msg)
	} else {
		current.Fatalw(msg, sanitizeKeyValues(args)...)
	}

}

// getLogger 返回当前可用的 logger。
func getLogger() *zap.SugaredLogger {
	return logger
}

// parseLevel 把历史配置里的字符串等级转换成 zap 等级。
func parseLevel(raw string) zapcore.Level {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "DEBUG":
		return zapcore.DebugLevel
	case "WARN", "WARNING":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// buildEncoder 根据现有配置决定使用文本格式还是 JSON 格式。
func buildEncoder(useText bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "msg"
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	if useText {
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
	return zapcore.NewJSONEncoder(encoderConfig)
}

// buildWriteSyncer 根据配置拼装输出目标。
// 当配置里包含 file 时，底层由 lumberjack 托管滚动策略。
func buildWriteSyncer() zapcore.WriteSyncer {
	names := strings.Split(strings.TrimSpace(viper.GetString("log.writers")), ",")
	syncers := make([]zapcore.WriteSyncer, 0, len(names))

	for _, name := range names {
		switch strings.TrimSpace(name) {
		case "stdout":
			syncers = append(syncers, zapcore.AddSync(os.Stdout))
		case "file":
			rotatingWriter := newRollingFileWriter()
			if rotatingWriter != nil {
				syncers = append(syncers, zapcore.AddSync(rotatingWriter))
			}
		}
	}

	if len(syncers) == 0 {
		return zapcore.AddSync(os.Stdout)
	}
	return zapcore.NewMultiWriteSyncer(syncers...)
}

// newRollingFileWriter 使用现有配置创建一个 lumberjack logger。
// 这样日志文件的滚动和底层文件打开/关闭逻辑都交给成熟库维护。
func newRollingFileWriter() *lumberjack.Logger {
	filename := strings.TrimSpace(viper.GetString("log.logger_file"))
	if filename == "" {
		return nil
	}

	maxSize := viper.GetInt("log.log_rotate_size")
	if maxSize <= 0 {
		maxSize = 100
	}

	maxBackups := viper.GetInt("log.log_backup_count")
	if maxBackups < 0 {
		maxBackups = 0
	}

	maxAge := viper.GetInt("log.log_rotate_date")
	if maxAge < 0 {
		maxAge = 0
	}

	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Compress:   false,
	}
}

// sanitizeKeyValues 保证传给 zap 的 key/value 列表是成对的。
// 如果调用方传了奇数个参数，就自动补一个占位 key，避免日志直接报错。
func sanitizeKeyValues(args []any) []any {
	if len(args)%2 == 0 {
		return args
	}

	result := make([]any, 0, len(args)+1)
	result = append(result, args...)
	result = append(result, "(missing)")
	return result
}
