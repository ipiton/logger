package logger

import "fmt"

// ConfigError представляет ошибку конфигурации логгера
type ConfigError struct {
	Err    error
	Reason string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("ошибка конфигурации: %s (причина: %v)", e.Reason, e.Err)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}

func (e *WriteError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("ошибка записи лога: %s (причина: %v)", e.Message, e.Cause)
	}
	return fmt.Sprintf("ошибка записи лога: %s", e.Message)
}

func (e *WriteError) Unwrap() error {
	return e.Cause
}

// LevelError представляет ошибку неверного уровня логирования
type LevelError struct {
	Level string
}

func (e *LevelError) Error() string {
	return fmt.Sprintf("неизвестный уровень логирования: %s", e.Level)
}
