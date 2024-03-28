package core

import (
	"fmt"

	"github.com/dop251/goja"
)

// JavaScriptExecutor структура, отвечающая за исполнение JavaScript.
type JavaScriptExecutor struct {
	runtime *goja.Runtime
}

func NewJavaScriptExecutor() *JavaScriptExecutor {
	executor := &JavaScriptExecutor{
		runtime: goja.New(),
	}
	executor.runtime.Set("console", map[string]interface{}{
		"log": func(call goja.FunctionCall) goja.Value {
			fmt.Println("LOG:", call.Argument(0).String())
			return nil
		},
		"warn": func(call goja.FunctionCall) goja.Value {
			fmt.Println("WARNING:", call.Argument(0).String())
			return nil
		},
	})
	return executor
}

// Execute выполняет JavaScript код.
func (executor *JavaScriptExecutor) Execute(script string) error {
	_, err := executor.runtime.RunString(script)
	return err
}
