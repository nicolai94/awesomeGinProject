package handlers

import (
	"awesomeProject/app/domain/dao"
	"awesomeProject/app/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func ProcessOrders(db *gorm.DB) {
	// Инициализируем Kafka Writer внутри функции
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-signals
		fmt.Println("Got signal: ", sig)
		cancel()
	}()
	kafkaWriter := kafka.Writer{
		Addr:  kafka.TCP("kafka:9092", "kafka:9093", "kafka:9094"),
		Topic: "orders",
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "orders",
		GroupID: "order-processors",
	})
	defer func() {
		err := r.Close()
		if err != nil {
			fmt.Println("Error closing consumer: ", err)
			return
		}
		fmt.Println("Consumer closed")
	}()
	for {
		// Читаем сообщение из Kafka
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Println("Error reading message:", err)
			fmt.Printf("message at offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))
		}

		log.Printf("Received message: %s\n", string(msg.Value))
		// Обновляем статус заказа в базе данных
		var order dao.Order
		if err := db.First(&order, "id = ?", string(msg.Key)).Error; err != nil {
			log.Printf("Order not found in database for ID: %s\n", msg.Key)
			// Не коммитим смещение, если не удалось найти заказ
			continue
		}

		order.Status = "Processed"
		if err := db.Save(&order).Error; err != nil {
			log.Printf("Error updating order in database for ID: %s\n", msg.Key)
			// Не коммитим смещение, если возникла ошибка при сохранении
			continue
		}

		orderJSON, err := json.Marshal(&order)
		if err != nil {
			log.Printf("Error marshalling order to JSON for ID: %s\n", msg.Key)
			continue
		}

		// Обновляем данные в Redis
		_, err = utils.AddToRedis(order.ID, orderJSON)
		if err != nil {
			log.Printf("Error setting order in Redis for ID: %s\n", msg.Key)
			// Не коммитим смещение, если возникла ошибка при сохранении в Redis
			continue
		}

		// Публикуем событие об обновлении заказа обратно в Kafka
		if err := kafkaWriter.WriteMessages(ctx, kafka.Message{
			Key:   []byte(order.ID),
			Value: []byte("Order Processed"),
		}); err != nil {
			log.Printf("Error publishing message to Kafka for ID: %s\n", order.ID)
			// Не коммитим смещение, если возникла ошибка при публикации
			continue
		}

		// Коммитим смещение после успешной обработки
		if err := r.CommitMessages(ctx, msg); err != nil {
			log.Printf("Error committing offset for ID: %s\n", order.ID)
		} else {
			log.Printf("Successfully committed offset: %d for ID: %s\n", msg.Offset, order.ID)
		}
	}
}
