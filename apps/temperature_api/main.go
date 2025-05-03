package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"temperature/handlers"
	"temperature/services"

	_ "temperature/docs"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Print("temperature run")

	// Initialize router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API routes
	apiRoutes := router.Group("")

	ch, err := channelRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}

	temperatureService := services.NewTemperatureService(ch)

	// Register sensor routes
	temperatureHandler := handlers.NewTemperatureHandler(temperatureService)
	temperatureHandler.RegisterRoutes(apiRoutes)

	// Start server
	srv := &http.Server{
		Addr:    getEnv("PORT", ":8081"),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited properly")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func TestTelemetry() {
	ch, err := channelRabbitMQ()
	queueName := "temperature-telemetry"

	// Запуск Consumer
	err = consumeMessages(ch, queueName)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	// Отправка тестовых сообщений (Producer)
	for i := 0; i < 5; i++ {
		msg := fmt.Sprintf("Sensor data %d", i)
		err = publishMessage(ch, queueName, msg)
		if err != nil {
			log.Printf("Failed to publish: %v", err)
		} else {
			log.Printf("Sent: %s", msg)
		}
		time.Sleep(2 * time.Second)
	}

	// Чтобы consumer не завершился сразу
	time.Sleep(10 * time.Second)
}

func connectRabbitMQ() (*amqp.Connection, error) {
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func channelRabbitMQ() (*amqp.Channel, error) {
	conn, err := connectRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}

	return ch, err
}

func publishMessage(ch *amqp.Channel, queueName string, message string) error {
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable (сохранять после перезагрузки)
		false,     // autoDelete (удалять, когда нет подписчиков)
		false,     // exclusive (только для текущего соединения)
		false,     // noWait (не ждать подтверждения)
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange (по умолчанию)
		q.Name, // routing key (имя очереди)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent, // сохранять на диск
		},
	)
	return err
}

func consumeMessages(ch *amqp.Channel, queueName string) error {
	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer (идентификатор)
		false,  // auto-ack (отключаем для ручного подтверждения)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return err
	}

	// Обработка сообщений в горутине
	go func() {
		for msg := range msgs {
			log.Printf("Received: %s", msg.Body)
			// Обработка сообщения...
			time.Sleep(1 * time.Second) // Имитация работы
			msg.Ack(false)              // Подтверждаем обработку
		}
	}()

	log.Println("Consumer started. Waiting for messages...")
	return nil
}
