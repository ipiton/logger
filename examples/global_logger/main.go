package main

import (
	"time"

	"github.com/ipiton/logger"
)

// Пример компонента приложения
type UserService struct {
	log logger.ILogger
}

func NewUserService() *UserService {
	return &UserService{
		log: logger.New().WithConfig(logger.DefaultConfig()).WithPrefix("UserService"),
	}
}

func (s *UserService) CreateUser(username string) error {
	s.log.WithFields(map[string]interface{}{
		"username": username,
		"time":     time.Now(),
	}).Info("Создание нового пользователя")

	// Имитация работы
	time.Sleep(100 * time.Millisecond)

	s.log.Infof("Пользователь %s успешно создан", username)
	return nil
}

func main() {
	// Настраиваем глобальный логгер
	cfg := logger.DefaultConfig()
	cfg.Level = "debug"

	globalLogger := logger.New().WithConfig(cfg)
	defer globalLogger.Close()

	// Используем глобальный логгер напрямую
	logger.Info("Запуск приложения")

	// Создаем сервис
	userService := NewUserService()

	// Выполняем операции
	if err := userService.CreateUser("john_doe"); err != nil {
		logger.WithFields(map[string]interface{}{
			"error": err,
			"user":  "john_doe",
		}).Error("Ошибка создания пользователя")
	}

	logger.Info("Завершение работы приложения")
}
