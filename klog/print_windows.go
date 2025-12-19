//go:build windows

package klog

import (
	"io"
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32")
var outputDebugStringW = kernel32.NewProc("OutputDebugStringW")

var debugViewWriter io.Writer

type DebugViewWriter struct {
}

func (w *DebugViewWriter) Write(s []byte) (n int, err error) {
	p, err := syscall.UTF16PtrFromString(string(s))
	if err == nil {
		outputDebugStringW.Call(uintptr(unsafe.Pointer(p)))
	}
	return len(s), nil
}

func init() {
	debugViewWriter = &DebugViewWriter{}
	WriterAdd("debugView", debugViewWriter)
}
