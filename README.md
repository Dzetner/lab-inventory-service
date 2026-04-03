# lab-inventory-service

Простой REST-сервис для управления лабораторным инвентарём: сотрудники, помещения, химические вещества и конкретные ёмкости (бутылки/флаконы) с реагентами.

## Стек

- Go + chi (HTTP роутер)
- PostgreSQL
- sqlc (генерация type-safe слоя доступа к БД)
- pgx / pgxpool

## Модель данных

Основные сущности:

- **employees** – сотрудники лаборатории  
  - `id`, `full_name`, `role`, `created_at`
- **rooms** – помещения / комнаты  
  - `id`, `name`, `description`, `created_at`
- **chemicals** – химические вещества  
  - `id`, `name`, `cas_number`, `formula`, `sds_url`, `created_at`
- **containers** – конкретные ёмкости с реагентами  
  - `id`, `chemical_id`, `room_id`, `label_code`, `quantity`, `unit`, `status`, `checked_out_by`, `created_at`

Одна запись в `chemicals` может иметь много `containers`.

## Запуск

### 1. PostgreSQL

Создать базу:

```sql
CREATE DATABASE lab_inventory;
```

Выполнить SQL-схему (из `sql/schema.sql`) в базе `lab_inventory` – можно через psql или pgAdmin.

По умолчанию сервис подключается к Postgres по DSN:

```text
postgres://postgres:postgres@localhost:5432/lab_inventory?sslmode=disable
```

Можно переопределить через переменную окружения:

```bash
export LAB_INVENTORY_DB_DSN="postgres://USER:PASS@HOST:PORT/lab_inventory?sslmode=disable"
```

### 2. Генерация кода sqlc

```bash
sqlc generate
```

### 3. Запуск сервиса

```bash
go run ./cmd/lab-inventory-service
```

По умолчанию сервис слушает `http://localhost:8080`.

## HTTP API (черновик)

### Health

```http
GET /health
```

Ответ: `200 OK`, текст `ok`.

---

### Employees

```http
GET /employees
```

Возвращает список сотрудников.

```http
POST /employees
Content-Type: application/json

{
  "full_name": "Ivan Petrov",
  "role": "chemist"
}
```

Ответ: `201 Created` и созданный сотрудник.

---

### Rooms

```http
GET /rooms
```

```http
POST /rooms
Content-Type: application/json

{
  "name": "Room 101",
  "description": "Organic synthesis lab"
}
```

---

### Chemicals

```http
GET /chemicals
```

```http
POST /chemicals
Content-Type: application/json

{
  "name": "Acetone",
  "cas_number": "67-64-1",
  "formula": "C3H6O",
  "sds_url": "https://example.com/acetone-sds"
}
```

#### Поиск химикатов

```http
GET /chemicals/search?query=acetone
```

Ищет по имени, CAS и формуле (ILIKE).

---

### Containers

```http
GET /containers
```

Список всех ёмкостей.

#### Фильтрация контейнеров

```http
GET /containers?status=available&room_id=1
```

`status` – строка (`available`, `checked_out` и т.п.),  
`room_id` – id помещения.

#### Создание контейнера

```http
POST /containers
Content-Type: application/json

{
  "chemical_id": 1,
  "room_id": 1,
  "label_code": "ACET-001",
  "quantity": 0.5,
  "unit": "L",
  "status": "available"
}
```

#### Взять контейнер (checkout)

```http
POST /containers/{id}/checkout
Content-Type: application/json

{
  "employee_id": 1
}
```

Меняет `status` на `checked_out` и ставит `checked_out_by`.

#### Вернуть контейнер

```http
POST /containers/{id}/return
```

Меняет `status` на `available` и сбрасывает `checked_out_by`.

---

## TODO / идеи

- Авторизация: передавать ID сотрудника через заголовок (`X-Employee-ID`) вместо JSON в body.
- История движений контейнеров (лог checkout/return).
- Больше фильтров и сортировки (по веществу, по количеству, по сроку годности и т.п.).
- Простая мобильная/веб-клиентская морда с QR-сканером.