package logger

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoggerScenarios(t *testing.T) {
	tempDir := t.TempDir()
	mainLogFile := filepath.Join(tempDir, "main.log")
	require.NoError(t, os.MkdirAll(filepath.Dir(mainLogFile), 0750))

	// Создаем базовый логгер
	logger := New().WithFile(mainLogFile).WithLevel("debug")
	defer func() {
		if err := logger.Close(); err != nil {
			t.Errorf("Failed to close logger: %v", err)
		}
	}()

	// Создаем дочерние логгеры
	prefixLogger := logger.WithPrefix("PREFIX")
	debugLogger := logger.WithLevel("debug")

	t.Run("базовое логирование", func(t *testing.T) {
		logger.Info("test info message")
		logger.WithFields(map[string]interface{}{"key": "value"}).Debug("test debug message")

		content := readLogFileSecure(t, mainLogFile)
		assert.Contains(t, content, "[INFO] test info message")
		assert.Contains(t, content, "[DEBUG] test debug message [key=value]")
	})

	t.Run("логирование с префиксом", func(t *testing.T) {
		prefixLogger.Info("test prefix message")
		prefixLogger.WithFields(map[string]interface{}{"id": 1}).Debug("test prefix debug")

		content := readLogFileSecure(t, mainLogFile)
		assert.Contains(t, content, "[INFO] [PREFIX] test prefix message")
		assert.Contains(t, content, "[DEBUG] [PREFIX] test prefix debug [id=1]")
	})

	t.Run("обработка ошибок", func(t *testing.T) {
		logger.Error("test error message")
		logger.WithFields(map[string]interface{}{"code": 1}).Error("test error with code")

		content := readLogFileSecure(t, mainLogFile)
		assert.Contains(t, content, "[ERROR] test error message")
		assert.Contains(t, content, "[ERROR] test error with code [code=1]")
	})

	t.Run("уровни логирования", func(t *testing.T) {
		debugLogger.Debug("test debug level")
		debugLogger.Info("test info level")

		content := readLogFileSecure(t, mainLogFile)
		assert.Contains(t, content, "[DEBUG] test debug level")
		assert.Contains(t, content, "[INFO] test info level")
	})

	t.Run("параллельная запись", func(t *testing.T) {
		var wg sync.WaitGroup
		iterations := 100

		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				logger.Infof("test message %d", i)
				prefixLogger.Debugf("test debug %d", i)
			}(i)
		}

		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			// Даем время на запись всех сообщений
			time.Sleep(100 * time.Millisecond)
			content := readLogFileSecure(t, mainLogFile)
			lines := strings.Split(strings.TrimSpace(content), "\n")
			assert.GreaterOrEqual(t, len(lines), iterations*2,
				"Ожидалось как минимум %d сообщений, получено %d", iterations*2, len(lines))
		case <-time.After(10 * time.Second):
			t.Fatal("Тест превысил таймаут")
		}
	})
}

func TestWriteToFile(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	logger := New().WithFile(logFile)
	defer func() {
		if err := logger.Close(); err != nil {
			t.Errorf("Failed to close logger: %v", err)
		}
	}()

	logger.Info("test message")

	assert.Contains(t, readLogFileSecure(t, logFile), "[INFO] test message")
}

func TestLoggerConfiguration(t *testing.T) {
	tempDir := t.TempDir()
	mainLogPath := filepath.Join(tempDir, "main.log")
	require.NoError(t, os.MkdirAll(filepath.Dir(mainLogPath), 0750))

	tests := []struct {
		name     string
		level    string
		messages []string
	}{
		{
			name:  "debug_level_configuration",
			level: "debug",
			messages: []string{
				"DEBUG] debug message",
				"INFO] info message",
				"ERROR] error message",
			},
		},
		{
			name:  "info level configuration",
			level: "info",
			messages: []string{
				"INFO] info message",
			},
		},
		{
			name:  "error level configuration",
			level: "error",
			messages: []string{
				"ERROR] error message",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New().WithLevel(tt.level)

			logger.Debug("debug message")
			logger.Info("info message")
			logger.Error("error message")

			messages := logger.GetMessages()
			for _, expectedMsg := range tt.messages {
				// Проверяем, что хотя бы одно сообщение содержит ожидаемый текст
				found := false
				for _, msg := range messages {
					if strings.Contains(msg, expectedMsg) {
						found = true
						break
					}
				}
				assert.True(t, found, "Сообщение не найдено: %s", expectedMsg)
			}
		})
	}
}

