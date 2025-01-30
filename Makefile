.PHONY: build test coverage lint clean examples

# Переменные сборки
BUILD_DIR = build
COVERAGE_FILE = coverage.out
COVERAGE_HTML = coverage.html
VERSION = $(shell git describe --tags --always --dirty)
COMMIT_HASH = $(shell git rev-parse --short HEAD)
BUILD_DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS = -X 'github.com/ipiton/logger.Version=$(VERSION)' \
          -X 'github.com/ipiton/logger.CommitHash=$(COMMIT_HASH)' \
          -X 'github.com/ipiton/logger.BuildDate=$(BUILD_DATE)' \
          -buildvcs=false

# Сборка
build:
	@echo "Сборка..."
	@go build -ldflags "$(LDFLAGS)" ./...

# Запуск тестов
test:
	@echo "Запуск тестов..."
	@go test -v -race --timeout 10s ./...

# Проверка покрытия кода тестами
coverage:
	@echo "Проверка покрытия..."
	@go test -v -race -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@go tool cover -func=$(COVERAGE_FILE)

# Запуск линтера
lint:
	@echo "Проверка кода линтером..."
	@golangci-lint run -c .golangci.yml

# Сборка примеров
examples:
	@echo "Сборка примеров..."
	@for dir in examples/*; do \
		if [ -d "$$dir" ]; then \
			echo "Сборка $$dir..."; \
			go build -o $(BUILD_DIR)/$$(basename $$dir) $$dir/main.go; \
		fi \
	done

# Очистка артефактов сборки
clean:
	@echo "Очистка..."
	@rm -rf $(BUILD_DIR)
	@rm -f $(COVERAGE_FILE)
	@rm -f $(COVERAGE_HTML)
	@find . -type f -name '*.log' -delete
	@rm -rf logs/

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  make build     - сборка проекта"
	@echo "  make test      - запуск тестов"
	@echo "  make coverage  - проверка покрытия кода тестами"
	@echo "  make lint      - проверка кода линтером"
	@echo "  make examples  - сборка примеров"
	@echo "  make clean     - очистка артефактов сборки"
