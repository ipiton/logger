package logger

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Импортируем переменную osExit
var _ = osExit // Используем для доступа к глобальной переменной

// MockLogger реализация для тестов
// nolint:govet
type MockLogger struct {
	mu         sync.RWMutex           // 24 bytes
	Messages   []string               // 24 bytes
	prefix     string                 // 16 bytes
	osExitFunc func(int)              // 8 bytes
	parent     *MockLogger            // 8 bytes
	fields     map[string]interface{} // 8 bytes
	level      Level                  // 4 bytes
	exitCode   int                    // 4 bytes
	hasFile    bool                   // 1 byte
}

// NewMockLogger creates a new instance of MockLogger.
func NewMockLogger() *MockLogger {
	return &MockLogger{
		Messages:   []string{},
		fields:     make(map[string]interface{}),
		level:      DebugLevel,
		mu:         sync.RWMutex{},
		osExitFunc: func(code int) { /* default - ничего не делаем */ },
	}
}

// formatMessage форматирует сообщение с учетом префикса и полей
func (m *MockLogger) formatMessage(level Level, args ...interface{}) string {
	msg := strings.TrimSuffix(fmt.Sprintln(args...), "\n")
	levelStr := getLevelString(level)
	prefixes, allFields := m.collectChainData()

	parts := []string{levelStr}
	if len(prefixes) > 0 {
		parts = append(parts, strings.Join(prefixes, "."))
	}
	parts = append(parts, msg)

	if len(allFields) > 0 {
		fieldParts := make([]string, 0, len(allFields))
		for k, v := range allFields {
			fieldParts = append(fieldParts, fmt.Sprintf("%s=%v", k, v))
		}
		sort.Strings(fieldParts)
		parts = append(parts, "["+strings.Join(fieldParts, " ")+"]")
	}

	return strings.Join(parts, " ")
}

// collectChainData собирает данные по всей цепочке логгеров
func (m *MockLogger) collectChainData() ([]string, map[string]interface{}) {
	var prefixes []string
	fields := make(map[string]interface{})

	current := m
	for current != nil {
		if current.prefix != "" {
			prefixes = append([]string{current.prefix}, prefixes...)
		}
		for k, v := range current.fields {
			fields[k] = v
		}
		current = current.parent
	}

	return prefixes, fields
}

// log добавляет сообщение в слайс сообщений
func (m *MockLogger) log(level Level, args ...interface{}) {
	// Всегда используем корневой логгер для проверки уровня
	root := m.getRootLogger()
	if level < root.level {
		return
	}

	root.mu.Lock()
	defer root.mu.Unlock()
	root.Messages = append(root.Messages, m.formatMessage(level, args...))
}

// logf добавляет форматированное сообщение в слайс сообщений
func (m *MockLogger) logf(level Level, format string, args ...interface{}) {
	// Всегда используем корневой логгер для проверки уровня
	root := m.getRootLogger()
	if level < root.level {
		return
	}

	root.mu.Lock()
	defer root.mu.Unlock()
	root.Messages = append(root.Messages, m.formatMessage(level, fmt.Sprintf(format, args...)))
}

// Debug логирует отладочное сообщение
func (m *MockLogger) Debug(args ...interface{}) {
	m.log(DebugLevel, args...)
}

// Debugf логирует отладочное сообщение с форматированием
func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.logf(DebugLevel, format, args...)
}

// Info логирует информационное сообщение
func (m *MockLogger) Info(args ...interface{}) {
	m.log(InfoLevel, args...)
}

// Infof логирует информационное сообщение с форматированием
func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.logf(InfoLevel, format, args...)
}

// Warning логирует предупреждение
func (m *MockLogger) Warning(args ...interface{}) {
	m.log(WarningLevel, args...)
}

// Warningf логирует предупреждение с форматированием
func (m *MockLogger) Warningf(format string, args ...interface{}) {
	m.logf(WarningLevel, format, args...)
}

// Error логирует ошибку
func (m *MockLogger) Error(args ...interface{}) {
	m.log(ErrorLevel, args...)
}

// Errorf логирует ошибку с форматированием
func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.logf(ErrorLevel, format, args...)
}

// Fatal логирует фатальную ошибку и завершает программу
func (m *MockLogger) Fatal(args ...interface{}) {
	m.log(FatalLevel, args...)
	if err := m.Close(); err != nil {
		m.Error("failed to close logger:", err)
	}
	m.exitCode = 1
	m.osExitFunc(1)
}

