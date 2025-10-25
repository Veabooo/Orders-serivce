# Orders service (NATS Streaming + PostgreSQL + Web UI)

Микросервис для обработки заказов через **NATS Streaming**, с хранением в **PostgreSQL**, кэшированием и веб-интерфейсом.

---

## Функции

- Подписка на канал `orders` (NATS Streaming)
- Валидация, запись в БД, кэширование
- Восстановление кэша при старте
- HTTP: `/:8080` — поиск заказов
- Паблишер: `publish/publish.go`

---

### Видео-демонстрация работы сервиса и интерфейса

- https://disk.yandex.ru/i/mGPLwmaF-TxQ3A

---

## Запуск (Docker)

```bash
docker-compose build
docker-compose up -d

#Выполнение скрипта для публикации сообщения в канал
docker-compose run --rm order-service ./publisher
