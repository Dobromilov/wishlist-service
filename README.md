# Wishlist API

REST API сервис для создания вишлистов (списков желаний) к праздникам и событиям.

## Возможности

- Регистрация и авторизация по email + пароль (JWT)
- CRUD вишлистов (название, описание, дата события)
- CRUD подарков внутри вишлиста (название, ссылка, приоритет)
- Публичный доступ к вишлисту по уникальной ссылке
- Бронирование подарков без авторизации

## Запуск

```bash
cp .env.example .env
docker-compose up --build
```

Сервер запустится на `http://localhost:8080`.

## API

### Авторизация

**POST /auth/register**
```json
{"email": "user@example.com", "password": "secret123"}
```

**POST /auth/login**
```json
{"email": "user@example.com", "password": "secret123"}
```

Ответ: `{"token": "jwt-token-here"}`

### Вишлисты (требует Authorization: Bearer)

**POST /api/wishlists**
```json
{"title": "День рождения", "description": "Мои подарки", "event_date": "2025-12-25"}
```

**GET /api/wishlists** — список моих вишлистов

**PUT /api/wishlists/{id}** — обновление

**DELETE /api/wishlists/{id}** — удаление

### Подарки (требует Authorization: Bearer)

**POST /api/wishlists/{id}/items**
```json
{"name": "Наушники", "url": "https://...", "priority": 5}
```

**GET /api/wishlists/{id}/items** — список подарков

**DELETE /api/wishlists/{id}/items/{itemId}** — удаление

### Публичный доступ

**GET /public/{token}** — просмотр вишлиста со всеми подарками

**POST /public/{token}/reserve**
```json
{"item_id": 1}
```

## Структура проекта

```
├── cmd/server/          — точка входа
├── internal/
│   ├── config/          — конфигурация
│   ├── domain/          — модели
│   ├── handler/         — HTTP handlers
│   ├── service/         — бизнес-логика
│   ├── storage/         — работа с БД
│   └── repository/      — интерфейсы репозиториев
├── migrations/          — SQL миграции
├── Dockerfile
├── docker-compose.yml
└── .env.example
```
