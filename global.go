package logger

import (
	"sync"
)

var (
	globalLogger ILogger
	globalMu     sync.RWMutex
)

// SetGlobalLogger устанавливает глобальный логгер
func SetGlobalLogger(l ILogger) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalLogger = l
}

// GetGlobalLogger возвращает текущий глобальный логгер
func GetGlobalLogger() ILogger {
	globalMu.RLock()
	defer globalMu.RUnlock()
	if globalLogger == nil {
		return New()
	}
	return globalLogger
}

// Debug логирует сообщение на уровне DEBUG
func Debug(args ...interface{}) { GetGlobalLogger().Debug(args...) }

// Debugf логирует форматированное сообщение на уровне DEBUG
func Debugf(format string, args ...interface{}) { GetGlobalLogger().Debugf(format, args...) }

// Info логирует сообщение на уровне INFO
func Info(args ...interface{}) { GetGlobalLogger().Info(args...) }

// Infof логирует форматированное сообщение на уровне INFO
func Infof(format string, args ...interface{}) { GetGlobalLogger().Infof(format, args...) }

// Warning логирует сообщение на уровне WARNING
func Warning(args ...interface{}) { GetGlobalLogger().Warning(args...) }

// Warningf логирует форматированное сообщение на уровне WARNING
func Warningf(format string, args ...interface{}) { GetGlobalLogger().Warningf(format, args...) }

// Error логирует сообщение на уровне ERROR
func Error(args ...interface{}) { GetGlobalLogger().Error(args...) }

// Errorf логирует форматированное сообщение на уровне ERROR
func Errorf(format string, args ...interface{}) { GetGlobalLogger().Errorf(format, args...) }

// Fatal логирует сообщение на уровне FATAL и завершает программу
func Fatal(args ...interface{}) { GetGlobalLogger().Fatal(args...) }

// Fatalf логирует форматированное сообщение на уровне FATAL и завершает программу
func Fatalf(format string, args ...interface{}) { GetGlobalLogger().Fatalf(format, args...) }

// WithPrefix создает новый логгер с указанным префиксом
func WithPrefix(prefix string) ILogger { return GetGlobalLogger().WithPrefix(prefix) }

// WithFields создает новый логгер с указанными полями
func WithFields(fields map[string]interface{}) ILogger { return GetGlobalLogger().WithFields(fields) }

// WithFile создает новый логгер с записью в указанный файл
func WithFile(filename string) ILogger { return GetGlobalLogger().WithFile(filename) }
