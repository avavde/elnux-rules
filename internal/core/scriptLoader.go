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
		return fmt.Errorf("не удалось прочитать файл скрипта %s: %w", filePath, err)
	}

	scriptName := filepath.Base(filePath)
	rm.scripts[scriptName] = Script{Name: scriptName, Code: string(code)}
	fmt.Printf("Скрипт %s успешно загружен из файла\n", scriptName)
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
	fmt.Printf("Скрипт %s успешно выполнен\n", scriptName)
	return nil
}

// RemoveScript удаляет скрипт из списка.
func (rm *RuleManager) RemoveScript(scriptName string) {
	rm.scriptsMutex.Lock()
	defer rm.scriptsMutex.Unlock()

	delete(rm.scripts, scriptName)
	fmt.Printf("Скрипт %s успешно удален\n", scriptName)
}

// UpdateScript обновляет существующий скрипт новым кодом.
func (rm *RuleManager) UpdateScript(scriptName, newCode string) {
	rm.scriptsMutex.Lock()
	defer rm.scriptsMutex.Unlock()

	if script, exists := rm.scripts[scriptName]; exists {
		script.Code = newCode
		rm.scripts[scriptName] = script
		fmt.Printf("Скрипт %s успешно обновлен\n", scriptName)
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
		return fmt.Errorf("ошибка при маршалинге правила: %w", err)
	}
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("ошибка при записи правила в файл: %w", err)
	}
	fmt.Printf("Правило %s успешно сохранено в файл %s\n", rule.ID, filePath)
	return nil
}

// LoadRulesFromDirectory загружает все правила из указанной директории.
func LoadRulesFromDirectory(directory string, executor *JavaScriptExecutor) ([]Rule, error) {
	fmt.Println("Начало загрузки правил из директории:", directory)
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("ошибка при чтении директории %s: %w", directory, err)
	}

	var rules []Rule
	for _, file := range files {
		filePath := filepath.Join(directory, file.Name())
		data, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Ошибка при чтении файла правила %s: %v\n", filePath, err)
			continue
		}
		var ruleDef RuleDefinition
		if err := json.Unmarshal(data, &ruleDef); err != nil {
			fmt.Printf("Ошибка при разборе JSON правила из файла %s: %v\n", filePath, err)
			continue
		}

		rule := Rule{
			ID:              ruleDef.ID,
			EventType:       ruleDef.Type,
			ConditionScript: ruleDef.Condition,
			ActionScript:    ruleDef.Script,
		}

		rules = append(rules, rule)
		fmt.Printf("Правило %s успешно загружено\n", rule.ID)
	}
	return rules, nil
}
