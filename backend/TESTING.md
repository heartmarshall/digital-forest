# Testing Guide

Этот документ описывает систему тестирования для Digital Forest Backend.

## Структура тестов

### Unit Tests
- **Расположение**: `./internal/*/..._test.go`
- **Покрытие**: Use cases, handlers, DTOs
- **Зависимости**: Моки (testutil/mocks.go)

### Integration Tests
- **Расположение**: `./internal/repository/.../..._test.go`, `./internal/transport/http/integration_test.go`
- **Покрытие**: Репозитории, HTTP API
- **Зависимости**: Testcontainers (PostgreSQL)

### End-to-End Tests
- **Расположение**: `./e2e_test.go`
- **Покрытие**: Полный workflow приложения
- **Зависимости**: Testcontainers, реальная база данных

## Запуск тестов

### Все тесты
```bash
make test
```

### Unit тесты
```bash
make test-unit
```

### Интеграционные тесты
```bash
make test-integration
```

### E2E тесты
```bash
make test-e2e
```

### Тесты с покрытием
```bash
make test-coverage
```

### Тесты в Docker
```bash
make docker-test
```

## Требования

### Локальная разработка
- Go 1.21+
- Docker (для testcontainers)
- Make

### CI/CD
- Docker
- Docker Compose

## Контейнеризация

### Testcontainers
Используется для создания изолированных тестовых окружений:
- PostgreSQL контейнер для интеграционных тестов
- Автоматическая очистка после тестов
- Изоляция тестов друг от друга

### Docker Compose для тестов
```yaml
# docker-compose.test.yml
services:
  test-postgres:
    image: postgres:15-alpine
    # ... конфигурация
```

## Покрытие кода

### Генерация отчета
```bash
make test-coverage
```

### Просмотр отчета
Откройте `coverage.html` в браузере.

### Целевое покрытие
- Unit тесты: >90%
- Интеграционные тесты: >80%
- E2E тесты: основные сценарии

## Тестовые данные

### Тестовые растения
```go
// internal/testutil/testdata.go
var TestPlants = []domain.Plant{
    // ... тестовые данные
}
```

### Моки
```go
// internal/testutil/mocks.go
type MockPlantRepository struct {
    mock.Mock
}
```

## Лучшие практики

### 1. Изоляция тестов
- Каждый тест должен быть независимым
- Используйте `testutil.TruncateTables()` для очистки БД
- Моки для внешних зависимостей

### 2. Именование
- `TestFunctionName_Scenario_ExpectedResult`
- Группировка тестов по функциональности

### 3. Структура теста
```go
func TestFunction(t *testing.T) {
    // Arrange
    // Act
    // Assert
}
```

### 4. Тестовые данные
- Используйте константы для повторяющихся данных
- Создавайте фабрики для сложных объектов
- Избегайте хардкода

## Отладка тестов

### Запуск конкретного теста
```bash
go test -v -run TestSpecificFunction ./path/to/package
```

### Запуск с детальным выводом
```bash
go test -v -race ./...
```

### Пропуск медленных тестов
```bash
go test -v -short ./...
```

## CI/CD интеграция

### GitHub Actions
```yaml
- name: Run tests
  run: make ci
```

### Локальная проверка
```bash
make ci
```

## Troubleshooting

### Проблемы с testcontainers
1. Убедитесь, что Docker запущен
2. Проверьте доступность портов
3. Очистите неиспользуемые контейнеры: `docker system prune`

### Проблемы с базой данных
1. Проверьте подключение к PostgreSQL
2. Убедитесь, что миграции применены
3. Проверьте права доступа

### Медленные тесты
1. Используйте `-short` флаг для unit тестов
2. Оптимизируйте тестовые данные
3. Используйте параллельное выполнение где возможно

## Расширение тестов

### Добавление нового теста
1. Создайте файл `*_test.go` в соответствующем пакете
2. Используйте существующие утилиты из `testutil`
3. Добавьте тест в соответствующий Makefile target

### Добавление нового мока
1. Создайте структуру в `testutil/mocks.go`
2. Реализуйте необходимые методы
3. Добавьте проверку интерфейса

### Добавление тестовых данных
1. Добавьте данные в `testutil/testdata.go`
2. Создайте фабричные функции при необходимости
3. Документируйте назначение данных
