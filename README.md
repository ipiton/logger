# Go Logger

[![Go Reference](https://pkg.go.dev/badge/github.com/ipiton/logger.svg)](https://pkg.go.dev/github.com/ipiton/logger)
[![Go Report Card](https://goreportcard.com/badge/github.com/ipiton/logger)](https://goreportcard.com/report/github.com/ipiton/logger)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Coverage Status](https://coveralls.io/repos/github/ipiton/logger/badge.svg?branch=main)](https://coveralls.io/github/ipiton/logger?branch=main)

Потокобезопасный структурированный логгер для Go с поддержкой:
- Уровней логирования (DEBUG, INFO, WARNING, ERROR, FATAL)
- Форматированного вывода
- Записи в файлы
- Контекстных полей и префиксов
- Кастомизации через конфигурацию

## Обзор

Пакет `logger` предоставляет потокобезопасную, гибкую и расширяемую систему логирования для Go приложений. Он поддерживает различные уровни логирования, форматирование сообщений, настройку префиксов, добавление полей и запись в файлы.

## Требования

- Go 1.18 или выше

## Возможности

- Несколько уровней логирования: DEBUG, INFO, WARNING, ERROR, FATAL (с завершением программы)
- Форматирование сообщений (Debugf, Infof, Warningf, Errorf, Fatalf)
- Настройка префиксов для сообщений
- Добавление дополнительных полей к сообщениям
- Возможность логирования в файл
- Потокобезопасность
- Глобальный логгер
- Поддержка цепочки вызовов (chaining)

## Установка

```bash
go get github.com/ipiton/logger
```

## Примеры использования

### Базовое использование
```go
package main

import "github.com/ipiton/logger"

func main() {
    log := logger.New()
    log.Info("Starting application")
    log.WithFields(map[string]interface{}{"version": "1.0.0"}).Debug("Debug info")
}
```

### Создание собственного логгера

```go
// Создание логгера с настройками
customLogger := logger.New().
    WithPrefix("SERVICE").
    WithLevel(logger.InfoLevel).
    WithFile("/var/log/service.log")

// Логирование с префиксом
serviceLogger := customLogger.WithPrefix("SERVICE")
serviceLogger.Info("Сервис запущен")
serviceLogger.Infof("Порт: %d", 8080)

// Логирование с дополнительными полями
userLogger := customLogger.WithFields(map[string]interface{}{
    "user_id": 123,
    "action": "login",
})
userLogger.Info("Пользователь авторизован")
userLogger.Infof("Время сессии: %v", duration)

// Цепочка вызовов
customLogger.WithPrefix("API").
    WithFields(map[string]interface{}{"method": "POST"}).
    Infof("Запрос обработан за %dms", 42)
```

### Логирование в файл

```go
// Создание логгера с записью в файл
fileLogger := customLogger.WithFile("/var/log/myapp.log")
fileLogger.Info("Сообщение будет записано в файл")
fileLogger.Infof("Статистика: обработано %d запросов", count)
```

### Глобальный логгер

```go
// Установка глобального логгера
globalLogger := logger.New()
logger.SetGlobalLogger(globalLogger)

// Использование глобального логгера
logger.Info("Сообщение через глобальный логгер")
logger.Infof("Версия: %s", version)
```

### Конфигурация по умолчанию

```go
// Значения по умолчанию
defaultConfig := logger.NewConfig()
// Level: INFO
// Prefix: ""
// TimeFormat: "2006-01-02 15:04:05"
// Fields: nil
// File: stdout
```

### Thread Safety

Все операции логгера потокобезопасны и могут использоваться из разных горутин без дополнительной синхронизации.

## Уровни логирования

- `DEBUG`: Детальная отладочная информация
- `INFO`: Информационные сообщения
- `WARNING`: Предупреждения
- `ERROR`: Ошибки
- `FATAL`: Критические ошибки, после которых приложение завершает работу

## Настройка

### Конфигурация логгера

```go
config := logger.NewConfig().
    WithLevel(logger.DEBUG).
    WithPrefix("APP").
    WithFile("/var/log/app.log").
    WithTime(true).
    WithFields(map[string]interface{}{
        "version": "1.0.0",
        "env": "production",
    })
```

## Интерфейс ILogger

Пакет предоставляет интерфейс `ILogger`, который включает следующие методы:

```go
type ILogger interface {
    Debug(args ...interface{})
    Debugf(format string, args ...interface{})
    Info(args ...interface{})
    Infof(format string, args ...interface{})
    Warning(args ...interface{})
    Warningf(format string, args ...interface{})
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
    Fatal(args ...interface{})
    Fatalf(format string, args ...interface{})
    WithPrefix(prefix string) ILogger
    WithFields(fields map[string]interface{}) ILogger
    WithFile(file string) ILogger
}
```

## Обработка ошибок

Пакет предоставляет специализированные типы ошибок:

- `ConfigError`: Ошибки конфигурации логгера
- `WriteError`: Ошибки записи лога
- `LevelError`: Ошибки неверного уровня логирования

## Тестирование

Пакет полностью покрыт unit-тестами с использованием `testify`. Включает тесты:
- Базового функционала
- Форматированного вывода
- Конкурентного доступа
- Работы с файлами
- Обработки ошибок
- Глобального логгера

## Производительность

Логгер оптимизирован для минимального overhead и потокобезопасен. Все операции логирования выполняются асинхронно, чтобы не блокировать основной поток выполнения.

## Лицензия

Этот проект распространяется под лицензией [MIT](LICENSE).

## Contributing

Мы приветствуем вклад в проект! Пожалуйста, ознакомьтесь с руководством перед началом работы.

### Процесс внесения изменений

1. Создайте форк репозитория
2. Создайте feature-ветку:
   `git checkout -b feature/AmazingFeature`
3. Установите зависимости:
   `go mod tidy`
4. Запустите тесты и линтеры:
   `make test lint`
5. Закоммитьте изменения:
   `git commit -m 'feat: add AmazingFeature'`
   (используйте [Conventional Commits](https://www.conventionalcommits.org/))
6. Запушьте ветку:
   `git push origin feature/AmazingFeature`
7. Откройте Pull Request

### Требования к коду

- 100% покрытие тестами для нового функционала
- Соответствие Go Code Review Comments
- Документация для публичных методов
- Отсутствие предупреждений линтеров

### Шаблоны Issues

При создании issue укажите:
- Версию Go: `go version`
- Версию пакета: `git rev-parse HEAD`
- Шаги для воспроизведения
- Ожидаемое и фактическое поведение

### Code of Conduct

Проект следует [Contributor Covenant](CODE_OF_CONDUCT.md). Участвуя, вы соглашаетесь соблюдать его правила.

## Примеры обработки ошибок

```go
config := logger.DefaultConfig().
    WithLevel("debug").
    WithFile("/var/log/app.log").
    WithTimeFormat(time.RFC3339)

log := logger.New().WithConfig(config)
```

## Документация

Полная документация доступна на [pkg.go.dev](https://pkg.go.dev/github.com/ipiton/logger).

Основные методы:
- `Debug/Debugf` - отладочные сообщения
- `Info/Infof` - информационные сообщения
- `Warning/Warningf` - предупреждения
- `Error/Errorf` - ошибки приложения
- `Fatal/Fatalf` - критические ошибки с завершением программы
- `WithPrefix` - создание дочернего логгера с префиксом
- `WithFields` - добавление контекстных полей
- `WithFile` - запись логов в файл
- `Sync` - принудительная синхронизация буферов
```go
</rewritten_file>
