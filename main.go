package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	pgHost := "postgres://user1:user1password@localhost:5432/order_db"
	clusterID := "cluster"
	clientID := "client"
	natsHost := "nats://localhost:4222"
	chanNATS := "orders"

	// Подключение к БД
	db := connectDB(pgHost)
	defer db.Close()

	// Восстановление кэша из БД
	if err := restoreCacheFromDB(db); err != nil {
		log.Printf("Ошибка восстановления кэша: %v", err)
	}

	// Подключение к NATSStreaming
	sc, err := connectNATSStreaming(clusterID, clientID, natsHost)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	// Подписка на канал
	if err := subscribeNATSStreaming(sc, chanNATS, db); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/order", orderHandler)

	go func() {
		fmt.Println("HTTP-сервер запущен на http://localhost:8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	fmt.Printf("Сервис запущен, слушает канал '%v'...\n", chanNATS)
	select {} // Бесконечный цикл
}
