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
	CallerLine int

	// 链式结构：指向下一级（更底层的错误）
	// 注意：设计为链表，方便遍历
}

// Error 实现 error 接口，返回当前错误的信息
func (e *Error) Error() string {
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
	_, _, callerLine, _ := runtime.Caller(2)

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
		File:       filepath.Base(file), // 只保留文件名，不保留完整路径
		ErrorLine:  errorLine,
		Func:       funcName,
		CallerLine: callerLine,
	}
}

// Stack 返回从最底层到当前层的完整堆栈信息（字符串切片）
// 每一条格式为："[函数名] 文件名:行号: 错误信息"
func (e *Error) Stack() []string {
	var stacks []string

	// 1. 收集所有层级（从最底层到最顶层）
	var chain []*Error
	current := e
	for current != nil {
		chain = append([]*Error{current}, chain...) // 头插法，使顺序为 底层 → 顶层
		// 尝试解包到下一层
		if next, ok := current.Base.(*Error); ok {
			current = next
		} else if current.Base != nil {
			// 如果 Base 是普通 error（非 xerror），直接作为叶子节点
			// 但我们不把它加入 chain，因为普通 error 没有堆栈信息
			break
		} else {
			current = nil
		}
	}

	// 2. 遍历 chain，生成堆栈字符串
	for _, err := range chain {
		// 格式：函数名 文件名:(调用行号~错误行号): 错误信息
		line := fmt.Sprintf("%s:%s:%d (~%d): %s", err.File, err.Func, err.ErrorLine, err.CallerLine, err.Message)
		stacks = append(stacks, line)
	}

	return stacks
}

// Error 辅助函数：直接打印完整堆栈（方便调试）
func (e *Error) StringWithStack() string {
	stacks := e.Stack()
	return strings.Join(stacks, "\n")
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
