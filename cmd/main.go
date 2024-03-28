package main

import (
	"elnux-rules/internal/core"
	"fmt"
	"log"
	"time"
)

func main() {
	eventBus := core.NewEventBus()
	javaScriptExecutor := core.NewJavaScriptExecutor()
	ruleEngine := core.NewRuleEngine(eventBus, javaScriptExecutor)

	// Загрузка и применение всех правил из директории .rules

	if err := ruleEngine.LoadAndApplyRules("./rules"); err != nil {
		log.Fatalf("Не удалось загрузить и применить правила: %v", err)
	}
	log.Println("Правила успешно загружены и применены")

	// Пример добавления пользовательского скрипта через RuleEngine
	script := `console.log("Обработано событие exampleEvent");`
	ruleEngine.AddRuleScript("exampleRule", "exampleEvent", "true", script)

	// Имитация события для активации правила
	event := core.Event{
		Type:    "exampleEvent",
		Payload: map[string]interface{}{"message": "test"},
	}
	ruleEngine.EventBus.Publish(event) // Публикация события
	fmt.Printf("Событие %s обработано подписчиком\n", event.Type)

	eventBus.Publish(core.Event{
		Type: "temperatureChange",
		Payload: map[string]interface{}{
			"temperature": 30, // Температура изменилась
		},
	})

	// Дать время для асинхронного выполнения скрипта
	time.Sleep(1 * time.Second)
}
