package core

import (
	"fmt"
	"strings"
)

// RegisterTestRules подписывается на события с типом, начинающимся на "test", и обрабатывает их.
func RegisterTestRules(eventBus *EventBus, ruleEngine *RuleEngine) {
	for _, rule := range ruleEngine.Rules {
		if strings.HasPrefix(rule.EventType, "test") {
			localRule := rule // Создаем локальную копию для использования в замыкании
			eventBus.Subscribe(localRule.EventType, func(event Event) {
				if ruleEngine.Executor.EvaluateCondition(localRule.ConditionScript, event.Payload) {
					fmt.Printf("Выполняется скрипт для правила %s\n", localRule.ID)
					err := ruleEngine.Executor.Execute(localRule.ActionScript)
					if err != nil {
						fmt.Printf("Ошибка при выполнении скрипта правила %s: %v\n", localRule.ID, err)
					} else {
						fmt.Printf("Скрипт для правила %s выполнен успешно\n", localRule.ID)
					}
				} else {
					fmt.Printf("Условие правила %s не выполнено\n", localRule.ID)
				}
			})
		}
	}
}
