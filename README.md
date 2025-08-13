# Balance Service

Коротко: сервис управления балансами с REST и gRPC API. REST и gRPC поднимаются одновременно. Для REST требуется API-ключ в заголовке.

## Стек
- Go (Gin, gRPC)
- PostgreSQL
- Swagger UI (`/swagger/index.html`)

## Быстрый старт (Docker Compose)
```bash
docker compose build
docker compose up
```
- REST: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- gRPC: localhost:9090

Переменные окружения (compose выставляет автоматически):
- `DATABASE_URL=postgresql://postgres:postgres@db:5432/postgres?sslmode=disable`
- `REST_ADDR=:8080`
- `GRPC_ADDR=:9090`

Миграции применяются автоматически при старте. Сиды добавляют сервисы с тестовыми API-ключами:
- payments: `2d9a5f20-16ac-4b47-85f4-1b62b2675c8f`
- shop: `cd0fbe13-7541-4fa7-94c8-774a9f9a0e01`

## Аутентификация (REST)
Передавайте заголовок API-ключа:
- `X-API-Key: <uuid>` (или `api_key: <uuid>`)

Без валидного ключа REST запросы вернут 401 Unauthorized. Swagger не требует ключа.

## REST API (основное)
Базовый путь: `/`

- PUT `/accounts/{account_id}/limit` — изменить лимит
  - Тело: `{ "delta": 1000 }`
  - Пример:
    ```bash
    curl -X PUT 'http://localhost:8080/accounts/1/limit' \
      -H 'Content-Type: application/json' \
      -H 'X-API-Key: 2d9a5f20-16ac-4b47-85f4-1b62b2675c8f' \
      -d '{"delta": 1000}'
    ```

- PUT `/accounts/{account_id}/balance` — изменить баланс
  - Тело: `{ "delta": -500 }`
  - Пример:
    ```bash
    curl -X PUT 'http://localhost:8080/accounts/1/balance' \
      -H 'Content-Type: application/json' \
      -H 'X-API-Key: 2d9a5f20-16ac-4b47-85f4-1b62b2675c8f' \
      -d '{"delta": -500}'
    ```

- POST `/accounts/{account_id}/reservation` — открыть резерв
  - Тело: `{ "owner_service_id": 1, "amount": 1500, "idempotency_key": "k1", "timeout": "1m" }`
  - Пример:
    ```bash
    curl -X POST 'http://localhost:8080/accounts/1/reservation' \
      -H 'Content-Type: application/json' \
      -H 'X-API-Key: 2d9a5f20-16ac-4b47-85f4-1b62b2675c8f' \
      -d '{"owner_service_id":1, "amount":1500, "idempotency_key":"k1", "timeout":"1m"}'
    ```

- POST `/reservations/{reservation_id}/confirm` — подтвердить резерв
  - Заголовок: `X-Owner-Service-ID: 1`
  - Пример:
    ```bash
    curl -X POST 'http://localhost:8080/reservations/10/confirm' \
      -H 'X-API-Key: 2d9a5f20-16ac-4b47-85f4-1b62b2675c8f' \
      -H 'X-Owner-Service-ID: 1'
    ```

- POST `/reservations/{reservation_id}/cancel` — отменить резерв
  - Заголовок: `X-Owner-Service-ID: 1`
  - Пример:
    ```bash
    curl -X POST 'http://localhost:8080/reservations/10/cancel' \
      -H 'X-API-Key: 2d9a5f20-16ac-4b47-85f4-1b62b2675c8f' \
      -H 'X-Owner-Service-ID: 1'
    ```

Подробная спецификация — в Swagger UI.

## gRPC
- Адрес: `localhost:9090`
- Прото: `backend/internal/api/grpc/balance.proto`
- Пример (grpcurl):
  ```bash
  grpcurl -plaintext localhost:9090 list
  grpcurl -plaintext -d '{"account_id":1, "delta":1000}' localhost:9090 balance.BalanceService/UpdateLimit
  ```

## Локальный запуск без Docker
```bash
export DATABASE_URL='postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable'
export REST_ADDR=':8080'
export GRPC_ADDR=':9090'

# В отдельном окне поднимите PostgreSQL
# затем:
go run ./backend
```

## Структура
- `backend/main.go` — запуск REST+gRPC, миграции
- `backend/internal/api/rest` — REST-роуты и middleware
- `backend/internal/api/grpc` — gRPC сервер и proto
- `backend/internal/service` — бизнес-логика
- `backend/internal/repository` — доступ к БД (PostgreSQL)
- `backend/migrations` — миграции и сиды
- `backend/docs` — Swagger (генерируется `swag init`) 