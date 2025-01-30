#!/bin/bash

# Запуск тестов с проверкой гонок и генерацией отчета о покрытии
go test -v -race -coverprofile=coverage.out ./...

# Генерация HTML-отчета
go tool cover -html=coverage.out -o coverage.html

# Вывод процента покрытия
go tool cover -func=coverage.out
