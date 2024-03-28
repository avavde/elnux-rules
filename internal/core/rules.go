package core

import (
	"fmt"

	"github.com/dop251/goja"
)

// ConditionFunc определяет тип функции условия.
type ConditionFunc func(event Event) bool

// ActionFunc определяет тип функции действия.
type ActionFunc func(event Event)

// Rule представляет правило.
type Rule struct {
	ID              string     // Уникальный идентификатор правила
	EventType       string     // Тип события, на который реагирует правило
	ConditionScript string     // JavaScript код условия, как строка
	Action          ActionFunc // Функция действия, выполняемая при активации правила
	Script          string     // JavaScript скрипт, который выполняется как действие
}

// RuleEngine управляет правилами и их выполнением.
type RuleEngine struct {
	Rules    map[string]Rule     // Набор правил, индексированных по ID
	EventBus *EventBus           // Ссылка на систему событий
	Executor *JavaScriptExecutor // Добавляем поле для хранения экземпляра JavaScriptExecutor
}

// NewRuleEngine создает новый экземпляр RuleEngine.
func NewRuleEngine(eventBus *EventBus, executor *JavaScriptExecutor) *RuleEngine { // Изменяем функцию NewRuleEngine для принятия экземпляра JavaScriptExecutor
	return &RuleEngine{
		Rules:    make(map[string]Rule),
		EventBus: eventBus,
		Executor: executor, // Присваиваем переданный экземпляр JavaScriptExecutor полю Executor
	}
}

// RegisterRule регистрирует новое правило в движке.
// Псевдокод
func (engine *RuleEngine) RegisterRule(rule Rule) {
	fmt.Printf("Регистрация правила: %s с скриптом: %s\n", rule.ID, rule.Script)
	engine.EventBus.Subscribe(rule.EventType, func(event Event) {
		if engine.EvaluateCondition(rule.ConditionScript, event) {
			fmt.Println("Выполнение скрипта правила:", rule.Script) // Для отладки
			engine.ExecuteScript(rule.Script)
		}
	})
}

// EvaluateCondition выполняет JavaScript условие и возвращает bool результат.
// EvaluateCondition выполняет JavaScript условие и возвращает bool результат.
func (engine *RuleEngine) EvaluateCondition(conditionScript string, event Event) bool {
	runtime := goja.New()
	runtime.Set("event", event.Payload) // Передача Payload для доступа к свойствам события
	wrappedConditionScript := fmt.Sprintf("(function() { return %s; })()", conditionScript)
	value, err := runtime.RunString(wrappedConditionScript)
	if err != nil {
		fmt.Println("Ошибка при оценке условия:", err)
		return false
	}
	result := value.ToBoolean() // Правильное использование результатов вызова
	fmt.Printf("Результат оценки условия для скрипта %s: %t\n", conditionScript, result)
	return result
}

// ExecuteRule выполняет указанное правило немедленно.
func (engine *RuleEngine) ExecuteRule(ruleID string, event Event) {
	if rule, exists := engine.Rules[ruleID]; exists {
		rule.Action(event)
	}
}

// ExecuteScript выполняет JavaScript код с помощью Goja.
// В предположении, что у вас есть метод Execute в JavaScriptExecutor:
func (engine *RuleEngine) ExecuteScript(script string) {
	fmt.Printf("Выполнение скрипта: %s\n", script)
	err := engine.Executor.Execute(script)
	if err != nil {
		panic(err)
	}
}

// AddRuleScript добавляет новое правило, выполняющее JavaScript.
func (engine *RuleEngine) AddRuleScript(id string, eventType string, conditionScript string, script string) {
	engine.RegisterRule(Rule{
		ID:              id,
		EventType:       eventType,
		ConditionScript: conditionScript,
		Script:          script,
		Action: func(event Event) {
			engine.ExecuteScript(script)
		},
	})
}

func (engine *RuleEngine) LoadAndApplyRules(directory string) error {

	// Загрузка правил
	rules, err := LoadRulesFromDirectory(directory, engine.Executor) // Передаем экземпляр JavaScriptExecutor вместо engine.executor

	if err != nil {
		return err
	}

	// Применение загруженных правил
	for _, rule := range rules {
		fmt.Printf("Перед регистрацией правила: ID=%s, EventType=%s, Script=%s\n", rule.ID, rule.EventType, rule.Script)
		engine.RegisterRule(rule)
	}

	return nil
}
