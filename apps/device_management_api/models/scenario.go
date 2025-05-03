package models

type Condition struct {
	DeviceID string      `json:"deviceId"` // ID устройства-источника
	Property string      `json:"property"` // Давление, температура и т.д.
	Operator string      `json:"operator"` // ">", "<", "=="
	Value    interface{} `json:"value"`    // Пороговое значение
}

type Action struct {
	DeviceID string      `json:"deviceId"` // Целевое устройство
	Command  string      `json:"command"`  // "turnOn", "setTemperature"
	Payload  interface{} `json:"payload"`  // { "temp": 25 }
}

type Rule struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Condition Condition `json:"condition"`
	Action    Action    `json:"action"`
	Enabled   bool      `json:"enabled"`
}
