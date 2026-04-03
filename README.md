# lab-inventory-service

Простой REST API для учёта лабораторного инвентаря: сотрудников, помещений, химических веществ и конкретных ёмкостей с реагентами. Сейчас сервис уже поддерживает создание и просмотр основных сущностей, поиск по химикатам, фильтрацию контейнеров, а также операции выдачи и возврата контейнера. [cite:12]

## Стек

- Go
- chi
- PostgreSQL
- pgx / pgxpool
- sqlc

## Что умеет сервис

Сервис работает с четырьмя основными сущностями:

- **employees** — сотрудники лаборатории. [cite:12]
- **rooms** — помещения / комнаты. [cite:12]
- **chemicals** — химические вещества с именем, CAS, формулой и SDS URL. [cite:12]
- **containers** — конкретные физические ёмкости с веществом, привязанные к комнате. [cite:12]

Дополнительно уже реализованы:

- поиск химикатов по имени, CAS или формуле через один query-параметр. [cite:12]
- фильтрация контейнеров по статусу и комнате. [cite:12]
- checkout / return контейнера. [cite:12]

## Структура данных

### employees

- `id`
- `full_name`
- `role`
- `created_at`

### rooms

- `id`
- `name`
- `description`
- `created_at`

### chemicals

- `id`
- `name`
- `cas_number`
- `formula`
- `sds_url`
- `created_at`

### containers

- `id`
- `chemical_id`
- `room_id`
- `label_code`
- `quantity`
- `unit`
- `status`
- `checked_out_by`
- `created_at`

`containers` — это не “вид вещества”, а конкретные физические бутылки / флаконы / канистры с реагентом. Одна запись в `chemicals` может иметь много `containers`. [cite:12]

## Запуск

### 1. Поднять PostgreSQL

Нужна база данных PostgreSQL с именем `lab_inventory`.

По умолчанию сервис подключается через DSN вида:

```text
postgres://postgres:postgres@localhost:5432/lab_inventory?sslmode=disable
```

Если нужно, можно задать свой DSN через переменную окружения:

```bash
LAB_INVENTORY_DB_DSN=postgres://USER:PASSWORD@HOST:PORT/lab_inventory?sslmode=disable
```

### 2. Применить схему

Создай таблицы из `sql/schema.sql`.

### 3. Сгенерировать sqlc-код

```bash
sqlc generate
```

### 4. Запустить сервис

```bash
go run ./cmd/lab-inventory-service
```

Сервис слушает:

```text
http://localhost:8080
```

## API

### Health

```http
GET /health
```

Ответ:

```text
ok
```

---

### Employees

#### Список сотрудников

```http
GET /employees
```

#### Создать сотрудника

```http
POST /employees
Content-Type: application/json
```

```json
{
  "full_name": "Ivan Petrov",
  "role": "chemist"
}
```

---

### Rooms

#### Список комнат

```http
GET /rooms
```

#### Создать комнату

```http
POST /rooms
Content-Type: application/json
```

```json
{
  "name": "Room 101",
  "description": "Organic synthesis lab"
}
```

---

### Chemicals

#### Список химикатов

```http
GET /chemicals
```

#### Создать химикат

```http
POST /chemicals
Content-Type: application/json
```

```json
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

Поиск работает по **одному** query-параметру `query`: сервис ищет это значение сразу в `name`, `cas_number` и `formula`. То есть можно передать либо имя вещества, либо CAS, либо формулу — не нужно заполнять всё сразу. [cite:12]

Примеры:

```http
GET /chemicals/search?query=acetone
GET /chemicals/search?query=67-64-1
GET /chemicals/search?query=C3H6O
```

---

### Containers

#### Список контейнеров

```http
GET /containers
```

#### Создать контейнер

```http
POST /containers
Content-Type: application/json
```

```json
{
  "chemical_id": 1,
  "room_id": 1,
  "label_code": "ACET-001",
  "quantity": 0.5,
  "unit": "L",
  "status": "available"
}
```

#### Фильтрация контейнеров

```http
GET /containers?status=available&room_id=1
```

Поддерживаются query-параметры:

- `status`
- `room_id`

Можно использовать как оба сразу, так и только один из них. [cite:12]

Примеры:

```http
GET /containers?status=available
GET /containers?room_id=1
GET /containers?status=checked_out&room_id=1
```

#### Взять контейнер

```http
POST /containers/{id}/checkout
Content-Type: application/json
```

```json
{
  "employee_id": 1
}
```

После этого контейнер переводится в статус `checked_out`, а поле `checked_out_by` заполняется id сотрудника. [cite:12]

#### Вернуть контейнер

```http
POST /containers/{id}/return
```

После этого контейнер переводится обратно в статус `available`, а `checked_out_by` сбрасывается в `null`. [cite:12]

## Примеры запросов

### PowerShell

```powershell
curl http://localhost:8080/health -UseBasicParsing

curl http://localhost:8080/employees `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"full_name":"Ivan Petrov","role":"chemist"}' `
  -UseBasicParsing

curl http://localhost:8080/rooms `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"name":"Room 101","description":"Organic synthesis lab"}' `
  -UseBasicParsing

curl http://localhost:8080/chemicals `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"name":"Acetone","cas_number":"67-64-1","formula":"C3H6O","sds_url":"https://example.com/acetone-sds"}' `
  -UseBasicParsing

curl "http://localhost:8080/chemicals/search?query=acetone" -UseBasicParsing

curl http://localhost:8080/containers `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"chemical_id":1,"room_id":1,"label_code":"ACET-001","quantity":0.5,"unit":"L","status":"available"}' `
  -UseBasicParsing

curl http://localhost:8080/containers/1/checkout `
  -Method Post `
  -ContentType "application/json" `
  -Body '{"employee_id":1}' `
  -UseBasicParsing

curl http://localhost:8080/containers/1/return `
  -Method Post `
  -UseBasicParsing

curl "http://localhost:8080/containers?status=available&room_id=1" -UseBasicParsing
```

## Текущее состояние

На текущем этапе это MVP лабораторного inventory-сервиса: есть CRUD-основа для основных сущностей, выдача/возврат контейнера, поиск химикатов и фильтрация контейнеров. [cite:12]

## Возможные следующие шаги

- история движений контейнеров;
- валидация бизнес-логики для checkout / return;
- простая авторизация;
- QR-коды для `label_code`;
- мобильный клиент или web UI. [cite:12]