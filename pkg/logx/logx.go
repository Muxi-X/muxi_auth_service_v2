package logx

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelInfo,
}))

// InitFromViper 从现有 viper 配置初始化日志系统。
// 这里使用标准库 slog，减少第三方依赖，同时尽量兼容项目原本的配置字段。
func InitFromViper() error {
	writers := buildWriters(viper.GetString("log.writers"), viper.GetString("log.logger_file"))
	level := parseLevel(viper.GetString("log.logger_level"))

	var handler slog.Handler
	if viper.GetBool("log.log_format_text") {
		handler = slog.NewTextHandler(writers, &slog.HandlerOptions{Level: level})
	} else {
		handler = slog.NewJSONHandler(writers, &slog.HandlerOptions{Level: level})
	}

	logger = slog.New(handler)
	slog.SetDefault(logger)
	return nil
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Infof(format string, args ...any) {
	logger.Info(fmt.Sprintf(format, args...))
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Errorf(format string, args ...any) {
	logger.Error(fmt.Sprintf(format, args...))
}

func Fatal(msg string, args ...any) {
	logger.Error(msg, args...)
	os.Exit(1)
}

func parseLevel(raw string) slog.Level {
	switch strings.ToUpper(strings.TrimSpace(raw)) {
	case "DEBUG":
		return slog.LevelDebug
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func buildWriters(rawWriters string, logFile string) io.Writer {
	names := strings.Split(strings.TrimSpace(rawWriters), ",")
	writers := make([]io.Writer, 0, len(names))

	for _, name := range names {
		switch strings.TrimSpace(name) {
		case "stdout":
			writers = append(writers, os.Stdout)
		case "file":
			if fileWriter := openLogFile(logFile); fileWriter != nil {
				writers = append(writers, fileWriter)
			}
		}
	}

	if len(writers) == 0 {
		return os.Stdout
	}
	if len(writers) == 1 {
		return writers[0]
	}
	return io.MultiWriter(writers...)
}

func openLogFile(logFile string) io.Writer {
	logFile = strings.TrimSpace(logFile)
	if logFile == "" {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(logFile), 0o755); err != nil {
		return nil
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil
	}

	return file
}
