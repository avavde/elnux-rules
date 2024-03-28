package core

import (
	"fmt"
	"github.com/dop251/goja"
)

func (executor *JavaScriptExecutor) RegisterAPI() {
	executor.runtime.Set("log", func(call goja.FunctionCall) goja.Value {
		msg := call.Argument(0).String()
		// Печатаем сообщение в лог системы
		fmt.Println("JavaScript log:", msg)
		executor.runtime.Set("warn", warnFunc)
		executor.runtime.Set("log", logFunc)
		return nil
	})


}

func logFunc(call goja.FunctionCall) goja.Value {
    fmt.Println("LOG:", call.Argument(0).String())
    return nil
}

func warnFunc(call goja.FunctionCall) goja.Value {
    fmt.Println("WARNING:", call.Argument(0).String())
    return nil
}