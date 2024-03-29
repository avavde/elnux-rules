package core

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/dop251/goja"
)

type JavaScriptExecutor struct {
	// Это поле больше не нужно, так как вы создаете новый runtime на каждый вызов
}

func NewJavaScriptExecutor() *JavaScriptExecutor {
	return &JavaScriptExecutor{}
}

func (executor *JavaScriptExecutor) initializeConsole(runtime *goja.Runtime) {

	runtime.Set("console", map[string]interface{}{
		"log": func(call goja.FunctionCall) goja.Value {
			fmt.Println("JS Log:", call.Argument(0).String())
			return nil
		},
	})
}

func (executor *JavaScriptExecutor) Execute(script string) error {
	// Создание нового экземпляра для каждого вызова
	runtime := goja.New()
	executor.initializeConsole(runtime)
	time.Sleep(1 * time.Second)
	_, err := runtime.RunString(script)
	if err != nil {
		fmt.Println("Ошибка при выполнении скрипта:", err)
	}
	return err
}

// EvaluateCondition оценивает JavaScript условие и возвращает bool результат.
func (executor *JavaScriptExecutor) EvaluateCondition(conditionScript string, eventPayload map[string]interface{}) bool {
	// Создание нового экземпляра для каждого вызова
	runtime := goja.New()
	executor.initializeConsole(runtime)

	payloadStr, err := json.MarshalIndent(eventPayload, "", "  ")
	if err != nil {
		log.Printf("Ошибка при преобразовании eventPayload в JSON: %v", err)
		return false
	}
	fmt.Printf("\neventPayload: %s\n", string(payloadStr)) // Исправлено для корректного вывода строки JSON
	fmt.Printf("\nevent.Payload: %s\n", eventPayload)
	if eventPayload == nil {
		log.Println("eventPayload равна nil")
		return false
	}

	runtime.Set("event", eventPayload)
	value, err := runtime.RunString(conditionScript)
	if err != nil {
		log.Printf("Ошибка при оценке условия: %v", err)
		return false
	}

	result := value.ToBoolean()
	if !result {
		log.Println("Условие не вернуло булево значение.")
		return false
	}

	return result
}
