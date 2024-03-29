package core

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateAndPublishEvents генерирует и публикует события на основе заданных правил.
func GenerateAndPublishEvents(eventBus *EventBus, ruleEngine *RuleEngine) {
	fmt.Println("Начало генерации и публикации событий")
	// Генерация тестовых событий на основе типов событий из правил
	for _, rule := range ruleEngine.Rules {
		// Проверяем, необходимо ли генерировать событие для данного правила
		if shouldGenerateEventForRule(rule) {
			// Генерация данных события на основе логики, связанной с типом события
			payload := generatePayloadForRule(rule)
			event := Event{
				Type:    rule.EventType,
				Payload: payload,
			}
			fmt.Printf("Сгенерировано событие: %v\n", event)
			// Публикация события
			eventBus.Publish(event)
		}
	}
	fmt.Println("Генерация и публикация событий завершены")
}

// shouldGenerateEventForRule определяет, следует ли генерировать событие для данного правила.
func shouldGenerateEventForRule(rule Rule) bool {

	fmt.Printf("Проверка генерации события для правила: %s\n", rule.ID)
	return true
}

// generatePayloadForRule генерирует пейлоад для события на основе правила.
func generatePayloadForRule(rule Rule) map[string]interface{} {

	fmt.Printf("Генерация пейлоада для события типа: %s\n", rule.EventType)
	switch rule.EventType {
	case "testmotionDetected":
		return map[string]interface{}{
			"motion": true,
			"time":   time.Now(),
		}
	case "testtemperatureChange":
		return map[string]interface{}{
			"temperature": rand.Intn(10) + 30, // Случайное значение температуры от 20 до 30
		}

	default:
		return map[string]interface{}{}
	}
}
