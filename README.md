# Order Service (NATS Streaming + PostgreSQL + Web UI)

Микросервис для обработки заказов через **NATS Streaming**, с хранением в **PostgreSQL**, кэшированием и веб-интерфейсом.

---

## Функции

- Подписка на канал `orders` (NATS Streaming)
- Валидация, запись в БД, кэширование
- Восстановление кэша при старте
- HTTP: `/:8080` — поиск, `/order?id=...` — просмотр
- Паблишер: `publish/publish.go`

---

## Запуск (Docker)

```bash
docker-compose up --build