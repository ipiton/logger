// Package logger предоставляет функциональность для логирования
package logger

import (
	"fmt"
	"strings"
)

// DefaultConfig возвращает конфигурацию логгера по умолчанию
func DefaultConfig() *Config {
	return &Config{
		Level: "info",
		Files: map[string]string{
			"main": "", // Убираем файл ошибок из конфигурации по умолчанию
		},
	}
}

// Override переопределяет значения конфигурации из другой конфигурации
func (c *Config) Override(other *Config) {
	if other == nil {
		return
	}

	if other.Level != "" {
		c.Level = other.Level
	}

	if other.Files != nil {
		if c.Files == nil {
			c.Files = make(map[string]string)
		}
		for k, v := range other.Files {
			c.Files[k] = v
		}
	}
}

// Validate проверяет корректность конфигурации
func (c *Config) Validate() error {
	// Проверяем уровень логирования
	validLevels := map[string]bool{
		"debug":   true,
		"info":    true,
		"warning": true,
		"error":   true,
		"fatal":   true,
	}
	if !validLevels[strings.ToLower(c.Level)] {
		return fmt.Errorf("некорректный уровень логирования: %s", c.Level)
	}

	// Проверяем наличие основного файла лога
	if _, ok := c.Files["main"]; !ok {
		return fmt.Errorf("не указан путь к основному файлу лога")
	}

	return nil
}
