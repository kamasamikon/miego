//go:build windows

package klog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32")
var outputDebugStringW = kernel32.NewProc("OutputDebugStringW")

func OutputDebugString(s string) {
	p, err := syscall.UTF16PtrFromString(s)
	if err == nil {
		outputDebugStringW.Call(uintptr(unsafe.Pointer(p)))
	}
}

// KLogLN : Log with CR
func KLogLN(dep int, shortPath bool, color string, class string, formating string, args ...interface{}) {
	if Conf.Mute {
		return
	}

	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(dep)

	if ok {
		funcname = runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")
	}

	if shortPath {
		filename = filepath.Base(filename)
	}

	aa := fmt.Sprintf(formating, args...)
	bb := fmt.Sprintf("|%s|F:%s|H:%s|L:%d| %s\n", class, filename, funcname, line, aa)
	OutputDebugString(bb)
}

// KLogLNS : Log with CR
func KLogLNS(dep int, shortPath bool, color string, class string, formating string, args ...interface{}) string {
	if Conf.Mute {
		return ""
	}

	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(dep)

	if ok {
		funcname = runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")
	}

	if shortPath {
		filename = filepath.Base(filename)
	}

	aa := fmt.Sprintf(formating, args...)
	return fmt.Sprintf("|%s|F:%s|H:%s|L:%d| %s\n", class, filename, funcname, line, aa)
}

// KLogX : No '\s' appended.
func KLog(dep int, shortPath bool, color string, class string, formating string, args ...interface{}) {
	if Conf.Mute {
		return
	}

	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(dep)

	if ok {
		funcname = runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")
	}

	if shortPath {
		filename = filepath.Base(filename)
	}

	aa := fmt.Sprintf(formating, args...)
	bb := fmt.Sprintf("|%s|F:%s|H:%s|L:%d| %s\n", class, filename, funcname, line, aa)
	OutputDebugString(bb)
}
