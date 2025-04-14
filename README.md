# Инструкция по запуску приложения

## Стек - Go, PostgreSQL, Docker, GRPC, Prometheus

### 1. Клонирование репозитория

Клонирование репозитория:

```bash
git clone https://github.com/w212w/avito-pvz-service.git
```

### 2. Настройка конфигурации
- Корневая директория - avito-pvz-service
- docker-compose.yml и Dockerfile расположен в avito-pvz-service/deployments
- База данных развертывается с помощью Docker Compose, создавать отдельный .env файл в корневой директории не обязательно. Будут использоваться параметры из config/config.go либо docker-compose.yml. Пример используемых параметров приведен ниже


```bash
# .env файл (локальное развертывание)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=shop
JWT_SECRET=supersecretkey
LOG_LEVEL=debug
```
```bash
# (Развертывание через docker, параметры из docker-compose)
APP_ENV=docker
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=shop
JWT_SECRET=supersecretkey
LOG_LEVEL=debug
```
### 3. Развертывание с Docker Compose
Для развертывания сервиса, базы данных, prometheus в контейнерах используйте Docker Compose:
```bash
docker-compose up -d
```
### 4. Применение миграций
При запуске приложения миграции для создания таблиц в базе данных будут применены автоматически. 
- Миграции расположены в avito-pvz-service/internal/storage/migrations

### 5. Запуск тестов
```bash
go test ./...
go test -cover  ./...
```
**Результаты тестов:** <br>
ok      avito-pvz-service/internal/handlers     1.060s  coverage: 79.9% of statements

ok      avito-pvz-service/internal/services     1.543s  coverage: 79.5% of statements

ok      avito-pvz-service/tests/integration     1.113s  coverage: [no statements]<br>

Тесты расположены в следующих директориях:
- avito-pvz-service/internal/service
- avito-pvz-service/internal/handlers
- avito-pvz-service/tests/integration


### 6. Prometheus, GRPC, логирование, кодогенерация DTO endpoint'ов

- Запускаются в docker-compose
- GRPC файлы расположены в директории avito-pvz-service/internal/grpc
- Prometheus файлы расположены в директории avito-pvz-service/internal/metrics
- Логирование осуществляется посредством "logrus" в директории avito-pvz-service/pkg
- Кодогенерация DTO endpoint'ов расположена в директории avito-pvz-service/api/gen

### 7. Вопросы и их решения

- Какой должен быть формат хранения времени ? Испоьзовал UTC, чтобы не иметь зависимостей от локации. Предположил, что теоретически с сервисом могут работать люди из разных городов, либо же ПВЗ будут открывать позже в городах с другим часовым поясом и, вероятно, лучше использовать UTC.