// Fatalf логирует фатальную ошибку с форматированием и завершает программу
func (m *MockLogger) Fatalf(format string, args ...interface{}) {
	m.logf(FatalLevel, format, args...)
	if err := m.Close(); err != nil {
		m.Error("failed to close logger:", err)
	}
	m.exitCode = 1
	m.osExitFunc(1)
}

// WithPrefix создает новый логгер с префиксом
func (m *MockLogger) WithPrefix(prefix string) ILogger {
	newLogger := &MockLogger{
		prefix:     prefix,
		fields:     make(map[string]interface{}),
		level:      m.level,
		mu:         sync.RWMutex{},
		parent:     m,
		osExitFunc: m.osExitFunc,
	}
	return newLogger
}

// WithFields создает новый логгер с полями
func (m *MockLogger) WithFields(fields map[string]interface{}) ILogger {
	newLogger := &MockLogger{
		prefix:     "",
		fields:     fields,
		level:      m.level,
		mu:         sync.RWMutex{},
		parent:     m,
		osExitFunc: m.osExitFunc,
	}
	return newLogger
}

// WithFile создает новый логгер с файлом
func (m *MockLogger) WithFile(filename string) ILogger {
	newLogger := &MockLogger{
		prefix:     "",
		fields:     make(map[string]interface{}),
		level:      m.level,
		mu:         sync.RWMutex{},
		hasFile:    true,
		parent:     m,
		osExitFunc: m.osExitFunc,
	}
	return newLogger
}

// SetLevel устанавливает уровень логирования
func (m *MockLogger) SetLevel(level string) error {
	switch strings.ToLower(level) {
	case "debug":
		m.level = DebugLevel
	case "info":
		m.level = InfoLevel
	case "warning":
		m.level = WarningLevel
	case "error":
		m.level = ErrorLevel
	case "fatal":
		m.level = FatalLevel
	default:
		return fmt.Errorf("неизвестный уровень логирования: %s", level)
	}
	return nil
}

// Close закрывает логгер
func (m *MockLogger) Close() error {
	return nil
}

// SetExitHandler устанавливает обработчик для os.Exit
func (m *MockLogger) SetExitHandler(fn func(int)) {
	m.osExitFunc = fn
}

// GetExitCode возвращает код выхода последнего вызова Fatal
func (m *MockLogger) GetExitCode() int {
	return m.exitCode
}

// getRootLogger возвращает корневой логгер в цепочке
func (m *MockLogger) getRootLogger() *MockLogger {
	if m.parent != nil {
		return m.parent.getRootLogger()
	}
	return m
}

// getLevelString возвращает строковое представление уровня логирования
func getLevelString(level Level) string {
	switch level {
	case DebugLevel:
		return "[DEBUG]"
	case InfoLevel:
		return "[INFO]"
	case WarningLevel:
		return "[WARNING]"
	case ErrorLevel:
		return "[ERROR]"
	case FatalLevel:
		return "[FATAL]"
	default:
		return "[UNKNOWN]"
	}
}

// WithConfig применяет конфигурацию к логгеру
func (m *MockLogger) WithConfig(cfg *Config) ILogger {
	return m // Простая реализация для мока
}

// WithLevel создает новый логгер с указанным уровнем
func (m *MockLogger) WithLevel(level string) ILogger {
	newLogger := m.clone()
	newLogger.level = parseLevel(level)
	return newLogger
}

// parseLevel преобразует строку в Level
func parseLevel(level string) Level {
	switch strings.ToLower(level) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warning":
		return WarningLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// clone создает копию MockLogger
func (m *MockLogger) clone() *MockLogger {
	return &MockLogger{
		Messages:   append([]string{}, m.Messages...),
		exitCode:   m.exitCode,
		osExitFunc: m.osExitFunc,
		prefix:     m.prefix,
		fields:     copyFields(m.fields),
		level:      m.level,
		mu:         sync.RWMutex{},
		hasFile:    m.hasFile,
		parent:     m.parent,
	}
}

// copyFields создает копию map полей
func copyFields(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// GetMessages возвращает список записанных сообщений
func (m *MockLogger) GetMessages() []string {
	return m.Messages
}

// WithTimeFormat устанавливает формат времени для логгера
func (m *MockLogger) WithTimeFormat(format string) ILogger {
	newLogger := m.clone()
	return newLogger
}
