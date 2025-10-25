package main

import (
	"database/sql"
	"fmt"
	"log"
)

// Подключение в БД
func connectDB(pgHost string) *sql.DB {
	db, err := sql.Open("pgx", pgHost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Подключение к БД успешно")
	return db
}

// Запись в БД
func insertDB(db *sql.DB, order Order) {
	query := `
            INSERT INTO orders (order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
            VALUES ($1, $2, $3, $4::jsonb, $5::jsonb, $6::jsonb, $7, $8, $9, $10, $11, $12, $13, $14)
            ON CONFLICT (order_uid) DO NOTHING`
	result, err := db.Exec(query,
		order.OrderUID,
		order.TrackNumber,
		order.Entry,
		order.Delivery,
		order.Payment,
		order.Items,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.Shardkey,
		order.SmID,
		order.DateCreated,
		order.OofShard)

	if err != nil {
		log.Printf("Ошибка записи в БД: %v", err)
		return
	}

	rowAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Ошибка получения количества строк: %v", err)
		return
	}

	// Проверка на дубликат заказа
	if rowAffected == 0 {
		log.Printf("Заказ %v уже существует, дубликат не был добавлен в БД", order.OrderUID)
	} else {
		log.Printf("Заказ %v успешно записан в БД", order.OrderUID)
	}
}
