package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"smarthome/models"
)

// SensorManagementClient is an HTTP client for the sensor microservice
type SensorManagementClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewSensorManagementClient creates a new SensorManagementClient
func NewSensorManagementClient(baseURL string) *SensorManagementClient {
	return &SensorManagementClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateSensor creates a new sensor
func (c *SensorManagementClient) CreateSensor(ctx context.Context, sensorCreate models.SensorCreate) (*models.Sensor, error) {
	body, err := json.Marshal(sensorCreate)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/devices", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var sensor models.Sensor
	if err := json.NewDecoder(resp.Body).Decode(&sensor); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &sensor, nil
}

// UpdateSensor updates an existing sensor
func (c *SensorManagementClient) UpdateSensor(ctx context.Context, id int, sensorUpdate models.SensorUpdate) (*models.Sensor, error) {
	body, err := json.Marshal(sensorUpdate)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/devices/%d", c.baseURL, id), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var sensor models.Sensor
	if err := json.NewDecoder(resp.Body).Decode(&sensor); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &sensor, nil
}

// DeleteSensor deletes a sensor
func (c *SensorManagementClient) DeleteSensor(ctx context.Context, id int) error {
	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/devices/%d", c.baseURL, id), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
