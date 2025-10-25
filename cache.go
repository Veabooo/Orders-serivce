package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

var cache = sync.Map{}

// Сохранение полученных данных в кэш
func cacheStore(order Order) {
	cache.Store(order.OrderUID, order)
	log.Printf("В кэш добавлен заказ: %v", order.OrderUID)
}

// Восстановление кэша из БД
func restoreCacheFromDB(db *sql.DB) error {
	var dbCount int

	query := `
        SELECT 
            order_uid, track_number, entry, locale, internal_signature,
            customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard,
            delivery, payment, items
        FROM orders`

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var order Order
		var deliveryJSON, paymentJSON, itemsJSON []byte

		err := rows.Scan(
			&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated,
			&order.OofShard, &deliveryJSON, &paymentJSON, &itemsJSON)
		if err != nil {
			log.Printf("Ошибка сканирования строки: %v", err)
			continue
		}

		// Восстанавливаем вложенные структуры
		if err := json.Unmarshal(deliveryJSON, &order.Delivery); err != nil {
			log.Printf("Ошибка парсинга delivery: %v", err)
			continue
		}
		if err := json.Unmarshal(paymentJSON, &order.Payment); err != nil {
			log.Printf("Ошибка парсинга payment: %v", err)
			continue
		}
		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			log.Printf("Ошибка парсинга items: %v", err)
			continue
		}

		cache.Store(order.OrderUID, order)
		count++
	}

	if err := db.QueryRow("SELECT COUNT(*) FROM orders").Scan(&dbCount); err != nil {
		return fmt.Errorf("ошибка подсчёта заказов в БД: %w", err)
	}

	fmt.Printf("Кэш восстановлен: %v заказов загружено из БД \nВсего заказов в БД: %v\n", count, dbCount)
	return nil
}
