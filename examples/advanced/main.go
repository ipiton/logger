package main

import (
	"github.com/ipiton/logger/v1"
)

func main() {
	// Пример обработки ошибок
	if err := initializeLogger(); err != nil {
		logger.Fatalf("Ошибка инициализации: %v", err)
	}
}

func initializeLogger() error {
	logger.New().
		WithFile("app.log").
		WithFields(map[string]interface{}{"version": "1.0"})
	return nil
}
