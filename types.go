package logger

import (
	"log"
	"os"
	"sync"
)

// Fields представляет собой набор дополнительных полей для логирования
type Fields map[string]interface{}

// ILogger интерфейс логгера
type ILogger interface {
	// Debug logs a message at DEBUG level
	Debug(args ...interface{})
	// Debugf logs a formatted message at DEBUG level
	Debugf(format string, args ...interface{})
	// Info logs a message at INFO level
	Info(args ...interface{})
	// Infof logs a formatted message at INFO level
	Infof(format string, args ...interface{})
	// Warning logs a message at WARNING level
	Warning(args ...interface{})
	// Warningf logs a formatted message at WARNING level
	Warningf(format string, args ...interface{})
	// Error logs a message at ERROR level
	Error(args ...interface{})
	// Errorf logs a formatted message at ERROR level
	Errorf(format string, args ...interface{})
	// Fatal logs a message at FATAL level and terminates the program
	Fatal(args ...interface{})
	// Fatalf logs a formatted message at FATAL level and terminates the program
	Fatalf(format string, args ...interface{})
	// WithPrefix creates a new logger with the specified prefix
	WithPrefix(prefix string) ILogger
	// WithFields creates a new logger with the specified fields
	WithFields(fields map[string]interface{}) ILogger
	// WithFile creates a new logger that writes to the specified file
	WithFile(filename string) ILogger
	// WithLevel sets the minimum logging level
	WithLevel(level string) ILogger
	// SetLevel sets the minimum logging level
	SetLevel(level string) error
	// Close closes all open file handles
	Close() error
	// WithConfig applies the specified configuration to the logger
	WithConfig(cfg *Config) ILogger
	// WithTimeFormat sets the time format for the logger
	WithTimeFormat(format string) ILogger
	// GetMessages returns the messages logged by the logger
	GetMessages() []string
}

// Level представляет уровень логирования
type Level int

// Уровни логирования
const (
	DebugLevel Level = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

// Алиасы для удобства использования
const (
	DEBUG   = DebugLevel
	INFO    = InfoLevel
	WARNING = WarningLevel
	ERROR   = ErrorLevel
	FATAL   = FatalLevel
)

// Logger реализует интерфейс ILogger и предоставляет функциональность для логирования
type Logger struct {
	mu         *sync.RWMutex
	logger     *log.Logger
	file       *os.File
	errorFile  *os.File
	fields     map[string]interface{}
	prefix     string
	timeFormat string
	messages   []string
	level      Level
	messagesMu sync.RWMutex
	fileMu     sync.Mutex
}

// Config представляет конфигурацию логгера
type Config struct {
	Files map[string]string
	Level string
}

// WriteError представляет ошибку записи в лог
type WriteError struct {
	Cause   error
	Message string
}
