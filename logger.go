package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Для возможности тестирования
var osExit = os.Exit

// New создает новый экземпляр логгера с конфигурацией по умолчанию
// Возвращает указатель на Logger, готовый к использованию
func New() *Logger {
	// Создаем конфигурацию по умолчанию
	defaultCfg := DefaultConfig()

	// Создаем базовый логгер
	l := &Logger{
		prefix:     "",
		fields:     make(map[string]interface{}),
		level:      DebugLevel,
		mu:         &sync.RWMutex{},
		fileMu:     sync.Mutex{},
		logger:     log.New(os.Stdout, "", 0),
		timeFormat: "",
		messages:   []string{},
		messagesMu: sync.RWMutex{},
	}

	// Если указан основной файл лога, добавляем его
	if mainLogFile, ok := defaultCfg.Files["main"]; ok && mainLogFile != "" {
		// Создаем директорию для файла логов
		dir := filepath.Dir(mainLogFile)
		if err := os.MkdirAll(dir, 0750); err != nil {
			l.Errorf("Ошибка создания директории для логов: %v", err)
			return l
		}

		f, err := os.OpenFile(filepath.Clean(mainLogFile), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			l.Errorf("Ошибка открытия файла лога: %v", err)
			return l
		}

		l.file = f
		l.logger = log.New(io.MultiWriter(os.Stdout, f), "", 0)
	}

	return l
}

// SetLevel устанавливает уровень логирования
func (l *Logger) SetLevel(level string) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	switch strings.ToLower(level) {
	case "debug":
		l.level = DebugLevel
	case "info":
		l.level = InfoLevel
	case "warning":
		l.level = WarningLevel
	case "error":
		l.level = ErrorLevel
	case "fatal":
		l.level = FatalLevel
	default:
		return fmt.Errorf("неизвестный уровень логирования: %s", level)
	}
	return nil
}

// WithFile создает новый логгер с записью в файл
func (l *Logger) WithFile(filename string) ILogger {
	// Создаем директорию для файла логов
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0750); err != nil {
		l.Errorf("Ошибка создания директории для логов: %v", err)
		return l
	}

	f, err := os.OpenFile(filepath.Clean(filename), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		l.Errorf("Ошибка открытия файла лога: %v", err)
		return l
	}

	newLogger := &Logger{
		prefix:    l.prefix,
		fields:    make(map[string]interface{}),
		file:      f,
		errorFile: l.errorFile,
		level:     l.level,
		mu:        &sync.RWMutex{},
		logger:    log.New(io.MultiWriter(os.Stdout, f), "", 0),
	}

	// Копируем поля
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

func (l *Logger) log(level Level, msg string) {
	if level < l.level {
		return
	}

	if l.prefix != "" {
		msg = "[" + l.prefix + "] " + msg
	}

	formattedMsg := l.formatMessage(level, msg)

	l.mu.Lock()
	defer l.mu.Unlock()

	l.messagesMu.Lock()
	l.messages = append(l.messages, formattedMsg)
	l.messagesMu.Unlock()

	if _, err := fmt.Fprintln(os.Stdout, formattedMsg); err != nil {
		l.handleWriteError(err)
	}

	l.fileMu.Lock()
	defer l.fileMu.Unlock()
	if l.file != nil {
		if _, err := l.file.WriteString(formattedMsg + "\n"); err != nil {
			l.handleWriteError(err)
		}
	}

	// Для ошибок и фатальных ошибок пишем также в файл ошибок
	if (level == ErrorLevel || level == FatalLevel) && l.errorFile != nil {
		if _, err := fmt.Fprintln(l.errorFile, formattedMsg); err != nil {
			l.handleWriteError(err)
		}
		if err := l.errorFile.Sync(); err != nil {
			l.handleWriteError(err)
		}
	}
}

func (l *Logger) logf(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}
	msg := fmt.Sprintf(format, args...)
	l.log(level, msg)
}

// Debug логирует сообщение на уровне DEBUG
func (l *Logger) Debug(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.log(DebugLevel, msg)
}

// Debugf логирует форматированное сообщение на уровне DEBUG
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(DebugLevel, format, args...)
}

// Info логирует сообщение на уровне INFO
func (l *Logger) Info(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.log(InfoLevel, msg)
}

// Infof логирует форматированное сообщение на уровне INFO
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logf(InfoLevel, format, args...)
}

// Warning логирует сообщение на уровне WARNING
func (l *Logger) Warning(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.log(WarningLevel, msg)
}

// Warningf логирует форматированное сообщение на уровне WARNING
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.logf(WarningLevel, format, args...)
}

// Error логирует сообщение на уровне ERROR
func (l *Logger) Error(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.log(ErrorLevel, msg)
}

// Errorf логирует форматированное сообщение на уровне ERROR
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logf(ErrorLevel, format, args...)
}

// Fatal логирует фатальную ошибку и завершает программу
func (l *Logger) Fatal(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.log(FatalLevel, msg)
	if err := l.Close(); err != nil {
		l.Error("failed to close logger:", err)
	}
	osExit(1)
}

