package services

import (
	"fmt"
	"math/rand"
	"temperature/models"
	"time"
)

// TemperatureService handles fetching temperature data from external API
type TemperatureService struct{}

// NewTemperatureService creates a new temperature service
func NewTemperatureService() *TemperatureService {
	return &TemperatureService{}
}

// GetTemperatureByLocation temperature data for a specific location
func (s *TemperatureService) GetTemperatureByLocation(location string) (*models.Temperature, error) {

	var temperatureResp = generateRandomTemperature(location)
	return &temperatureResp, nil
}

// GetTemperatureByID temperature data for a specific sensor ID
func (s *TemperatureService) GetTemperatureByID(sensorID string) (*models.Temperature, error) {

	var temperatureResp = generateRandomTemperature("")
	return &temperatureResp, nil
}

func generateRandomTemperature(location string) models.Temperature {
	rand.Seed(time.Now().UnixNano())

	// Список возможных значений для строковых полей
	sensorTypes := []string{"Digital", "Analog", "Infrared", "Thermocouple"}
	descriptions := []string{"Main sensor", "Backup sensor", "External sensor", "Internal sensor"}

	return models.Temperature{
		Value:       rand.Float64()*50 - 10, // Диапазон от -10 до +40
		Unit:        "°C",
		Timestamp:   time.Now(),
		Location:    location,
		Status:      "OK",
		SensorID:    fmt.Sprintf("SENS-%04d", rand.Intn(10000)),
		SensorType:  sensorTypes[rand.Intn(len(sensorTypes))],
		Description: descriptions[rand.Intn(len(descriptions))],
	}
}
