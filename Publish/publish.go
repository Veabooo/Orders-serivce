package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nats-io/stan.go"
)

// Публикация сообщений в канал "orders"
func main() {
	file := "model.json"

	// Читаем JSON из файла
	data, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("Не могу прочитать '%v': %v", file, err)
	}

	// Подключаемся к NATS Streaming
	sc, err := stan.Connect("cluster", "publisher", stan.NatsURL("nats://nats:4222"))
	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
	defer sc.Close()

	// Отправка в канал
	if err := sc.Publish("orders", data); err != nil {
		log.Fatal("Ошибка отправки:", err)
	}

	fmt.Println("Сообщение успешно отправлено в канал 'orders'!")
}
