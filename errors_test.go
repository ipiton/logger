package logger

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigError(t *testing.T) {
	err := &ConfigError{
		Err:    fmt.Errorf("invalid log level"),
		Reason: "некорректное значение уровня логирования",
	}

	assert.Equal(t,
		"ошибка конфигурации: некорректное значение уровня логирования (причина: invalid log level)",
		err.Error(),
		"Сообщение об ошибке должно быть сформировано корректно",
	)
}

func TestWriteError(t *testing.T) {
	t.Run("WithCause", func(t *testing.T) {
		originalErr := errors.New("permission denied")
		err := &WriteError{
			Message: "не удалось записать лог",
			Cause:   originalErr,
		}

		assert.Equal(t,
			"ошибка записи лога: не удалось записать лог (причина: permission denied)",
			err.Error(),
			"Сообщение об ошибке должно включать причину",
		)

		assert.Equal(t, originalErr, err.Unwrap(), "Unwrap должен возвращать исходную ошибку")
	})

	t.Run("WithoutCause", func(t *testing.T) {
		err := &WriteError{
			Message: "не удалось записать лог",
		}

		assert.Equal(t,
			"ошибка записи лога: не удалось записать лог",
			err.Error(),
			"Сообщение об ошибке должно быть корректным без причины",
		)

		assert.Nil(t, err.Unwrap(), "Unwrap должен возвращать nil без причины")
	})
}

func TestLevelError(t *testing.T) {
	err := &LevelError{
		Level: "SUPER_DEBUG",
	}

	assert.Equal(t,
		"неизвестный уровень логирования: SUPER_DEBUG",
		err.Error(),
		"Сообщение об ошибке должно быть сформировано корректно",
	)
}

func TestErrorWrapping(t *testing.T) {
	originalErr := fmt.Errorf("original error")
	configErr := &ConfigError{Err: originalErr, Reason: "invalid config"}

	if !errors.Is(configErr, originalErr) {
		t.Error("Ошибка должна содержать оригинальную ошибку")
	}
}

func TestNew(t *testing.T) {
	err := New().SetLevel("invalid")
	assert.Equal(t,
		"неизвестный уровень логирования: invalid",
		err.Error(),
		"Сообщение об ошибке должно быть сформировано корректно",
	)
}
