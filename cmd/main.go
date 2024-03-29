package main

import (
	"elnux-rules/internal/core"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка при загрузке файла .env: %v", err)
	}

	// Пример чтения переменной окружения из .env
	dbHost := os.Getenv("DB_HOST")
	fmt.Println("DB_HOST:", dbHost)

	eventBus := core.NewEventBus() // шина событий
	time.Sleep(1 * time.Second)
	javaScriptExecutor := core.NewJavaScriptExecutor() // executor.og
	time.Sleep(1 * time.Second)
	ruleEngine := core.NewRuleEngine(eventBus, javaScriptExecutor) // движок правил

	// Загрузка и применение всех правил из директории .rules

	if err := ruleEngine.LoadAndApplyRules("./rules"); err != nil {
		log.Fatalf("Не удалось загрузить и применить правила: %v", err)
	}
	log.Println("Правила успешно загружены и применены")

	// Регистрируем тестовые правила
	core.RegisterTestRules(eventBus, ruleEngine)

	testEvents := []core.Event{
		{
			Type: "testhumidityChange",
			Payload: map[string]interface{}{
				"humidity": 20, // должно запустить "Внимание: влажность ниже 30%!"
			},
		},
		{
			Type: "testmotionDetected",
			Payload: map[string]interface{}{
				"motion": true,
				"time":   "2024-03-27T12:00:00Z", // должно запустить "Движение обнаружено. Включение освещения."
			},
		},
		{
			Type: "testtemperatureChange",
			Payload: map[string]interface{}{
				"temperature": 30, // должно запустить "Внимание: температура превысила 25 градусов!"
			},
		},
		{
			Type: "test-ExampleEvent",
			Payload: map[string]interface{}{
				"message": "test", // должно запустить "Температура превысила 30 градусов!"
			},
		},
	}

	// опубликовать все тестовые события
	for _, event := range testEvents {
		eventBus.Publish(event)
	}

	// Дать время для асинхронного выполнения скрипта

	time.Sleep(3 * time.Second)
}
