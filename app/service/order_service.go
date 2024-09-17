package service

import (
	"awesomeProject/app/domain/dao"
	"awesomeProject/app/pkg"
	"awesomeProject/app/repository"
	"awesomeProject/app/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type OrderService interface {
	CreateOrder(c *gin.Context)
}

type OrderServiceImpl struct {
	orderRepository repository.OrderRepository
}

func (u OrderServiceImpl) CreateOrder(c *gin.Context) {
	defer pkg.PanicHandler(c)

	var order dao.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := u.orderRepository.CreateOrder(&order)
	if err != nil {
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	orderJSON, err := json.Marshal(&order)
	if err != nil {
		c.JSON(500, gin.H{"error": "JSON marshal error"})
		return
	}

	// Сохраняем JSON-строку в Redis
	ctx, err := utils.AddToRedis(order.ID, string(orderJSON))
	if err != nil {
		c.JSON(500, gin.H{"error": "Cache error"})
		return
	}

	kafkaWriter := kafka.Writer{
		Addr:  kafka.TCP("kafka:9092", "kafka:9093", "kafka:9094"),
		Topic: "orders",
	}
	// Публикуем событие о создании заказа в Kafka
	message := kafka.Message{
		Key:   []byte(order.ID),
		Value: []byte("Order Created"),
	}
	if err := kafkaWriter.WriteMessages(ctx, message); err != nil {
		log.Println("Failed to send message to Kafka:", err)
	}

	c.JSON(200, gin.H{"status": "Order created"})
}

func OrderServiceInit(orderRepository repository.OrderRepository) *OrderServiceImpl {
	return &OrderServiceImpl{
		orderRepository: orderRepository,
	}
}
