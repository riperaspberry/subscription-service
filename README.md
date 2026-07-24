# Subscription Service

REST API для хранения подписок пользователей и расчёта суммарных трат за период.

## Stack

- Go
- Gin
- PostgreSQL
- pgx
- Docker

## Требования

- Go 1.25+
- Docker и Docker Compose (для PostgreSQL)

## Быстрый старт

### 1. Клонировать репозиторий

```bash
git clone <repo-url>
cd subscription-service
```

### 2. Настроить переменные окружения

Создай файл `.env` в корне проекта (или скопируй из примера):

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions
DB_SSLMODE=disable
```

### 3. Запустить PostgreSQL

```bash
docker compose up -d
```

База будет доступна на `localhost:5433`.

### 4. Применить миграции

```bash
docker exec -i subscriptions-db psql -U postgres -d subscriptions < migrations/000001_create_subscriptions_table.up.sql
```

### 5. Запустить сервис

```bash
go run ./cmd/app
```

Сервис стартует на `http://localhost:8080`.

### Проверка сборки и тестов

```bash
go build ./...
go test ./...
```

## API

Базовый URL: `http://localhost:8080`

Формат дат во всех запросах: `MM-YYYY` (например, `07-2026`).

### Создать подписку

`POST /subscriptions`

**Body (JSON):**

```json
{
  "service_name": "Netflix",
  "price": 999,
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "start_date": "01-2026",
  "end_date": "12-2026"
}
```

| Поле | Тип | Обязательное | Описание |
|------|-----|--------------|----------|
| `service_name` | string | да | Название сервиса |
| `price` | int | да | Цена за месяц |
| `user_id` | string (UUID) | да | ID пользователя |
| `start_date` | string | да | Дата начала (`MM-YYYY`) |
| `end_date` | string | нет | Дата окончания (`MM-YYYY`) |

**Ответ:** `201 Created` — объект подписки.

---

### Получить список подписок

`GET /subscriptions`

**Ответ:** `200 OK` — массив подписок.

---

### Получить подписку по ID

`GET /subscriptions/{id}`

**Ответ:** `200 OK` — объект подписки.

---

### Обновить подписку

`PUT /subscriptions/{id}`

**Body (JSON):**

```json
{
  "service_name": "Netflix Premium",
  "price": 1299,
  "start_date": "03-2026",
  "end_date": "12-2026"
}
```

**Ответ:** `200 OK`

---

### Удалить подписку

`DELETE /subscriptions/{id}`

**Ответ:** `204 No Content`

---

### Рассчитать траты за период

`GET /subscriptions/calculate`

Считает сумму трат пользователя по конкретному сервису за указанный период.  
Если `to` позже текущего месяца, период обрезается до текущего месяца.

**Query-параметры:**

| Параметр | Тип | Обязательный | Описание |
|----------|-----|--------------|----------|
| `user_id` | string (UUID) | да | ID пользователя |
| `service_name` | string | да | Название сервиса |
| `from` | string | да | Начало периода (`MM-YYYY`) |
| `to` | string | да | Конец периода (`MM-YYYY`) |

**Пример:**

```bash
curl "http://localhost:8080/subscriptions/calculate?user_id=550e8400-e29b-41d4-a716-446655440000&service_name=Netflix&from=01-2026&to=06-2026"
```

**Ответ:** `200 OK`

```json
{
  "total": 5994
}
```

## Примеры curl

```bash
# Создать подписку
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Spotify",
    "price": 299,
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "start_date": "01-2026"
  }'

# Получить все подписки
curl http://localhost:8080/subscriptions

# Получить подписку по ID
curl http://localhost:8080/subscriptions/{id}

# Обновить подписку
curl -X PUT http://localhost:8080/subscriptions/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "service_name": "Spotify Family",
    "price": 399
  }'

# Удалить подписку
curl -X DELETE http://localhost:8080/subscriptions/{id}
```

## Структура проекта

```
cmd/app/              — точка входа
internal/
  config/             — конфигурация из .env
  database/           — подключение к PostgreSQL
  handler/            — HTTP-обработчики
  model/              — модели и DTO
  repository/         — работа с БД
  router/             — маршруты
  service/            — бизнес-логика
migrations/           — SQL-миграции
```
