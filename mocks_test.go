package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockLogger(t *testing.T) {
	mock := NewMockLogger()

	// Настраиваем перехватчик
	exitCode := 0
	mock.SetExitHandler(func(code int) {
		exitCode = code
	})

	t.Run("Debug", func(t *testing.T) {
		mock.Debug("test")
		require.Len(t, mock.Messages, 1)
		assert.Equal(t, "[DEBUG] test", mock.Messages[0])
		assert.Equal(t, 0, exitCode)
	})

	t.Run("Info", func(t *testing.T) {
		mock.Info("info message")
		require.Len(t, mock.Messages, 2)
		assert.Equal(t, "[INFO] info message", mock.Messages[1])
		assert.Equal(t, 0, exitCode)
	})

	t.Run("Warning", func(t *testing.T) {
		mock.Warning("warning")
		require.Len(t, mock.Messages, 3)
		assert.Equal(t, "[WARNING] warning", mock.Messages[2])
		assert.Equal(t, 0, exitCode)
	})

	t.Run("Error", func(t *testing.T) {
		mock.Error("error")
		require.Len(t, mock.Messages, 4)
		assert.Equal(t, "[ERROR] error", mock.Messages[3])
		assert.Equal(t, 0, exitCode)
	})

	t.Run("Fatal", func(t *testing.T) {
		mock.Fatal("fatal")
		assert.Equal(t, 1, exitCode, "Код выхода должен быть 1")
	})
}

func TestMockLoggerMultipleMessages(t *testing.T) {
	mockLogger := NewMockLogger()

	// Настраиваем перехватчик
	exitCode := 0
	mockLogger.SetExitHandler(func(code int) {
		exitCode = code
	})

	// Логируем несколько сообщений
	mockLogger.Debug("debug1")
	mockLogger.Info("info1")
	mockLogger.Warning("warning1")
	mockLogger.Error("error1")
	mockLogger.Fatal("fatal1")

	// Проверяем код выхода после Fatal
	assert.Equal(t, 1, exitCode, "Должен быть установлен код выхода 1")

	// Проверяем, что все сообщения были добавены
	expected := []string{
		"[DEBUG] debug1",
		"[INFO] info1",
		"[WARNING] warning1",
		"[ERROR] error1",
		"[FATAL] fatal1",
	}

	assert.Equal(t, expected, mockLogger.Messages, "Сообщения должны полностью совпадать с ожидаемыми")
}

func TestMockLoggerWithMultipleArguments(t *testing.T) {
	mockLogger := NewMockLogger()

	// Логируем сообщение с несколькими аргументами
	mockLogger.Info("test", 123, true)

	// Проверяем, что сообщение было корректно сформировано
	assert.Len(t, mockLogger.Messages, 1, "Должно быть одно сообщение")
	assert.Equal(t, "[INFO] test 123 true", mockLogger.Messages[0], "Сообщение должно содержать все аргументы с пробелами")
}

func TestMockLoggerWithPrefix(t *testing.T) {
	mockLogger := NewMockLogger()
	prefixedLogger := mockLogger.WithPrefix("TEST")
	prefixedLogger.Info("message")

	assert.Len(t, mockLogger.Messages, 1)
	assert.Equal(t, "[INFO] TEST message", mockLogger.Messages[0])
}

func TestMockLoggerWithFields(t *testing.T) {
	mockLogger := NewMockLogger()

	// Создаем логгер с полями
	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 2,
	}
	fieldLogger := mockLogger.WithFields(fields)
	fieldLogger.Info("message")

	// Проверяем, что сообщение содержит поля
	assert.Len(t, mockLogger.Messages, 1, "Должно быть одно сообщение")
	assert.Equal(t, "[INFO] message [key1=value1 key2=2]", mockLogger.Messages[0], "Сообщение должно содержать поля")
}

func TestMockLoggerChainedMethods(t *testing.T) {
	mockLogger := NewMockLogger()

	logger := mockLogger.
		WithPrefix("PREFIX").
		WithFields(map[string]interface{}{"key": "value"}).
		WithFile("test.log")

	logger.Info("chained message")

	expected := "[INFO] PREFIX chained message [key=value]"
	assert.Equal(t, []string{expected}, mockLogger.Messages)
}

