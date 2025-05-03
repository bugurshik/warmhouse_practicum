package handlers

import (
	"fmt"
	"net/http"
	"temperature/services"

	"github.com/gin-gonic/gin"
)

// TemperatureHandler handles temperature-related requests
type TemperatureHandler struct {
	TemperatureService *services.TemperatureService
}

// NewSensorHandler creates a new SensorHandler
func NewTemperatureHandler(temperatureService *services.TemperatureService) *TemperatureHandler {
	return &TemperatureHandler{
		TemperatureService: temperatureService,
	}
}

// RegisterRoutes registers the temperature routes
func (h *TemperatureHandler) RegisterRoutes(router *gin.RouterGroup) {
	sensors := router.Group("temperature")
	{
		sensors.GET("", h.GetTemperatureByLocation)
		sensors.GET("/:sensorID", h.GetTemperatureById)
	}
}

// @Summary Получить данные температуры по локации датчика
// @Description Возвращает температурные данные для указанного датчика, включая местоположение, значение, единицы измерения и временную метку
// @Tags Temperature
// @Accept json
// @Produce json
// @Param location query string true "Локация датчика"
// @Success 200 {object} map[string]interface{} "Успешный ответ" Example({"location": "Комната 5", "value": 23.5, "unit": "C", "status": "normal", "timestamp": "2023-05-15T14:30:00Z", "description": "Текущая температура в помещении"})
// @Failure 400 {object} map[string]string "Неверный запрос" Example({"error": "Location is required"})
// @Failure 404 {object} map[string]string "Датчик не найден"
// @Failure 500 {object} map[string]string "Ошибка сервера" Example({"error": "Failed to fetch temperature data: connection timeout"})
// @Router /temperature [get]
// GetTemperatureByLocation handles GET /temperature
func (h *TemperatureHandler) GetTemperatureByLocation(c *gin.Context) {
	location := c.Query("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location is required"})
		return
	}

	// Fetch temperature data from the external API
	tempData, err := h.TemperatureService.GetTemperatureByLocation(location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch temperature data: %v", err),
		})
		return
	}

	// Return the temperature data
	c.JSON(http.StatusOK, gin.H{
		"location":    tempData.Location,
		"value":       tempData.Value,
		"unit":        tempData.Unit,
		"status":      tempData.Status,
		"timestamp":   tempData.Timestamp,
		"description": tempData.Description,
	})
}

// GetTemperatureById godoc
// @Summary Получить данные температуры по ID датчика
// @Description Возвращает температурные данные для указанного датчика, включая местоположение, значение, единицы измерения и временную метку
// @Tags Temperature
// @Accept json
// @Produce json
// @Param sensorId path string true "Уникальный идентификатор датчика температуры"
// @Success 200 {object} map[string]interface{} "Успешный ответ" Example({"location": "Комната 5", "value": 23.5, "unit": "C", "status": "normal", "timestamp": "2023-05-15T14:30:00Z", "description": "Текущая температура в помещении"})
// @Failure 400 {object} map[string]string "Неверный запрос" Example({"error": "Location is required"})
// @Failure 404 {object} map[string]string "Датчик не найден"
// @Failure 500 {object} map[string]string "Ошибка сервера" Example({"error": "Failed to fetch temperature data: connection timeout"})
// @Router /temperature/{sensorId} [get]
// GetTemperatureByLocation handles GET /temperature/:sensorId
func (h *TemperatureHandler) GetTemperatureById(c *gin.Context) {
	sensorID := c.Param("sensorID")
	if sensorID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location is required"})
		return
	}

	// Fetch temperature data from the external API
	tempData, err := h.TemperatureService.GetTemperatureByID(sensorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to fetch temperature data: %v", err),
		})
		return
	}

	// Return the temperature data
	c.JSON(http.StatusOK, gin.H{
		"location":    tempData.Location,
		"value":       tempData.Value,
		"unit":        tempData.Unit,
		"status":      tempData.Status,
		"timestamp":   tempData.Timestamp,
		"description": tempData.Description,
	})
}