func TestLoggerFormatMethods(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	logger := New().WithFile(logFile)
	defer func() {
		if f := logger.(*Logger).file; f != nil {
			if err := f.Sync(); err != nil {
				t.Errorf("Failed to sync file: %v", err)
			}
		}
		if err := logger.Close(); err != nil {
			t.Errorf("Failed to close logger: %v", err)
		}
	}()

	// Перехват os.Exit только для теста Fatalf
	originalOsExit := osExit
	defer func() { osExit = originalOsExit }()
	exited := false
	osExit = func(code int) { exited = true }

	tests := []struct {
		name     string
		logFunc  func(l ILogger)
		expected string
		isFatal  bool // Добавляем флаг для fatal-тестов
	}{
		{
			name: "Debugf",
			logFunc: func(l ILogger) {
				l.Debugf("debug message %d", 1)
			},
			expected: `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[DEBUG\] debug message 1$`,
			isFatal:  false,
		},
		{
			name: "Infof",
			logFunc: func(l ILogger) {
				l.Infof("info message %s", "test")
			},
			expected: `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[INFO\] info message test$`,
			isFatal:  false,
		},
		{
			name: "Warningf",
			logFunc: func(l ILogger) {
				l.Warningf("warning message %t", true)
			},
			expected: `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[WARNING\] warning message true$`,
			isFatal:  false,
		},
		{
			name: "Errorf",
			logFunc: func(l ILogger) {
				l.Errorf("error message %q", "test")
			},
			expected: `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[ERROR\] error message "test"$`,
			isFatal:  false,
		},
		{
			name: "Fatalf",
			logFunc: func(l ILogger) {
				l.Fatalf("fatal message %d", 1)
			},
			expected: `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[FATAL\] fatal message 1$`,
			isFatal:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isFatal {
				exited = false // Сбрасываем флаг перед каждым fatal-тестом
			}

			// Очищаем файл перед каждым тестом
			err := os.WriteFile(logFile, []byte{}, 0600)
			require.NoError(t, err)

			// Вызываем тестируемый метод
			tt.logFunc(logger)

			if tt.isFatal {
				// Проверяем, что osExit был вызван
				assert.True(t, exited, "osExit должен быть вызван для FATAL уровня")
			}

			// Проверяем содержимое файла
			content := readLogFileSecure(t, logFile)
			logLine := strings.TrimSpace(content)

			// Проверяем формат лога с помощью регулярного выражения
			matched, err := regexp.MatchString(tt.expected, logLine)
			require.NoError(t, err)
			assert.True(t, matched, "Строка лога не соответствует ожидаемому формату")
		})
	}
}

func TestLoggerWithPrefixAndFields(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	tests := []struct {
		name     string
		setup    func(logger ILogger) ILogger
		messages []string
	}{
		{
			name: "WithPrefix and Debugf",
			setup: func(logger ILogger) ILogger {
				return logger.WithPrefix("TEST")
			},
			messages: []string{
				"[DEBUG] [TEST] message 1",
			},
		},
		{
			name: "WithFields and Infof",
			setup: func(logger ILogger) ILogger {
				return logger.WithFields(map[string]interface{}{"key": "value"})
			},
			messages: []string{
				"[INFO] message test [key=value]",
			},
		},
		{
			name: "WithPrefix, WithFields and Warningf",
			setup: func(logger ILogger) ILogger {
				return logger.WithPrefix("TEST").WithFields(map[string]interface{}{"key": "value"})
			},
			messages: []string{
				"[WARNING] [TEST] message true [key=value]",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New().WithFile(logFile)
			defer func() {
				if err := logger.Close(); err != nil {
					t.Errorf("Failed to close logger: %v", err)
				}
			}()

			configuredLogger := tt.setup(logger)

			// Тестируем только один метод для каждого случая
			switch tt.name {
			case "WithPrefix and Debugf":
				configuredLogger.Debugf("message %d", 1)
			case "WithFields and Infof":
				configuredLogger.Infof("message %s", "test")
			case "WithPrefix, WithFields and Warningf":
				configuredLogger.Warningf("message %t", true)
			}

			// Получаем сообщения из настроенного логгера
			messages := configuredLogger.GetMessages()

			for _, expected := range tt.messages {
				found := false
				for _, msg := range messages {
					if strings.Contains(msg, expected) {
						found = true
						break
					}
				}
				assert.True(t, found, "Сообщение должно содержать: %s", expected)
			}
		})
	}
}

