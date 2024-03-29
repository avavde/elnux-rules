package core

import (
	"fmt"
)

// Rule представляет правило.
type Rule struct {
	ID              string // Уникальный идентификатор правила
	EventType       string // Тип события, на который реагирует правило
	ConditionScript string // JavaScript код условия, как строка
	ActionScript    string // JavaScript код действия, как строка
}

// RuleEngine управляет правилами и их выполнением.
type RuleEngine struct {
	Rules    map[string]Rule     // Набор правил, индексированных по ID
	EventBus *EventBus           // Ссылка на систему событий
	Executor *JavaScriptExecutor // Экземпляр исполнителя JavaScript
}

// NewRuleEngine создает новый экземпляр RuleEngine.
func NewRuleEngine(eventBus *EventBus, executor *JavaScriptExecutor) *RuleEngine {
	return &RuleEngine{
		Rules:    make(map[string]Rule),
		EventBus: eventBus,
		Executor: executor,
	}
}

// RegisterRule регистрирует новое правило в движке.
func (engine *RuleEngine) RegisterRule(rule Rule) {
	engine.EventBus.Subscribe(rule.EventType, func(event Event) {
		if engine.Executor.EvaluateCondition(rule.ConditionScript, event.Payload) {
			fmt.Printf("Выполняется правило %s для события типа %s\n", rule.ID, rule.EventType)
			err := engine.Executor.Execute(rule.ActionScript)
			if err != nil {
				fmt.Printf("Ошибка при выполнении действия для правила %s: %v\n", rule.ID, err)
			} else {
				fmt.Printf("Действие для правила %s выполнено успешно\n", rule.ID)
			}
		} else {
			fmt.Printf("Условие для правила %s не выполнено для события типа %s\n", rule.ID, rule.EventType)
		}
	})
}

// AddRuleScript добавляет новое правило, выполняющее JavaScript.
func (engine *RuleEngine) AddRuleScript(id string, eventType string, conditionScript string, actionScript string) {
	rule := Rule{
		ID:              id,
		EventType:       eventType,
		ConditionScript: conditionScript,
		ActionScript:    actionScript,
	}
	engine.RegisterRule(rule)
	fmt.Printf("Добавлено новое правило %s для события типа %s\n", id, eventType)
}

func (engine *RuleEngine) LoadAndApplyRules(directory string) error {
	fmt.Println("Загрузка и применение правил из директории:", directory)
	rules, err := LoadRulesFromDirectory(directory, engine.Executor)
	if err != nil {
		fmt.Println("Ошибка при загрузке правил:", err)
		return err
	}

	for _, rule := range rules {
		engine.RegisterRule(rule)
	}
	fmt.Println("Правила успешно загружены и применены")
	return nil
}
