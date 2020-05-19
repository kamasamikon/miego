package klog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/mgutz/ansi"
)

var ColorType_F = ansi.ColorCode("red+b:black")
var ColorType_A = ansi.ColorCode("red+h:black")
var ColorType_C = ansi.ColorCode("cyan+b:black")
var ColorType_E = ansi.ColorCode("cyan+h:black")
var ColorType_W = ansi.ColorCode("yellow+b:black")
var ColorType_N = ansi.ColorCode("yellow+h:black")
var ColorType_I = ansi.ColorCode("green+b:black")
var ColorType_D = ansi.ColorCode("green+h:black")
var ColorType_Reset = ansi.ColorCode("reset")

var Conf struct {
	ShortPath bool
	NoColor   bool
}

func KLog(dep int, color string, class string, formating string, args ...interface{}) {
	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(dep)

	if ok {
		funcname = runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")
	}

	cEnd := ColorType_Reset
	if Conf.ShortPath {
		filename = filepath.Base(filename)
	}
	if Conf.NoColor {
		color = ""
		cEnd = ""
	}

	now := time.Now()
	nowQ := now.Format("2006/01/02 15:04:05.")
	nowH := now.Nanosecond() / 1000 / 1000 % 1000
	fmt.Printf("%s|%s|S:%s%03d|F:%s|H:%s|L:%d|%s %s\n", color, class, nowQ, nowH, filename, funcname, line, cEnd, fmt.Sprintf(formating, args...))
}

// F :Fatal
func F(formating string, args ...interface{}) {
	KLog(2, ColorType_F, "F", formating, args...)
}

// A :Alert
func A(formating string, args ...interface{}) {
	KLog(2, ColorType_A, "A", formating, args...)
}

// C :Critical conditions
func C(formating string, args ...interface{}) {
	KLog(2, ColorType_C, "C", formating, args...)
}

// E :Error
func E(formating string, args ...interface{}) {
	KLog(2, ColorType_E, "E", formating, args...)
}

// W :Warning
func W(formating string, args ...interface{}) {
	KLog(2, ColorType_W, "W", formating, args...)
}

// N :Notice
func N(formating string, args ...interface{}) {
	KLog(2, ColorType_N, "N", formating, args...)
}

// I :Information
func I(formating string, args ...interface{}) {
	KLog(2, ColorType_I, "I", formating, args...)
}

// D :Debug message
func D(formating string, args ...interface{}) {
	KLog(2, ColorType_D, "D", formating, args...)
}

func Dump(obj interface{}) {
	s := spew.Sdump(obj)
	KLog(2, ColorType_D, "D", s)
}