func TestLoggerFatalMethod(t *testing.T) {
	tempDir := t.TempDir()
	mainLogPath := filepath.Join(tempDir, "main.log")
	require.NoError(t, os.MkdirAll(filepath.Dir(mainLogPath), 0750))

	// Создаем файлы для логов
	mainLogPath = filepath.Join(tempDir, "main.log")

	// Создаем конфигурацию
	config := &Config{
		Level: "debug",
		Files: map[string]string{
			"main": mainLogPath,
		},
	}

	// Создаем логгер с конфигурацией
	logger := New().WithConfig(config)
	logger = logger.WithFile(mainLogPath)
	defer func() {
		if err := logger.Close(); err != nil {
			t.Errorf("Failed to close logger: %v", err)
		}
	}()

	// Перехватываем osExit
	originalOsExit := osExit
	defer func() { osExit = originalOsExit }()
	exited := false
	osExit = func(code int) { exited = true }

	// Вызываем Fatal
	logger.Fatal("test fatal message")

	// Синхронизируем запись в файл
	if f := logger.(*Logger).file; f != nil {
		if err := f.Sync(); err != nil {
			t.Errorf("failed to sync file: %v", err)
		}
	}

	// Проверяем, что osExit был вызван
	assert.True(t, exited, "osExit должен быть вызван")

	// Проверяем содержимое файла
	content := readLogFileSecure(t, mainLogPath)
	logLine := strings.TrimSpace(content)

	const pattern = `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[FATAL\] test fatal message`

	t.Logf("Actual log line: %s", logLine)
	t.Logf("Expected pattern: %s", pattern)

	// Разделяем проверку на части
	assert.Contains(t, logLine, "[FATAL] test fatal message", "Отсутствует FATAL сообщение")
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}`, logLine, "Неверный формат времени")
}

func TestLoggerWithFile(t *testing.T) {
	tempDir := t.TempDir()
	logFile := filepath.Join(tempDir, "test.log")

	logger := New().WithLevel("debug").WithFile(logFile)
	defer func() {
		if err := logger.Close(); err != nil {
			t.Errorf("Failed to close logger: %v", err)
		}
	}()
	logger.Info("test message")

	// Синхронизируем запись в файл
	if f := logger.(*Logger).file; f != nil {
		if err := f.Sync(); err != nil {
			t.Errorf("failed to sync file: %v", err)
		}
	}

	// Проверяем с учетом временной метки
	assert.Contains(t, readLogFileSecure(t, logFile), "[INFO] test message")
}

func TestTimeFormatting(t *testing.T) {
	logger := New().WithTimeFormat(time.RFC3339).WithLevel("debug")
	logger.Info("Тест времени")

	messages := logger.GetMessages()
	require.Greater(t, len(messages), 0, "Нет записанных сообщений")

	logEntry := messages[0]
	parts := strings.Split(logEntry, " ")
	require.Greater(t, len(parts), 0, "Неверный формат сообщения")

	_, err := time.Parse(time.RFC3339, parts[0])
	assert.NoError(t, err, "Неверный формат времени")
}

// readLogFileSecure безопасно читает файл лога
func readLogFileSecure(t *testing.T, filename string) string {
	// #nosec G304
	content, err := os.ReadFile(filename)
	require.NoError(t, err)
	return string(content)
}
