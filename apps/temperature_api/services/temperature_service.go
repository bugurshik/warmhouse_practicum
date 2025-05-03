package services

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"temperature/models"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// TemperatureService handles fetching temperature data from external API
type TemperatureService struct{ Channel *amqp.Channel }

// NewTemperatureService creates a new temperature service
func NewTemperatureService(ch *amqp.Channel) *TemperatureService {
	return &TemperatureService{Channel: ch}
}

// GetTemperatureByLocation temperature data for a specific location
func (s *TemperatureService) GetTemperatureByLocation(location string) (*models.Temperature, error) {
	var temperatureResp = generateRandomTemperature(location)
	err := s.PublishTemperature(&temperatureResp)
	if err != nil {
		log.Fatalf("Not sended on message brocker: %v\n", err)
	}
	return &temperatureResp, nil
}

// GetTemperatureByID temperature data for a specific sensor ID
func (s *TemperatureService) GetTemperatureByID(sensorID string) (*models.Temperature, error) {
	var temperatureResp = generateRandomTemperature("")
	err := s.PublishTemperature(&temperatureResp)
	if err != nil {
		log.Fatalf("Not sended on message brocker: %v\n", err)
	}
	return &temperatureResp, nil
}

func generateRandomTemperature(location string) models.Temperature {
	rand.Seed(time.Now().UnixNano())

	// Список возможных значений для строковых полей
	sensorTypes := []string{"Digital", "Analog", "Infrared", "Thermocouple"}
	descriptions := []string{"Main sensor", "Backup sensor", "External sensor", "Internal sensor"}

	sensorID := ""
	switch location {
	case "Living Room":
		sensorID = "1"
	case "Bedroom":
		sensorID = "2"
	case "Kitchen":
		sensorID = "3"
	default:
		sensorID = "0"
	}

	return models.Temperature{
		Value:       rand.Float64()*50 - 10, // Диапазон от -10 до +40
		Unit:        "°C",
		Timestamp:   time.Now(),
		Location:    location,
		Status:      "OK",
		SensorID:    sensorID,
		SensorType:  sensorTypes[rand.Intn(len(sensorTypes))],
		Description: descriptions[rand.Intn(len(descriptions))],
	}
}

func (s *TemperatureService) PublishTemperature(temp *models.Temperature) error {
	fmt.Println("send temperature!")

	q, err := s.Channel.QueueDeclare(
		"temperature-t",
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Not declare: %v\n", err)
		return err
	}

	// Сериализация в JSON
	body, err := json.Marshal(temp)
	if err != nil {
		log.Fatalf("Not serialized: %v\n", err)
		return err
	}

	err = s.Channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	return err
}
