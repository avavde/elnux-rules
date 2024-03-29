package core

import (
	"fmt"

	"github.com/dop251/goja"
)

func (executor *JavaScriptExecutor) RegisterAPI(runtime *goja.Runtime) {
	runtime.Set("log", func(call goja.FunctionCall) goja.Value {
		msg := call.Argument(0).String()
		fmt.Println("JavaScript log:", msg)
		return nil
	})

	runtime.Set("warn", func(call goja.FunctionCall) goja.Value {
		msg := call.Argument(0).String()
		fmt.Println("JavaScript warning:", msg)
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
