package weboffice

import (
	"fmt"
	"time"
)

// Logger 定义log接口
type Logger interface {
	Info(format string, args ...any)
	Error(format string, args ...any)
}

// DefaultLogger 默认logger
func DefaultLogger() Logger {
	return &stdLogger{}
}

// stdLogger 标准logger
type stdLogger struct {
}

// Info info日志
func (*stdLogger) Info(format string, args ...any) {
	tm := time.Now().Format("2006-01-02T15:04:05.000")
	println(tm, "[INFO]", fmt.Sprintf(format, args...))
}

// Error error日志
func (*stdLogger) Error(format string, args ...any) {
	tm := time.Now().Format("2006-01-02T15:04:05.000")
	println(tm, "[ERROR]", fmt.Sprintf(format, args...))
	
}

// noopLogger 空logger
type noopLogger struct {
}

// Info info日志
func (*noopLogger) Info(format string, args ...any) {
}

// Error error日志
func (*noopLogger) Error(format string, args ...any) {
}