func TestMockLoggerWithFile(t *testing.T) {
	mockLogger := NewMockLogger()

	// Проверяем, что WithFile не влияет на работу логгера
	fileLogger := mockLogger.WithFile("test.log")
	fileLogger.Info("message")

	assert.Len(t, mockLogger.Messages, 1, "Должно быть одно сообщение")
	assert.Contains(t, mockLogger.Messages[0], "message", "Сообщение должно быть записано")
}

func TestMockLoggerNestedFields(t *testing.T) {
	mockLogger := NewMockLogger()

	// Тестируем вложенные поля
	logger1 := mockLogger.WithFields(map[string]interface{}{"key1": "value1"})
	logger2 := logger1.WithFields(map[string]interface{}{"key2": "value2"})
	logger2.Info("message")

	assert.Len(t, mockLogger.Messages, 1, "Должно быть одно сообщение")
	assert.Contains(t, mockLogger.Messages[0], "message [key1=value1 key2=value2]", "Сообщение должно содержать все поля")
}

func TestMockLoggerNestedPrefixes(t *testing.T) {
	mockLogger := NewMockLogger()
	logger1 := mockLogger.WithPrefix("PREFIX1")
	logger2 := logger1.WithPrefix("PREFIX2")
	logger2.Info("message")

	assert.Contains(t, mockLogger.Messages[0], "[INFO] PREFIX1.PREFIX2 message")
}

func TestMockLoggerComplexScenario(t *testing.T) {
	mockLogger := NewMockLogger()
	logger := mockLogger.WithPrefix("SERVICE").WithFields(map[string]interface{}{"version": "1.0"})

	logger.Debug("debug message")
	logger.WithFields(map[string]interface{}{"user": "admin"}).Info("info message")

	assert.Contains(t, mockLogger.Messages[0], "[DEBUG] SERVICE debug message [version=1.0]")
	assert.Contains(t, mockLogger.Messages[1], "[INFO] SERVICE info message [user=admin version=1.0]")
}

func TestMockLoggerFormatMethods(t *testing.T) {
	mockLogger := NewMockLogger()

	t.Run("Debugf", func(t *testing.T) {
		mockLogger.Debugf("debug %s", "test")
		assert.Contains(t, mockLogger.Messages[0], "[DEBUG] debug test")
	})

	t.Run("Infof", func(t *testing.T) {
		mockLogger.Infof("info %d", 42)
		assert.Contains(t, mockLogger.Messages[1], "[INFO] info 42")
	})

	t.Run("Warningf", func(t *testing.T) {
		mockLogger.Warningf("warning %t", true)
		assert.Contains(t, mockLogger.Messages[2], "[WARNING] warning true")
	})

	t.Run("Errorf", func(t *testing.T) {
		mockLogger.Errorf("error %v", "test")
		assert.Contains(t, mockLogger.Messages[3], "[ERROR] error test")
	})

	t.Run("Fatalf", func(t *testing.T) {
		exitCode := 0
		mockLogger.SetExitHandler(func(code int) { exitCode = code })

		mockLogger.Fatalf("fatal %s", "error")
		assert.Contains(t, mockLogger.Messages[4], "[FATAL] fatal error")
		assert.Equal(t, 1, exitCode)
	})
}

func TestMockLoggerLevelHandling(t *testing.T) {
	mockLogger := NewMockLogger()

	t.Run("SetLevel valid", func(t *testing.T) {
		err := mockLogger.SetLevel("info")
		assert.NoError(t, err)
		assert.Equal(t, InfoLevel, mockLogger.level)
	})

	t.Run("SetLevel invalid", func(t *testing.T) {
		err := mockLogger.SetLevel("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "неизвестный уровень логирования")
	})
}

func TestMockLoggerExitCode(t *testing.T) {
	mockLogger := NewMockLogger()
	mockLogger.SetExitHandler(func(code int) {})
	mockLogger.Fatal("test")
	assert.Equal(t, 1, mockLogger.GetExitCode(), "Код выхода должен быть 1")
}
