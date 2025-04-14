# Инструкция по запуску приложения

## Стек - Go, PostgreSQL, Docker

### 1. Клонирование репозитория

Клонирование репозитория:

```bash
git clone https://github.com/w212w/avito-pvz-service.git
```

### 2. Настройка конфигурации
- Корневая директория - avito-pvz-service.
- docker-compose.yml и Dockerfile расположен в avito-pvz-service/deployments
- База данных развертывается с помощью Docker Compose, создавать отдельный .env файл в корневой директории не обязательно. Будут использоваться параметры из config/config.go либо docker-compose.yml. Пример используемых параметров приведен ниже.


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
### 3. Запуск базы данных с Docker Compose
Для развертывания базы данных используйте Docker Compose:
```bash
docker-compose up -d db
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
- avito-shop-service/internal/service
- avito-shop-service/internal/handlers
- avito-shop-service/tests/integration


### 6. Prometheus, GRPC

- Запускаются в docker-compose
- GRPC файлы расположены в директории avito-pvz-service/internal/grpc
- Prometheus файлы расположены в директории avito-pvz-service/internal/metrics
