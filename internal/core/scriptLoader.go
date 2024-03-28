package core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"
)

// Script содержит информацию о пользовательском скрипте.
type Script struct {
	Name string
	Code string
}

// RuleManager управляет правилами и скриптами.
type RuleManager struct {
	executor     *JavaScriptExecutor
	scripts      map[string]Script
	scriptsMutex sync.Mutex
}

// NewRuleManager создает новый экземпляр RuleManager.
func NewRuleManager(executor *JavaScriptExecutor) *RuleManager {
	return &RuleManager{
		executor: executor,
		scripts:  make(map[string]Script),
	}
}

// LoadScriptFromFile загружает скрипт из файла и добавляет его к списку скриптов.
func (rm *RuleManager) LoadScriptFromFile(filePath string) error {
	rm.scriptsMutex.Lock()
	defer rm.scriptsMutex.Unlock()

	code, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл скрипта: %w", err)
	}

	scriptName := filepath.Base(filePath)
	rm.scripts[scriptName] = Script{Name: scriptName, Code: string(code)}
	return nil
}

// ExecuteScriptByName выполняет скрипт по его имени.
func (rm *RuleManager) ExecuteScriptByName(scriptName string) error {
	rm.scriptsMutex.Lock()
	script, exists := rm.scripts[scriptName]
	rm.scriptsMutex.Unlock()

	if !exists {
		return fmt.Errorf("скрипт %s не найден", scriptName)
	}

	err := rm.executor.Execute(script.Code)
	if err != nil {
		return fmt.Errorf("ошибка при выполнении скрипта %s: %w", scriptName, err)
	}
	return nil
}

// RemoveScript удаляет скрипт из списка.
func (rm *RuleManager) RemoveScript(scriptName string) {
	rm.scriptsMutex.Lock()
	defer rm.scriptsMutex.Unlock()

	delete(rm.scripts, scriptName)
}

// UpdateScript обновляет существующий скрипт новым кодом.
func (rm *RuleManager) UpdateScript(scriptName, newCode string) {
	rm.scriptsMutex.Lock()
	defer rm.scriptsMutex.Unlock()

	if script, exists := rm.scripts[scriptName]; exists {
		script.Code = newCode
		rm.scripts[scriptName] = script
	}
}

// RuleDefinition представляет определение правила для сохранения в файле.
type RuleDefinition struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Condition string `json:"condition"`
	Script    string `json:"script"`
}

// SaveRuleToFile сохраняет правило в JSON-файл.
func SaveRuleToFile(rule RuleDefinition, directory string) error {
	filePath := filepath.Join(directory, rule.ID+".json")
	data, err := json.MarshalIndent(rule, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, data, 0644)
}

// LoadRulesFromDirectory загружает все правила из указанной директории.
func LoadRulesFromDirectory(directory string, executor *JavaScriptExecutor) ([]Rule, error) {
	fmt.Println("Начало загрузки правил из директории:", directory)
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Println("Ошибка при чтении директории:", err)
		return nil, err
	}

	var rules []Rule
	for _, file := range files {
		filePath := filepath.Join(directory, file.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Ошибка при чтении файла правила:", err)
			continue
		}

		var ruleDef RuleDefinition
		if err := json.Unmarshal(data, &ruleDef); err != nil {
			fmt.Println("Ошибка при десериализации JSON правила:", err)
			continue
		}

		// Создание объекта Rule из RuleDefinition
		rule := Rule{
			ID:              ruleDef.ID,
			EventType:       ruleDef.Type, // ПРОВЕРИТЬ ПОЛЯЯЁ!
			ConditionScript: ruleDef.Condition,
			Script:          ruleDef.Script,
		}

		fmt.Printf("Подготовлено к регистрации правило: ID=%s, EventType=%s, Script=%s\n", rule.ID, rule.EventType, rule.Script)
		rules = append(rules, rule)
	}
	return rules, nil
}