// Fatalf логирует фатальную ошибку с форматированием и завершает программу
func (l *Logger) Fatalf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.log(FatalLevel, msg)
	if err := l.Close(); err != nil {
		l.Error("failed to close logger:", err)
	}
	osExit(1)
}

// WithPrefix создает новый логгер с добавленным префиксом
// Префиксы объединяются через точку при вложенных вызовах
// Пример: logger.WithPrefix("API").WithPrefix("V1") -> "[API.V1]"
func (l *Logger) WithPrefix(prefix string) ILogger {
	newLogger := l.clone()
	if l.prefix != "" {
		newLogger.prefix = l.prefix + "." + prefix
	} else {
		newLogger.prefix = prefix
	}
	return newLogger
}

// WithFields создает новый логгер с дополнительными полями
func (l *Logger) WithFields(fields map[string]interface{}) ILogger {
	newLogger := l.clone()
	// Копируем существующие поля
	for k, v := range l.fields {
		newLogger.fields[k] = v
	}
	// Добавляем новые поля
	for k, v := range fields {
		newLogger.fields[k] = v
	}
	return newLogger
}

// Close закрывает файлы логов, если они открыты
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	var closeErr error

	if l.file != nil {
		if syncErr := l.file.Sync(); syncErr != nil {
			l.handleWriteError(syncErr)
		}
		if err := l.file.Close(); err != nil {
			closeErr = err
		}
		l.file = nil
	}

	if l.errorFile != nil {
		if syncErr := l.errorFile.Sync(); syncErr != nil {
			l.handleWriteError(syncErr)
		}
		if err := l.errorFile.Close(); err != nil {
			closeErr = err
		}
		l.errorFile = nil
	}

	return closeErr
}

// clone создает копию логгера
func (l *Logger) clone() *Logger {
	return &Logger{
		prefix:     l.prefix,
		fields:     copyFields(l.fields),
		level:      l.level,
		mu:         l.mu,
		fileMu:     sync.Mutex{},
		logger:     l.logger,
		file:       l.file,
		errorFile:  l.errorFile,
		timeFormat: l.timeFormat,
		messages:   append([]string{}, l.messages...),
		messagesMu: sync.RWMutex{},
	}
}

func (l *Logger) formatMessage(level Level, msg string) string {
	timeFormat := "2006-01-02 15:04:05"
	if l.timeFormat != "" {
		timeFormat = l.timeFormat
	}
	timestamp := time.Now().Format(timeFormat)

	// Получаем строковое представление уровня
	levelStr := ""
	switch level {
	case DebugLevel:
		levelStr = "[DEBUG]"
	case InfoLevel:
		levelStr = "[INFO]"
	case WarningLevel:
		levelStr = "[WARNING]"
	case ErrorLevel:
		levelStr = "[ERROR]"
	case FatalLevel:
		levelStr = "[FATAL]"
	default:
		levelStr = "[UNKNOWN]"
	}

	// Формируем строку с полями
	if len(l.fields) > 0 {
		var fields []string
		for k, v := range l.fields {
			fields = append(fields, fmt.Sprintf("%s=%v", k, v))
		}
		msg += " [" + strings.Join(fields, " ") + "]"
	}

	// Формируем части сообщения
	parts := []string{timestamp, levelStr, msg}

	return strings.Join(parts, " ")
}

// WithLevel создает новый логгер с указанным уровнем логирования
func (l *Logger) WithLevel(level string) ILogger {
	newLogger := l.clone()
	if err := newLogger.SetLevel(level); err != nil {
		if err := newLogger.SetLevel("info"); err != nil {
			newLogger.Error("failed to set level:", err)
		}
	}
	return newLogger
}

// WithConfig применяет конфигурацию к логгеру
func (l *Logger) WithConfig(cfg *Config) ILogger {
	newLogger := l.clone()
	if err := newLogger.SetLevel(cfg.Level); err != nil {
		if err := newLogger.SetLevel("info"); err != nil {
			newLogger.Error("failed to set level:", err)
		}
	}
	// Применяем другие настройки конфигурации...
	return newLogger
}

// WithTimeFormat устанавливает формат времени для логгера
func (l *Logger) WithTimeFormat(format string) ILogger {
	newLogger := l.clone()
	newLogger.timeFormat = format
	return newLogger
}

// GetMessages возвращает список сообщений логгера
func (l *Logger) GetMessages() []string {
	l.messagesMu.RLock()
	defer l.messagesMu.RUnlock()
	return append([]string{}, l.messages...)
}

// Добавляем метод обработки ошибок записи
func (l *Logger) handleWriteError(err error) {
	// Логирование ошибки или другая обработка
}

// Sync принудительно синхронизирует буферы файла логов
func (l *Logger) Sync() error {
	l.fileMu.Lock()
	defer l.fileMu.Unlock()
	if l.file != nil {
		return l.file.Sync()
	}
	return nil
}
