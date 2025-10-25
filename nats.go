package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

// Подключение к NATS
// Подключение к NATS Streaming
func connectNATSStreaming(clusterID, clientID, natsURL string) (stan.Conn, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Соединение с NATS Streaming потеряно: %v", reason)
		}))
	if err != nil {
		return nil, err
	}
	fmt.Println("Подключение к NATS Streaming успешно")
	return sc, nil
}

// Подписка на канал
func subscribeNATSStreaming(sc stan.Conn, chanNATS string, db *sql.DB) error {
	var err error
	_, err = sc.Subscribe(chanNATS, func(msg *stan.Msg) {
		var order Order

		// Парсинг JSON
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("Ошибка при парсинге JSON: %v", err)
			return
		}

		// Валидация заказа
		if err := validateOrder(order); err != nil {
			log.Printf("Заказ не прошёл валидацию: %v", err)
			return
		}

		// Проверка дубликата в кэше
		if _, exists := cache.Load(order.OrderUID); exists {
			log.Printf("Заказ %v уже в кэше, дубликат проигнорирован", order.OrderUID)
			if err := msg.Ack(); err != nil {
				log.Printf("Ошибка подтверждения дубликата: %v", err)
			}
			return
		}

		// Запись в БД
		insertDB(db, order)

		// Сохранение полученных данных в кэш
		cacheStore(order)

		// Подтвреждение получения и обработки сообщения
		if err := msg.Ack(); err != nil {
			log.Printf("Ошибка подтверждения сообщения: %v", err)
		}
	},
		stan.SetManualAckMode(),
		stan.DurableName("order-durable"),
		stan.AckWait(10*time.Second),
		stan.StartWithLastReceived(),
	)
	return err
}

// Валидация заказа
func validateOrder(order Order) error {
	if order.OrderUID == "" {
		return fmt.Errorf("order_uid пустой")
	}
	if order.TrackNumber == "" {
		return fmt.Errorf("track_number пустой")
	}
	return nil
}
