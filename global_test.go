package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetAndGetGlobalLogger(t *testing.T) {
	// Сохраняем текущий глобальный логгер
	originalLogger := globalLogger

	defer func() {
		// Восстанавливаем оригинальный логгер после теста
		globalLogger = originalLogger
	}()

	// Создаем тестовый логгер
	testLogger := NewMockLogger()
	SetGlobalLogger(testLogger)

	// Проверяем, что глобальный логгер установлен
	retrievedLogger := GetGlobalLogger()
	assert.Equal(t, testLogger, retrievedLogger, "Глобальный логгер должен быть установлен корректно")
}

func TestGetGlobalLoggerDefault(t *testing.T) {
	// Сохраняем текущий глобальный логгер
	originalLogger := globalLogger

	defer func() {
		// Восстанавливаем оригинальный логгер после теста
		globalLogger = originalLogger
	}()

	// Сбрасываем глобальный логгер
	globalLogger = nil

	// Получаем глобальный логгер
	defaultLogger := GetGlobalLogger()

	// Проверяем, что возвращен новый логгер по умолчанию
	assert.NotNil(t, defaultLogger, "Должен быть возвращен новый логгер")
	assert.IsType(t, &Logger{}, defaultLogger, "Должен быть возвращен логгер типа Logger")
}

func TestGlobalLogger(t *testing.T) {
	// Сохраняем текущий глобальный логгер
	originalLogger := globalLogger
	defer func() {
		globalLogger = originalLogger
	}()

	// Создаем мок логгер
	mockLogger := NewMockLogger()
	SetGlobalLogger(mockLogger)

	t.Run("базовые методы логирования", func(t *testing.T) {
		Debug("debug message")
		Info("info message")
		Warning("warning message")
		Error("error message")

		assert.Equal(t, "[DEBUG] debug message", mockLogger.Messages[0])
		assert.Equal(t, "[INFO] info message", mockLogger.Messages[1])
		assert.Equal(t, "[WARNING] warning message", mockLogger.Messages[2])
		assert.Equal(t, "[ERROR] error message", mockLogger.Messages[3])
	})

	t.Run("методы с префиксом и полями", func(t *testing.T) {
		WithPrefix("TEST").Info("test prefix")
		assert.Equal(t, "[INFO] TEST test prefix", mockLogger.Messages[4])

		WithPrefix("TEST").WithFields(map[string]interface{}{"test": "field"}).Info("with fields")
		assert.Equal(t, "[INFO] TEST with fields [test=field]", mockLogger.Messages[5])
	})

	t.Run("получение глобального логгера", func(t *testing.T) {
		logger := GetGlobalLogger()
		assert.NotNil(t, logger)
	})

	t.Run("создание логгера по умолчанию", func(t *testing.T) {
		SetGlobalLogger(nil)
		logger := GetGlobalLogger()
		assert.NotNil(t, logger)
	})
}

func TestGlobalLoggerWithMethods(t *testing.T) {
	originalLogger := globalLogger
	defer func() {
		globalLogger = originalLogger
	}()

	mockLogger := NewMockLogger()
	SetGlobalLogger(mockLogger)

	t.Run("логирование с префиксом", func(t *testing.T) {
		WithPrefix("TEST").Info("test message")
		assert.Equal(t, "[INFO] TEST test message", mockLogger.Messages[0])
	})

	t.Run("логирование с полями", func(t *testing.T) {
		WithFields(map[string]interface{}{"key": "value"}).Info("test message")
		assert.Equal(t, "[INFO] test message [key=value]", mockLogger.Messages[1])
	})

	t.Run("логирование в файл", func(t *testing.T) {
		fileLogger := WithFile("/tmp/test.log")
		assert.NotNil(t, fileLogger)

		fileLogger.Info("test message")
		assert.Contains(t, mockLogger.Messages[2], "[INFO] test message")
	})
}

func TestFatalMethod(t *testing.T) {
	originalLogger := globalLogger
	defer func() {
		globalLogger = originalLogger
	}()

	// Создаем мок логгер
	mockLogger := NewMockLogger()
	SetGlobalLogger(mockLogger)

	// Сохраняем оригинальную функцию osExit
	originalOsExit := osExit
	defer func() { osExit = originalOsExit }()

	// Подменяем osExit для тестирования
	var exitCalled bool
	mockLogger.SetExitHandler(func(code int) {
		exitCalled = true
		assert.Equal(t, 1, code, "Код выхода должен быть 1")
	})

	// Проверяем, что вызов Fatal не вызывает панику
	assert.NotPanics(t, func() {
		Fatal("fatal message")
	}, "Метод Fatal не должен вызывать панику")

	// Проверяем, что сообщение было добавлено
	assert.Len(t, mockLogger.Messages, 1, "Должно быть одно сообщение")
	assert.Equal(t, "[FATAL] fatal message", mockLogger.Messages[0], "Сообщение должно быть с уровнем FATAL")
	assert.True(t, exitCalled, "Должен быть вызван os.Exit")
}

func TestGlobalLoggerFormatMethods(t *testing.T) {
	originalLogger := globalLogger
	defer func() { globalLogger = originalLogger }()

	mockLogger := NewMockLogger()
	SetGlobalLogger(mockLogger)

	t.Run("Debugf", func(t *testing.T) {
		Debugf("debug %s", "test")
		assert.Contains(t, mockLogger.Messages[0], "[DEBUG] debug test")
	})

	t.Run("Infof", func(t *testing.T) {
		Infof("info %d", 42)
		assert.Contains(t, mockLogger.Messages[1], "[INFO] info 42")
	})

	t.Run("Warningf", func(t *testing.T) {
		Warningf("warning %t", true)
		assert.Contains(t, mockLogger.Messages[2], "[WARNING] warning true")
	})

	t.Run("Errorf", func(t *testing.T) {
		Errorf("error %v", "test")
		assert.Contains(t, mockLogger.Messages[3], "[ERROR] error test")
	})

	t.Run("Fatalf", func(t *testing.T) {
		exitCode := 0
		mockLogger.SetExitHandler(func(code int) { exitCode = code })

		Fatalf("fatal %s", "error")
		assert.Contains(t, mockLogger.Messages[4], "[FATAL] fatal error")
		assert.Equal(t, 1, exitCode)
	})
}
