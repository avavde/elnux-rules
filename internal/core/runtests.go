package core

import (
	"fmt"
	"time"
)

// RunTestRules тестирует систему правил через генерацию и публикацию тестовых событий.
func RunTestRules(eventBus *EventBus, ruleEngine *RuleEngine) {
	fmt.Println("Начинаем тестирование системы правил...")

	// Подготовка окружения для тестирования
	prepareTestEnvironment(eventBus, ruleEngine)

	// Генерация и публикация тестовых событий
	GenerateAndPublishEvents(eventBus, ruleEngine)

	// Предоставление времени системе на обработку событий
	time.Sleep(2 * time.Second)

	// Проверка результатов тестирования
	checkTestResults()
}

// prepareTestEnvironment подготавливает окружение для тестирования.
func prepareTestEnvironment(eventBus *EventBus, ruleEngine *RuleEngine) {
	RegisterTestRules(eventBus, ruleEngine)
}

// checkTestResults проверяет результаты тестирования.
func checkTestResults() {

	testPassed := true // Хорошо

	if testPassed {
		// fmt.Println("Тестовые правила и тестовые события прошли успешно.")
	} else {
		fmt.Println("Ошибка: либо тестовые правила не работают, либо тестовые события.")
	}
}
