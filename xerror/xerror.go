package xerror

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// Error 是 xerror 的错误类型，实现了 error 接口
type Error struct {
	// 上一级错误（被包装的底层错误）
	Base error

	// 当前错误信息（由 fmt 生成）
	Message string

	// 堆栈信息：xerror.New的位置
	File      string
	ErrorLine int
	Func      string

	// 堆栈信息：函数调用的行
	CallerFile string
	CallerLine int
}

// Error 实现 error 接口，返回当前错误的信息
func (e *Error) Error() string {
	return e.StringWithStack()
	if e.Base != nil {
		return e.Message + ": " + e.Base.Error()
	}
	return e.Message
}

// Unwrap 实现 errors.Unwrap 接口，支持 errors.Is / errors.As
func (e *Error) Unwrap() error {
	return e.Base
}

func New(base error, format string, args ...interface{}) *Error {
	_, callerFile, callerLine, _ := runtime.Caller(2)

	pc, file, errorLine, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		errorLine = 0
	}

	funcName := "unknown"
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			parts := strings.Split(fn.Name(), ".")
			funcName = parts[len(parts)-1]
		}
	}

	msg := fmt.Sprintf(format, args...)

	return &Error{
		Base:       base,
		Message:    msg,
		File:       filepath.Base(file),
		ErrorLine:  errorLine,
		Func:       funcName,
		CallerLine: callerLine,
		CallerFile: filepath.Base(callerFile),
	}
}

func (e *Error) Stack() []string {
	var stacks []string

	var chain []*Error
	current := e
	for current != nil {
		chain = append(
			chain,
			&Error{
				Base:       current.Base,
				Message:    current.Message,
				File:       current.File,
				ErrorLine:  current.ErrorLine,
				Func:       current.Func,
				CallerFile: current.CallerFile,
				CallerLine: current.CallerLine,
			},
		)

		if next, ok := current.Base.(*Error); ok {
			current = next
		} else if current.Base != nil {
			chain = append(
				chain,
				&Error{
					Message: current.Base.Error(),
				},
			)
			break
		} else {
			current = nil
		}
	}

	for _, err := range chain {
		line := fmt.Sprintf(
			"(%s:%d => %s:%d:%s) %s",
			err.CallerFile, err.CallerLine,
			err.File, err.ErrorLine, err.Func,
			err.Message,
		)
		stacks = append(stacks, line)
	}

	return stacks
}

// Error 辅助函数：直接打印完整堆栈（方便调试）
func (e *Error) StringWithStack() string {
	var arr []string
	for i, item := range e.Stack() {
		arr = append(arr, fmt.Sprintf("[%d] %s", i+1, item))
	}
	return strings.Join(arr, "\n")
}

// 确保 *Error 实现了 error 接口
var _ error = (*Error)(nil)

/*
package main

import (
	"fmt"
	"miego/xerror"
)

func Layer0() error {
	return xerror.New(nil, "Layer0: database connection timeout")
}

func Layer1() error {
	if err := Layer0(); err != nil {
		return xerror.New(err, "Layer1: query failed")
	}
	return nil
}

func Layer2() error {
	if err := Layer1(); err != nil {
		return xerror.New(err, "Layer2: failed to process")
	}
	return nil
}

func main() {
	err := Layer2()
	if err != nil {
		fmt.Println("// 方式1：打印错误信息（含链）.................")
		fmt.Println(err)

		fmt.Println("// 方式2：打印堆栈列表.................")
		if ne, ok := err.(*xerror.Error); ok {
			stacks := ne.Stack()
			for i, s := range stacks {
				fmt.Printf("%d: %s\n", i+1, s)
			}
		}

		fmt.Println("// 方式3：直接打印完整堆栈（方式2的简化）.................")
		if ne, ok := err.(*xerror.Error); ok {
			fmt.Println(ne.StringWithStack())
		}
	}
}


:!go run main.go
// 方式1：打印错误信息（含链）.................
Layer2: failed to process: Layer1: query failed: Layer0: database connection timeout
// 方式2：打印堆栈列表.................
1: main.go:Layer0:9 (~13): Layer0: database connection timeout
2: main.go:Layer1:14 (~20): Layer1: query failed
3: main.go:Layer2:21 (~27): Layer2: failed to process
// 方式3：直接打印完整堆栈（方式2的简化）.................
main.go:Layer0:9 (~13): Layer0: database connection timeout
main.go:Layer1:14 (~20): Layer1: query failed
main.go:Layer2:21 (~27): Layer2: failed to process
*/
