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
	sensors := router.Group("/temperature")
	{
		sensors.GET("", h.GetTemperatureByLocation)
		//sensors.GET("/:sensorId", h.GetTemperatureByLocation)
	}
}

// @Summary GetTemperatureByLocation
// @Description Получить температуру по локации.
// @Tags Temperature
// @Accept  json
// @Produce  json
// @Param   location      query    string  true  "Локация"
// @Success 200 {string} string "Возвращает текущие показания сенсора"
// @Router /example [get]
// GetTemperatureByLocation handles GET /api/v1/sensors/temperature/:location
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
