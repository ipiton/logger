package main

import (
	"github.com/ipiton/logger"
)

func main() {
	// Создаем конфигурацию
	cfg := logger.DefaultConfig()
	cfg.Level = "debug"

	// Создаем логгер
	log := logger.New()
	defer log.Close()

	// Базовое логирование
	log.Debug("Это отладочное сообщение")
	log.Info("Это информационное сообщение")

	// Логирование с префиксом
	serviceLog := log.WithPrefix("SERVICE")
	serviceLog.Info("Сервис запущен")

	// Логирование с полями
	serviceLog.WithFields(map[string]interface{}{
		"port": 8080,
		"host": "localhost",
	}).Info("Сервер слушает")

	// Новый пример использования
	logger.New().WithPrefix("APP").Info("message")
}
