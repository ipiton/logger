package main

import (
	"os"
	"path/filepath"

	"github.com/ipiton/logger"
)

func main() {
	// Создаем директорию для логов
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	// Создаем конфигурацию с файлами логов
	cfg := logger.DefaultConfig()
	cfg.Files = map[string]string{
		"main":  filepath.Join(logDir, "app.log"),
		"error": filepath.Join(logDir, "error.log"),
	}

	// Создаем и настраиваем логгер
	log := logger.New().
		WithFile(cfg.Files["main"]).
		WithLevel(cfg.Level)
	defer log.Close()

	// Логируем в файл
	log.Info("Приложение запущено")
	log.WithFields(map[string]interface{}{
		"module": "api",
		"method": "POST",
	}).Info("Обработка запроса")

	// Логируем ошибку
	log.Error("Пример ошибки")
}
