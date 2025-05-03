package handlers

import (
	"log"
	"device_management/models"
)

type RuleManager struct {
    rules map[string]models.Rule
}

func NewRuleManager() *RuleManager {
    return &RuleManager{
        rules: make(map[string]models.Rule),
    }
}

// Добавление правила
func (rm *RuleManager) AddRule(rule models.Rule) {
    rm.rules[rule.ID] = rule
    log.Printf("Правило добавлено: %s", rule.Name)
}

// Проверка условия и выполнение действия
func (rm *RuleManager) Evaluate(deviceID string, property string, value interface{}) {
    for _, rule := range rm.rules {
        if !rule.Enabled {
            continue
        }
        if rule.Condition.DeviceID == deviceID && rule.Condition.Property == property {
            if checkCondition(rule.Condition, value) {
                executeAction(rule.Action)
            }
        }
    }
}

// Проверка условия (">", "<", "==")
func checkCondition(cond models.Condition, value interface{}) bool {
    switch cond.Operator {
    case ">":
        return value.(float64) > cond.Value.(float64)
    case "<":
        return value.(float64) < cond.Value.(float64)
    case "==":
        return value == cond.Value
    }
    return false
}

// Отправка команды на устройство (через брокер/HTTP)
func executeAction(action models.Action) {
    log.Printf("Выполняется действие: %s -> %s", action.DeviceID, action.Command)
    // Реализация: NATS/Kafka/REST вызов
}